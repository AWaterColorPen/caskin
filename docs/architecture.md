# Architecture

This document describes the internal design of caskin for contributors and integrators
who want to understand how the library works under the hood.

## Table of Contents

- [Overview](#overview)
- [Core Concepts](#core-concepts)
- [Layer Architecture](#layer-architecture)
- [Key Interfaces](#key-interfaces)
- [Type Registration System](#type-registration-system)
- [Casbin Integration](#casbin-integration)
- [Permission Model (casbin_model.conf)](#permission-model-casbin_modelconf)
- [Dictionary System](#dictionary-system)
- [Metadata Layer](#metadata-layer)
- [Request Flow: A Permission Check](#request-flow-a-permission-check)
- [Directory: Object Tree Traversal](#directory-object-tree-traversal)
- [Multi-Instance Sync (Redis Watcher)](#multi-instance-sync-redis-watcher)
- [Package Layout](#package-layout)

---

## Overview

Caskin is a **multi-domain RBAC authorization library** for Go, built on top of
[casbin](https://github.com/casbin/casbin). It adds three key capabilities that
plain casbin does not provide:

1. **Domain isolation** — each tenant/organization has its own completely isolated
   role and policy space, enforced at the casbin model level.
2. **Hierarchy-aware objects** — objects (resources) form a tree; a permission
   granted on a parent automatically applies to all descendants.
3. **Dictionary-driven setup** — features, backend APIs, frontend UI items, and
   initial roles/policies are declared in a TOML config file, not hard-coded.

---

## Core Concepts

| Concept     | Type interface | Description |
|-------------|----------------|-------------|
| `User`      | `caskin.User`  | An actor (real person or service account). Identified by an integer ID and a string encoding used as casbin `sub`. |
| `Role`      | `caskin.Role`  | A named permission bundle within one domain. Also an `ObjectData`, so it is bound to an `Object` in the permission tree. |
| `Object`    | `caskin.Object`| A resource or resource group. Forms a tree via `GetParentID()`. The unit of policy assignment. |
| `Domain`    | `caskin.Domain`| An isolated tenant/organization. All RBAC is scoped to a domain. |
| `Policy`    | `*caskin.Policy`| A binding: `(Role, Object, Action)` — "role R can do action A on object O". |
| `ObjectData`| `caskin.ObjectData`| Any application entity (e.g., a document, a role itself) that is protected by an `Object`. |

---

## Layer Architecture

```
┌──────────────────────────────────────────────────┐
│                  Application Code                 │
│         (your handlers, middleware, etc.)          │
└───────────────────────┬──────────────────────────┘
                        │  IService / ICurrentService
┌───────────────────────▼──────────────────────────┐
│                   server (caskin)                  │
│   IBaseService  ·  IFeatureService                 │
│   IDirectoryService  ·  ICurrentService            │
└──────┬────────────────────────────────┬───────────┘
       │ MetaDB (GORM)                  │ IEnforcer
┌──────▼────────┐              ┌────────▼────────────┐
│  RDBMS        │              │  casbin SyncedEnforcer│
│  (SQLite /    │              │  + gorm-adapter       │
│   MySQL /     │              │  + optional watcher   │
│   Postgres)   │              └─────────────────────┘
└───────────────┘
         ▲
┌────────┴──────────┐
│ Dictionary (TOML) │
│  features / pkgs  │
│  backend/frontend │
│  creator_* init   │
└───────────────────┘
```

The `server` struct is the single concrete implementation of `IService`. It holds:
- `Enforcer IEnforcer` — wraps casbin's `SyncedEnforcer` with caskin-specific helpers
- `DB MetaDB` — wraps GORM for domain/user/role/object metadata storage
- `Dictionary IDictionary` — parsed TOML config for features, packages, and initial data

---

## Key Interfaces

### `IService`

The primary interface your application code uses:

```go
type IService interface {
    IBaseService      // CRUD for User, Domain, Role, Object; user-role and policy management
    IDirectoryService // directory (tree) operations for Object and ObjectData
    IFeatureService   // backend/frontend permission management driven by Dictionary
    ICurrentService   // scoped service for a fixed (User, Domain) pair
}
```

`ICurrentService` is returned by `svc.SetCurrent(user, domain)`. It binds a specific
operator and domain, providing convenience methods like `GetCurrentRole()`,
`GetCurrentBackend()`, and `GetCurrentFrontend()` without repeating the user+domain arguments.

### `IEnforcer`

An internal interface wrapping casbin's `SyncedEnforcer`. It exposes:
- `Enforce(user, object, domain, action)` — the core permission check
- Casbin policy/grouping management methods

`IEnforcer` is not exposed via `IService`. If you need raw enforcement (e.g., for
middleware), use the service's high-level methods (`GetObject`, `GetCurrentBackend`, etc.)
or the `caskin.Check` / `caskin.Filter` helpers (which accept an `IEnforcer`).

### `Factory`

Caskin uses a generic type registration system (see next section) to convert
casbin's string tokens back into your concrete `User`, `Role`, `Object`, and `Domain`
values. `Factory` is the interface that drives this.

---

## Type Registration System

Caskin stores casbin policy subjects/objects as encoded strings (e.g., `"user_42"`,
`"role_7"`). To convert these strings back to your concrete types, caskin needs to
know their types at startup.

```go
// Must be called once before caskin.New
caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()
```

Internally, `Register` builds a `Factory` implementation using reflection (via `register.go`).
The factory:
1. Creates zero-value instances of your types (`NewUser()`, `NewRole()`, etc.)
2. Decodes strings into concrete values using each type's `Decode(string) error` method
3. Provides a `MetadataDB` that knows which GORM models correspond to which caskin types

**Common pitfall:** forgetting `Register` before `New` causes a runtime panic.
Every new test setup function must include `caskin.Register[...]()`.

---

## Casbin Integration

### Model (`configs/casbin_model.conf`)

Caskin uses a custom casbin model with **two role-definition tables** (`g` and `g2`):

```ini
[request_definition]
r = sub, dom, obj, act

[policy_definition]
p = sub, dom, obj, act

[role_definition]
g  = _, _, _    # user → role inheritance, scoped to domain
g2 = _, _, _    # object → object inheritance (hierarchy), scoped to domain

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub, r.dom) && g2(r.obj, p.obj, r.dom) && r.dom == p.dom && r.act == p.act
  || g(r.sub, "superadmin", "superdomain")
```

Key design choices:
- `g(sub, role, domain)` — user-to-role membership, isolated per domain
- `g2(obj, parentObj, domain)` — object hierarchy (child inherits parent's policies)
- The `|| g(r.sub, "superadmin", "superdomain")` clause allows superadmins to bypass
  all domain checks by having a synthetic role in a synthetic "superdomain"

### Adapter

Caskin uses [`casbin/gorm-adapter`](https://github.com/casbin/gorm-adapter) to persist
policies in the same RDBMS as your domain metadata. The `casbin_rule` table is created
automatically by the adapter.

---

## Dictionary System

The `Dictionary` (`IDictionary`) is loaded from a TOML config file at startup. It defines:

| Section          | Purpose |
|------------------|---------|
| `feature`        | Named features of your application |
| `backend`        | API endpoints (path + method), used for backend permission checks |
| `frontend`       | UI items (name + type, e.g. `menu`/`button`), used for frontend visibility |
| `package`        | Logical groupings that bundle `backend` + `frontend` items under a `feature` key |
| `creator_object` | Objects created automatically when `ResetDomain` is called |
| `creator_role`   | Roles created automatically when `ResetDomain` is called |
| `creator_policy` | Policies seeded automatically when `ResetDomain` is called |

**`ResetDomain` vs `ResetFeature`:**
- `ResetDomain` — creates `creator_object`, `creator_role`, and `creator_policy` entries
  in the database for a specific domain. Used at domain bootstrap.
- `ResetFeature` — syncs `feature`, `backend`, and `frontend` definitions into the database
  for a specific domain. Call this after modifying the dictionary.

---

## Metadata Layer

`MetaDB` (implemented in `metadata_database.go` / `metadata_imp.go`) wraps GORM and provides
typed CRUD for caskin's core entities:

- Domain, User records
- Role records (scoped to domain + object)
- Object records (forming the permission tree)
- ObjectData records (any application entity protected by an Object)

All soft-delete patterns use GORM's `DeletedAt`. Recover methods (`RecoverUser`,
`RecoverObject`, etc.) un-delete records.

---

## Request Flow: A Permission Check

Here is what happens when your application calls `svc.GetObject(alice, domain, caskin.Read)`:

```
1. server.GetObject(alice, domain, Read)
   │
2. Calls MetaDB to fetch all Object records for the domain
   │
3. Calls IEnforcer.Enforce(alice, obj, domain, Read) for each object
   │   └─ casbin evaluates:
   │       g(alice, role, domain)   — is alice in a role?
   │       g2(obj, parent, domain)  — is obj in a hierarchy?
   │       p(role, obj, domain, act)— does a policy allow role+obj+domain+read?
   │
4. Returns only objects where alice has at least Read access
```

This "filter by permission" pattern is used throughout — `GetRole`, `GetPolicy`,
`GetCurrentBackend`, etc. all filter their results to what the calling user can see.

---

## Directory: Object Tree Traversal

Objects form a tree via `GetParentID()`. The `IDirectoryService` interface provides
tree-aware operations:

- `GetObjectDirectory` — returns an object and all its descendants (subtree)
- `GetObjectHierarchyLevel` — returns the depth of an object in the tree
- Object-data methods (`CreateObjectData`, `GetObjectData`, etc.) — CRUD for
  application entities bound to objects in the tree

Internally, tree traversal uses **iterative BFS** (breadth-first search) rather than
recursive queries, making it safe for deep hierarchies without stack overflow risk.

---

## Multi-Instance Sync (Redis Watcher)

When running caskin in a horizontally-scaled deployment (multiple instances sharing
one database), each instance has its own in-memory casbin enforcer. To keep them in
sync when one instance modifies a policy, caskin supports a Redis pub/sub watcher.

```go
svc, _ := caskin.New(&caskin.Options{
    Dictionary: &caskin.DictionaryOption{Dsn: "caskin.toml"},
    DB:         dbOption,
    Watcher: &caskin.WatcherOption{
        Type:     "redis",
        Address:  "localhost:6379",
        Password: "",
        Channel:  "/caskin",
    },
})
```

When a policy changes, the watcher publishes a notification on the configured channel.
All other instances receive it and reload their enforcer from the database.

If Redis is not available, set `AutoLoad` to a positive integer to enable periodic
policy reload (polling):

```go
Watcher: &caskin.WatcherOption{AutoLoad: 30} // reload every 30 seconds
```

---

## Package Layout

```
caskin/
├── casbin.go              # casbin model loading, WatcherOption, SetWatcher
├── constant.go            # DefaultSuperadminDomainName, DefaultSuperadminRoleName, Actions
├── dictionary.go          # IDictionary interface + NewDictionary
├── dictionary_adaptor.go  # TOML → IDictionary adapter
├── dictionary_model.go    # TOML struct definitions (Feature, Backend, Frontend, Package, ...)
├── doc.go                 # package-level godoc
├── error.go               # error types and sentinel errors
├── inheritance.go         # role/object hierarchy helpers (AddRoleG, BFS traversal)
├── metadata.go            # MetaDB interface
├── metadata_database.go   # GORM-backed MetaDB implementation
├── metadata_imp.go        # MetaDB helper implementations
├── object_deleter.go      # cascade object deletion logic
├── object_directory.go    # object tree directory helpers
├── object_updater.go      # object update/parent-reparent logic
├── options.go             # Options struct (DB, Dictionary, Watcher)
├── register.go            # Factory interface + Register[U,R,O,D]() generic function
├── schema.go              # Core interfaces: User, Role, Object, Domain, ObjectData, ...
├── schema_buildin.go      # Built-in policy/object types used internally
├── schema_private.go      # Internal interface fragments (idInterface, codeInterface, ...)
├── server.go              # server struct + New() constructor + Check/Filter helpers
├── server_check.go        # Permission validation helpers used across service methods
├── server_current.go      # ICurrentService implementation (SetCurrent)
├── server_directory.go    # IDirectoryService implementation
├── server_domain.go       # Domain CRUD
├── server_domain_reset.go # ResetDomain logic
├── server_feature.go      # IFeatureService — backend/frontend permission queries
├── server_feature_reset.go# ResetFeature logic
├── server_object.go       # Object CRUD
├── server_object_data.go  # ObjectData CRUD
├── server_policy.go       # Policy CRUD and modification
├── server_role.go         # Role CRUD
├── server_role_g.go       # Role inheritance (AddRoleG / RemoveRoleG)
├── server_superadmin.go   # Superadmin management
├── server_user.go         # User CRUD
├── server_user_domain.go  # User-domain queries
├── server_user_role.go    # User-role pair management
├── service.go             # IService and sub-interface definitions
│
├── configs/
│   ├── casbin_model.conf  # Embedded casbin RBAC+domain model
│   └── caskin.toml        # Default/example dictionary config
│
├── example/               # Reference implementations of User, Role, Object, Domain
│   ├── domain.go
│   ├── object.go
│   ├── object_data.go
│   ├── role.go
│   └── user.go
│
└── playground/            # Ready-to-use test environment (SQLite + bootstrapped data)
    └── playground.go
```

---

## Contributing

See [CONTRIBUTING.md](../CONTRIBUTING.md) for the contribution workflow,
coding conventions, and how to run the test suite.
