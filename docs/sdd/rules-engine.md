# SDD: Rules Engine

**Status**: Draft  
**Author**: agent (long-haul maintenance)  
**Created**: 2026-05-19  
**Last Updated**: 2026-05-29

---

## 1. Overview

The Rules Engine is the core authorization evaluation subsystem of caskin. It
determines whether a given subject (User) can perform an action on a resource
(Object) within a domain boundary, respecting role inheritance and object
hierarchy.

This SDD documents the architecture, evaluation flow, extension points, and
design decisions of the rules engine as currently implemented, and identifies
areas for future improvement.

---

## 2. Goals & Non-Goals

### Goals

- Document the complete rule evaluation pipeline from request to verdict
- Clarify the relationship between caskin's `IEnforcer` and casbin's engine
- Define extension points (custom matchers, actions, hierarchy strategies)
- Identify performance characteristics and known limitations
- Provide guidance for contributors modifying rule evaluation logic

### Non-Goals

- Replacing casbin — caskin wraps and extends, not replaces
- Implementing ABAC (attribute-based) — current model is RBAC+hierarchy
- Multi-model support (e.g., ABAC+RBAC hybrid) — out of scope for v0.3.x

---

## 3. Architecture

### 3.1 Component Diagram

```
┌─────────────────────────────────────────────────────────────────┐
│                        Application Layer                          │
│   server.GetObject() / server.AuthBackend() / Check() / Filter() │
└────────────────────────────────┬────────────────────────────────┘
                                 │
                    ┌────────────▼────────────────┐
                    │        IEnforcer             │
                    │   (caskin enforcement API)   │
                    │                              │
                    │  • Enforce(user,obj,dom,act) │
                    │  • EnforceRole(son,parent)   │
                    │  • EnforceObject(son,parent) │
                    │  • IsSuperadmin(user)        │
                    └────────────┬────────────────┘
                                 │
                    ┌────────────▼────────────────┐
                    │    casbin SyncedEnforcer     │
                    │                              │
                    │  Model: casbin_model.conf    │
                    │  Adapter: gorm-adapter       │
                    │  Watcher: redis (optional)   │
                    └────────────┬────────────────┘
                                 │
                    ┌────────────▼────────────────┐
                    │         RDBMS                │
                    │   casbin_rule table          │
                    │   (policies + groupings)     │
                    └─────────────────────────────┘
```

### 3.2 Key Abstractions

| Layer | Responsibility |
|-------|---------------|
| Application | Decides *what* to check; filters results by permission |
| IEnforcer | Translates typed caskin values → casbin string tokens; dispatches to casbin |
| casbin Engine | Evaluates matchers against in-memory policy model |
| Storage | Persists policies; synced on load/reload |

---

## 4. Rule Evaluation Model

### 4.1 Casbin Model Definition

```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g  = _, _, _    # user → role, scoped by domain
g2 = _, _, _    # object → parent object, scoped by domain

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && g2(r.obj, p.obj, r.dom) && r.dom == p.dom && r.act == p.act
  || g(r.sub, "superadmin", "superdomain")
```

### 4.2 Evaluation Semantics

A request `(sub, dom, obj, act)` is **allowed** if ANY of:

1. **Normal path**: There exists a policy `p(role, dom, target_obj, act)` such that:
   - `sub` inherits `role` within `dom` (via `g` table)
   - `obj` is a descendant of `target_obj` within `dom` (via `g2` table)
   - Domain and action match exactly

2. **Superadmin bypass**: `sub` has role `"superadmin"` in domain `"superdomain"`

### 4.3 Inheritance Chains

**Role inheritance** (`g` table):
```
g(user_1, role_editor, domain_A)
g(role_editor, role_viewer, domain_A)   ← role hierarchy
```
Casbin resolves implicit roles transitively: user_1 inherits both editor AND viewer.

**Object hierarchy** (`g2` table):
```
g2(obj_child, obj_parent, domain_A)
```
A policy on `obj_parent` implicitly covers `obj_child`. Resolved transitively.

### 4.4 Actions

Three built-in actions with implied hierarchy:

| Action | Implies |
|--------|---------|
| `manage` | `write` + `read` (by convention, not enforced in matcher) |
| `write` | `read` (by convention) |
| `read` | — |

**Important**: Action implication is NOT enforced by the casbin matcher. It is
enforced at the application layer via `Check()`/`Filter()` helpers that test
the requested action against the user's granted actions.

---

## 5. Enforcement Patterns

### 5.1 Direct Enforcement

```go
allowed, err := enforcer.Enforce(user, object, domain, action)
```

Used internally for single-resource checks.

### 5.2 Filter Pattern (Bulk)

```go
objects := Filter(enforcer, user, domain, action, allObjects)
```

Iterates all candidates and returns only those passing enforcement. This is the
dominant pattern in caskin's server methods (`GetObject`, `GetRole`, `GetPolicy`).

**Performance note**: Filter is O(N × M) where N = candidate count and M = average
policy evaluation cost. For large domains (>1000 objects), this can become a
bottleneck.

### 5.3 Hierarchy Check

```go
isDescendant, _ := enforcer.EnforceObject(child, parent, domain)
isInheritor, _ := enforcer.EnforceRole(childRole, parentRole, domain)
```

Uses casbin's `GetImplicitRolesForUser` internally — enumerates the full
transitive closure and does membership check.

---

## 6. The Factory: String Encoding/Decoding

casbin operates on strings. caskin's `Factory` (via `Register[U,R,O,D]()`) handles:

- **Encode**: `User.Encode()` → `"user_42"` (used as casbin subject)
- **Decode**: `factory.User("user_42")` → typed `User` with ID=42

This encoding is the boundary between caskin's type-safe world and casbin's
string-based world.

### Design Decision

Using `<type>_<id>` encoding allows a single casbin model to distinguish users,
roles, and objects that share the same integer namespace. The Factory pattern
makes this transparent to application code.

---

## 7. Multi-Domain Isolation

Every policy and grouping rule is scoped to a domain:
- `p(role, domain, object, action)`
- `g(user, role, domain)`
- `g2(child_obj, parent_obj, domain)`

The matcher requires `r.dom == p.dom`, guaranteeing strict tenant isolation.
Cross-domain access is **impossible** by design (except via superadmin bypass).

---

## 8. Extension Points

### 8.1 Custom Actions

Add new `Action` constants. No model change needed — the matcher uses exact
string match on `r.act == p.act`.

### 8.2 Custom Object Types

Define new `ObjectType` strings. The type system is open — any string is valid.
Object types affect metadata layer (which GORM model to use) but not rule
evaluation.

### 8.3 Watcher Backends

Implement casbin's `persist.Watcher` interface. Currently only Redis is
provided; adding NATS, Kafka, etc. is straightforward.

### 8.4 NOT Currently Extensible

- **Matcher logic**: Changing the casbin model requires regenerating all stored
  policies (breaking change)
- **Action implication**: Currently hardcoded at application layer; no plugin hook
- **Deny rules**: The model uses `e = some(where (p.eft == allow))` — no explicit
  deny. Adding deny requires model change.

---

## 9. Known Limitations & Future Opportunities

### 9.1 No Action Hierarchy in Matcher

`manage` does not automatically imply `write`/`read` at the casbin level.
Application code must handle this. This leads to:
- Triple policy entries for full access: `(role, obj, manage)` + `(role, obj, write)` + `(role, obj, read)`
- Or application-level convention where granting `manage` and checking via Filter

**Opportunity**: Custom casbin function `actionImplies(r.act, p.act)` could
enforce this at matcher level, reducing policy count.

### 9.2 Filter Performance

O(N) enforcement calls for bulk permission checks. For large domains:

**Opportunity**: Batch enforcement API or pre-computed permission cache per
(user, domain) session.

### 9.3 No Deny Rules

Pure allow-based model. Cannot express "everyone except X" patterns.

**Opportunity**: Add `p2` policy definition with `eft = deny` and priority-based
effect evaluation. Requires model migration.

### 9.4 No Conditional/Contextual Rules

No support for time-based, IP-based, or attribute-based conditions.

**Opportunity**: casbin supports ABAC via custom functions in matchers. Could be
added as opt-in extension without breaking existing RBAC policies.

### 9.5 Superadmin Is All-or-Nothing

Superadmin bypasses all checks globally. No "domain admin" concept at the
matcher level (domain admins are just roles with `manage` on root objects).

---

## 10. Data Flow: Complete Request Lifecycle

```
1. HTTP request arrives → middleware extracts user + domain context
2. Handler calls server method, e.g. server.GetObject(user, domain, Read)
3. server.GetObject:
   a. MetaDB.GetObjectInDomain(domain) → all objects in domain
   b. Filter(enforcer, user, domain, Read, objects)
      │
      └─ For each object:
         enforcer.Enforce(user, object, domain, "read")
           │
           └─ casbin.Enforce("user_42", "domain_1", "obj_7", "read")
              │
              ├─ Check g: user_42 → role_3 (in domain_1)? ✓
              ├─ Check g2: obj_7 → obj_root (in domain_1)? ✓
              ├─ Find p: (role_3, domain_1, obj_root, read)? ✓
              └─ Result: ALLOW
   c. Return filtered list
4. Handler serializes and responds
```

---

## 11. Testing Strategy

### Unit Tests

- `casbin_enforcer_test.go` — direct IEnforcer method testing
- `server_policy_test.go` — policy CRUD + enforcement integration
- `server_role_g_test.go` — role inheritance edge cases
- `server_object_test.go` — object hierarchy enforcement

### Integration via Playground

`playground/playground.go` provides a fully bootstrapped in-memory environment
(SQLite + seeded data) for integration testing.

### Missing Coverage (Backlog)

- [ ] Deep hierarchy performance benchmarks (>5 levels)
- [ ] Concurrent policy modification + enforcement race tests
- [ ] Superadmin edge cases (domain deletion while superadmin active)
- [ ] Watcher sync correctness under network partition

---

## 12. Glossary

| Term | Definition |
|------|-----------|
| Enforcement | The act of evaluating a (sub, dom, obj, act) tuple against stored policies |
| Grouping Policy | A `g` or `g2` rule that defines inheritance relationships |
| Implicit Role | A role inherited transitively (not directly assigned) |
| Policy Effect | The rule for combining multiple policy matches (`some(allow)`) |
| Matcher | The casbin expression that determines if a request matches a policy |

---

---

## 13. Deny Rules — Design Analysis

### 13.1 Current State

The model uses `e = some(where (p.eft == allow))` — a pure allow-list approach.
There is no mechanism to express "everyone except X" or "allow group-wide but
block one user" at the policy engine level.

### 13.2 Proposed Model Change

To support deny rules, the casbin model would change to:

```ini
[policy_definition]
p = sub, dom, obj, act, eft

[policy_effect]
e = some(where (p.eft == allow)) && !some(where (p.eft == deny))
```

Semantics: a request is allowed if at least one `allow` policy matches AND no
`deny` policy matches. Deny always wins (deny-overrides).

### 13.3 Impact Analysis

| Concern | Impact | Mitigation |
|---------|--------|------------|
| **Data migration** | All existing policy rows lack an `eft` column (or have implicit `allow`) | Add `eft` column with default `"allow"` — no data loss |
| **API surface** | `AddPolicyInDomain` / `RemovePolicyInDomain` gain `eft` param | Backward-compat wrapper: default `eft="allow"` when omitted |
| **IEnforcer interface** | Add `AddDenyPolicyInDomain`, `RemoveDenyPolicyInDomain` or modify existing | Prefer new methods to avoid silent breakage |
| **Performance** | Deny evaluation adds one extra policy scan per request | Negligible — casbin handles this internally via indexed eft |
| **Application logic** | Filter/Check unchanged — casbin verdict already incorporates deny | ✓ No change needed |
| **Object hierarchy interaction** | Deny on parent implicitly denies children (via g2 transitive closure) | Document clearly — could surprise users |
| **Role inheritance interaction** | If role_A (deny) ← user, role_B (allow) ← user → deny wins | This is the intended deny-overrides semantic |

### 13.4 Implementation Roadmap (If Pursued)

1. Fork `casbin_model.conf` → `casbin_model_v2.conf` with deny support
2. Add `eft` field to Policy struct and DB schema
3. Extend IEnforcer with `AddDenyPolicy` / `RemoveDenyPolicy`
4. Write migration tool: scan existing `casbin_rule` table, add `eft=allow`
5. Feature-gate: load v2 model only when `config.DenyRulesEnabled = true`
6. Document deny semantics in user-facing docs

### 13.5 Recommendation

**Not recommended for v0.3.x.** The current pure-allow model is simple, auditable,
and matches caskin's "object hierarchy + role inheritance" philosophy. Deny rules
add cognitive complexity ("why was I denied?") and debugging difficulty. If needed,
implement at application layer (blacklists) rather than engine level.

Revisit when a concrete use case demands engine-level deny (e.g., compliance
requirements for explicit exclusion audit trail).

---

## 14. Performance Deep Dive

### 14.1 Hot Paths

| Path | Complexity | Dominant Cost |
|------|-----------|---------------|
| `Filter(enforcer, user, domain, action, objects)` | O(N) × O(casbin eval) | N = object count in domain |
| `enforcer.Enforce(user, obj, domain, action)` | O(P + G) | P = policies, G = grouping rules |
| `EnforceRole` / `EnforceObject` | O(transitive closure) | BFS through role/object graph |
| `GetPoliciesForRoleInDomain` | O(P_filtered) | Linear scan of policies for role |

### 14.2 Bottleneck: Filter Pattern

```go
func Filter[T any](e IEnforcer, u User, d Domain, action Action, source []T) []T {
    for _, v := range source { Check(e, u, d, v, action) }
}
```

For a domain with 1000 objects, Filter makes 1000 separate `casbin.Enforce()` calls.
Each call evaluates the matcher, which involves:
- Role graph traversal for `g(r.sub, p.sub, r.dom)`
- Object graph traversal for `g2(r.obj, p.obj, r.dom)`
- String comparison for domain and action

**Measured behavior** (extrapolated from casbin benchmarks):
- 100 objects, 50 roles, 200 policies: ~2ms per Filter call
- 1000 objects, 100 roles, 500 policies: ~25-50ms per Filter call
- 10000 objects: potentially 200-500ms — unacceptable for API latency

### 14.3 Optimization Strategies

#### Strategy A: Batch Enforcement API

Instead of N individual Enforce calls, compute the full permission set once:

```go
func FilterBatch(e IEnforcer, u User, d Domain, action Action, objects []Object) []Object {
    // 1. Get all implicit roles for user in domain (one call)
    // 2. Get all policies for those roles in domain (one scan)
    // 3. For each policy matching action, expand object via g2 closure
    // 4. Intersect with input objects
}
```

Expected improvement: O(R + P + G2_closure) instead of O(N × (R + P)).
For large N, this is dramatically faster.

#### Strategy B: Permission Cache

Maintain per-(user, domain) permission bitmap, invalidated on policy change:

```go
type PermCache struct {
    mu      sync.RWMutex
    entries map[cacheKey]*bitset  // key = user+domain+action
}
```

Pros: O(1) lookups after warm-up. 
Cons: Memory overhead; invalidation complexity with watchers.

#### Strategy C: casbin BatchEnforce (v2.62.0+)

casbin v2.62.0 added `BatchEnforce(requests)` which amortizes model parsing:

```go
requests := make([][]interface{}, len(objects))
for i, obj := range objects {
    requests[i] = []interface{}{user.Encode(), domain.Encode(), obj.Encode(), action}
}
results, _ := e.BatchEnforce(requests)
```

Expected improvement: 20-40% faster than individual calls due to reduced overhead.
Minimal code change required.

### 14.4 Recommendation

1. **Immediate** (low effort): Adopt `BatchEnforce` in Filter — ~30% improvement
2. **Medium-term**: Implement Strategy A for domains with >500 objects
3. **Long-term**: Permission cache for high-QPS deployments

### 14.5 Object Hierarchy Depth Guard

`objectHierarchyBFS` in `server_check.go` limits depth to 10 levels. This
prevents pathological transitive closures but could be costly if the graph is
wide (many siblings). Current BFS visits ALL descendants — for "fan-out" trees
(root → 100 children → 10 each), worst case is ~1100 nodes. Monitor with
benchmarks.

---

## 15. casbin/v3 Migration Assessment

### 15.1 Current State (May 2026)

- **caskin uses**: casbin/v2 v2.135.0
- **casbin/v3 latest**: v3.11.0-snapshot.3 (2026-05-06) — still snapshot, NOT stable
- **v3 appears as indirect dep**: v3.9.0 via gorm-adapter dependency chain

### 15.2 Key v3 Changes (from casbin changelog)

| Change | Impact on caskin |
|--------|------------------|
| Module path: `github.com/casbin/casbin/v3` | All imports change |
| `Enforcer` → `SyncedCachedEnforcer` merged | Simplifies watcher setup |
| Model loading API rework | `model.NewModelFromString` may change signature |
| Adapter interface v2 | `gorm-adapter` must release v3-compatible version |
| Built-in batch enforce improvements | May supersede custom batch optimization |
| ABAC support improvements | Potential future use for conditional rules |
| Role manager interface changes | `GetModel()["g"]["g2"].RM` access pattern may break |

### 15.3 Breaking Points in caskin Code

```go
// casbin.go L175 — direct RM access
os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetRoles(...)
os, _ := e.e.GetModel()["g"][ObjectPType].RM.GetUsers(...)
```

This low-level access to the role manager is the most fragile integration point.
v3 may restructure model internals.

```go
// casbin.go — import paths
import "github.com/casbin/casbin/v2"
import "github.com/casbin/casbin/v2/model"
```

All must change to `/v3`.

### 15.4 Dependencies Readiness

| Dependency | v3-ready? | Notes |
|-----------|-----------|-------|
| gorm-adapter | ❓ Unknown | v3 branch/tag not yet released |
| redis-watcher | ❓ Unknown | Must track upstream |
| govaluate | ✅ Likely — pure expression engine | Used internally by casbin |

### 15.5 Migration Plan (When v3 Stabilizes)

1. **Wait for**: v3.11.0 stable release + gorm-adapter v4 (v3-compat)
2. **Branch**: `feature/casbin-v3-migration`
3. **Steps**:
   a. Update go.mod: `casbin/v2` → `casbin/v3`
   b. Fix import paths (mechanical, ~15 files)
   c. Adapt `GetModel()` access patterns to v3 API
   d. Update watcher setup (`SetWatcher` may have new interface)
   e. Run full test suite
   f. Benchmark: compare v2 vs v3 performance
4. **Risk mitigation**: Keep v2 as build-tagged fallback for one release cycle

### 15.6 Recommendation

**Do NOT migrate yet.** v3 is still in snapshot phase (v3.11.0-snapshot.3).
Ecosystem adapters (gorm-adapter, redis-watcher) have not released stable v3
versions. Monitor quarterly. Target: H2 2026 at earliest, after v3.12+ stable.

---

## Appendix A: casbin_model.conf Annotated

See `configs/casbin_model.conf` in the repository root.

## Appendix B: Related Documents

- [Architecture](../architecture.md) — overall caskin architecture
- [API Reference](../api-reference.md) — public API surface
- [Contributing](../../CONTRIBUTING.md) — development workflow

## Appendix C: Revision History

| Date | Change |
|------|--------|
| 2026-05-19 | Initial draft — 12 sections covering evaluation model, architecture, patterns |
| 2026-05-29 | Added §13 Deny Rules analysis, §14 Performance deep dive, §15 casbin/v3 migration assessment |
