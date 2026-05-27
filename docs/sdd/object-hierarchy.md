# System Design Document: Object Hierarchy

**Author**: slyao (via long-haul agent)  
**Created**: 2026-05-25  
**Status**: Draft  
**Cross-references**: [Rules Engine SDD](rules-engine.md)

---

## 1. Overview

The Object Hierarchy subsystem manages **tree-structured resource ownership** within caskin. Every `Object` may have a parent object, forming a rooted tree per domain per object type. Permissions on a parent implicitly cascade to all descendants via casbin's `g2` grouping policy.

This is the foundational mechanism that enables hierarchical permission inheritance — a user with "manage" on a folder implicitly has "manage" on all sub-folders within.

---

## 2. Goals & Non-Goals

### Goals

- Provide a tree-structured ownership model for resources
- Enforce safety invariants (no cycles, bounded depth, type consistency)
- Keep GORM metadata and casbin policy in sync as a dual-store system
- Support directory operations (search, move, delete subtrees)
- Allow generic graph utilities (topological sort, edge ordering)

### Non-Goals

- Cross-type hierarchies (parent and child must share the same `ObjectType`)
- Graph structures beyond trees (DAGs or arbitrary digraphs)
- Lazy/on-demand hierarchy loading (full domain objects are loaded per operation)
- Caching of transitive closures

---

## 3. Architecture

### 3.1 Dual-Store Design

```
Object.parentInterface
       │
       ▼
┌─────────────────────┐
│  GORM Metadata DB   │   parent_id column on object records
│  (source of truth)  │
└──────────┬──────────┘
           │ objectUpdater syncs
           ▼
┌─────────────────────┐
│  casbin g2 rules    │   g2(child_encode, parent_encode, domain_encode)
│  (derived state)    │   used for permission evaluation
└─────────────────────┘
```

**Core Invariant**: `parent_id` in GORM ≡ corresponding `g2` rule in casbin. `objectUpdater.Run()` maintains this by diffing current g2 parents against the new `parent_id`.

### 3.2 Key Interfaces

| Interface | File | Purpose |
|-----------|------|---------|
| `parentInterface` | `schema_private.go` | `GetParentID()` / `SetParentID()` |
| `Object` | `schema.go` | Embeds `parentInterface` + `objectTypeInterface` |
| `IDirectory` | (implicit in `server_directory.go`) | Directory CRUD operations |

### 3.3 Key Implementation Files

| File | Responsibility |
|------|---------------|
| `inheritance.go` | Generic `InheritanceEdge[T]`, `InheritanceGraph[T]`, `TopSort`, `EdgeSorter` |
| `object_updater.go` | Syncs parent changes to casbin g2 rules (add/remove) |
| `object_directory.go` | Tree construction from flat object list, BFS search, DFS counting |
| `server_object.go` | Object CRUD with hierarchy validation |
| `server_check.go` | Safety guards (cycle, depth, type, permission checks) |
| `server_directory.go` | Directory-level operations (create/update/delete/get/move) |
| `casbin.go` | `EnforceObject`, `GetParentsForObjectInDomain`, `GetChildrenForObjectInDomain` |

---

## 4. Object Tree Operations

### 4.1 Create

1. `ObjectParentCheck` — validates parent exists, same type, user has manage permission on parent
2. `ObjectHierarchyCheck` — BFS depth validation (total depth ≤ 10)
3. Insert into GORM
4. `objectUpdater.Run()` — add g2 rule `(child, parent, domain)`

### 4.2 Update (re-parent)

1. `ObjectUpdateCheck` — no self-loop, type immutability, manage permission on old parent
2. `ObjectParentToDescendantCheck` — cycle detection via `EnforceObject`
3. `ObjectParentCheck` — validate new parent
4. `ObjectHierarchyCheck` — depth validation with new parent
5. Update GORM `parent_id`
6. `objectUpdater.Run()` — remove old g2 rule, add new g2 rule

### 4.3 Delete

1. `DeleteDirectory` collects all descendants via BFS (`DirectorySearchAll`)
2. Executes `ActionDirectory` callback on all collected IDs (application-level cleanup)
3. Soft-deletes the root object via `DeleteObject`
4. Casbin g2 rules cleaned up by `objectUpdater`

### 4.4 Move

`MoveDirectory` = Update with re-parenting; same validation chain applies.

---

## 5. Safety Guards & Invariants

| Guard | Location | Check | Error |
|-------|----------|-------|-------|
| Self-loop prevention | `ObjectUpdateCheck` | `object.GetID() == object.GetParentID()` | `ErrParentCanNotBeItself` |
| Descendant→ancestor cycle | `ObjectParentToDescendantCheck` | `EnforceObject(object, newParent)` — if child already inherits new parent, it's a cycle | `ErrParentToDescendant` |
| Max depth 10 | `ObjectHierarchyCheck` | BFS up from parent + BFS down from object; sum > 10 → error | `"max directory depth ... too large than 10"` |
| Type consistency | `ObjectParentCheck` | `parent.ObjectType != child.ObjectType` | `ErrInValidObjectType` |
| Root protection | `ObjectParentCheck` | `parentID == 0` | `ErrCantOperateRootObject` |
| Type immutability | `ObjectUpdateCheck` | cannot change type after creation | `ErrCantChangeObjectType` |
| Permission gate | `ObjectParentCheck` | user must have `manage` on parent | `ErrNoWritePermission` |

### Invariant Enforcement Order

Guards are evaluated **before** any mutation, ensuring the tree can never enter an invalid state. The enforcement order is: permission → self-loop → type → cycle → depth → mutation.

---

## 6. Performance Analysis

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| Create/Update depth check | O(tree_size) | 2× BFS traversal per operation |
| objectUpdater sync | O(1) | Single g2 rule add/remove |
| Directory search (all) | O(N) | N = total objects in domain of that type |
| Directory tree build | O(N) | Single pass to construct adjacency map |
| DFS counting | O(N) | Aggregate counts for all nodes |
| Permission evaluation | O(g2_rules) | Casbin walks full g2 table per Enforce |

### Known Bottlenecks

1. **`GetParentsForObjectInDomain` / `GetChildrenForObjectInDomain`** enumerate the full g2 table for the domain — no index on parent_id within casbin's in-memory model.
2. **`ObjectHierarchyCheck`** performs 2 BFS traversals per create/update — acceptable for depth ≤ 10 but costly if called in tight loops.
3. **No batched hierarchy validation** — moving N objects requires N independent BFS checks.
4. **No caching of transitive closure** — every `EnforceObject` recomputes hierarchy traversal.

---

## 7. Directory Service

### 7.1 `objectDirectory` (internal)

Constructed from a flat `[]*Directory` slice via `NewObjectDirectory`:

- `Tree map[uint64][]*Directory` — children indexed by parent ID
- `Node map[uint64]*Directory` — all nodes indexed by own ID
- `visit map[uint64]bool` — DFS visited tracker

### 7.2 Search Types

| Type | Behavior |
|------|----------|
| `DirectorySearchTop` (default) | Direct children of target only |
| `DirectorySearchAll` | BFS to collect entire subtree |

### 7.3 DFS Counting

`dfs()` recursively aggregates:
- `TopDirectoryCount` — direct child directories
- `AllDirectoryCount` — total descendant directories
- `AllItemCount` — total items (leaf data count provided by `CountDirectory` callback)

---

## 8. InheritanceGraph Utilities

Generic graph utilities in `inheritance.go` support both role and object hierarchies:

| Type / Function | Purpose |
|-----------------|---------|
| `InheritanceEdge[T]` | Directed edge (U→V) with JSON encode/decode |
| `InheritanceGraph[T]` | Adjacency-list graph (map[T][]T) |
| `TopSort()` | Kahn's algorithm BFS topological sort |
| `EdgeSorter[T]` | Sort edges root-first or leaf-first based on topo order |
| `MergeInheritanceGraph` | Combine multiple graphs, deduplicate, sort |
| `distinct[T]` | Helper to deduplicate ordered slices |

These utilities are used by `objectUpdater` and role updaters to determine the correct order of g2 rule modifications (root-first for inserts, leaf-first for deletes).

---

## 9. Extension Points

1. **`ActionDirectory` callback** — application layer defines how to clean up items when a directory subtree is deleted.
2. **`CountDirectory` callback** — application layer provides item counts per object ID for directory listing.
3. **`parentInterface`** — any struct embedding this interface can participate in the hierarchy.
4. **Custom `ObjectType`** — new resource types automatically get hierarchy support by implementing the Object interface.
5. **`InheritanceGraph[T cmp.Ordered]`** — generic over any ordered type, reusable beyond caskin's uint64 IDs.

---

## 10. Known Limitations & Future Opportunities

### Limitations

- **Type-homogeneous trees only** — parent and child must share the same `ObjectType`
- **Single parent** — true tree, not DAG; an object can have exactly one parent
- **Global depth limit** — hardcoded to 10, not configurable per domain/type
- **No incremental sync** — `objectUpdater` reads all g2 rules for comparison on each run
- **Cycle detection depends on casbin state** — if g2 is out of sync, `ObjectParentToDescendantCheck` may give wrong results

### Future Opportunities

| # | Opportunity | Impact |
|---|-------------|--------|
| 1 | **Materialized path / nested set encoding** | O(1) ancestor queries, eliminates BFS depth checks |
| 2 | **Lazy depth validation** | Skip BFS if parent hasn't changed (common case) |
| 3 | **Batch move optimization** | Single BFS for entire subtree move instead of per-object |
| 4 | **Configurable depth limit** | Per-domain or per-type max depth via configuration |
| 5 | **Cross-type hierarchies** | Mixed resource trees for complex permission models |
| 6 | **Explicit orphan handling** | Policy-driven cascade (reassign to grandparent vs error vs soft-delete) |
| 7 | **g2 transitive closure cache** | Precomputed ancestry sets invalidated on tree mutation |

---

## 11. Testing Strategy

| Layer | What to test | Approach |
|-------|-------------|----------|
| Unit | `InheritanceGraph.TopSort`, `EdgeSorter`, `objectDirectory.Search` | Table-driven tests with known graphs |
| Integration | Safety guards (cycle, depth, type) | Build tree in test DB, attempt invalid mutations |
| Integration | objectUpdater sync correctness | Assert g2 rules match GORM parent_id after operations |
| E2E | Directory CRUD with permission checks | Full server setup with user/domain/policy fixtures |

---

## 12. Glossary

| Term | Definition |
|------|-----------|
| **g2** | Casbin grouping policy type for object inheritance (child, parent, domain) |
| **objectUpdater** | Component that syncs GORM parent_id changes to casbin g2 rules |
| **Directory** | Wrapper around Object adding tree-aggregate counts (items, subdirectories) |
| **BFS depth check** | Bi-directional BFS to compute total tree depth at a node |
| **EdgeSorter** | Utility that orders inheritance edges by topological position |
| **parentInterface** | Private interface providing GetParentID/SetParentID |
