# SDD: Rules Engine

**Status**: Draft  
**Author**: agent (long-haul maintenance)  
**Created**: 2026-05-19  
**Last Updated**: 2026-05-19

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

## Appendix A: casbin_model.conf Annotated

See `configs/casbin_model.conf` in the repository root.

## Appendix B: Related Documents

- [Architecture](../architecture.md) — overall caskin architecture
- [API Reference](../api-reference.md) — public API surface
- [Contributing](../../CONTRIBUTING.md) — development workflow
