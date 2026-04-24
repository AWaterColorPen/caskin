# Common Use Cases

This document walks through the most common authorization patterns in caskin.
All examples build on the [Getting Started](./getting-started.md) guide and use the
`example` and `playground` packages.

## Table of Contents

- [Multi-Domain Management](#multi-domain-management)
- [Role Hierarchy (Inheritance)](#role-hierarchy-inheritance)
- [Permission Checks](#permission-checks)
- [Frontend / Backend Permission Separation](#frontend--backend-permission-separation)

---

## Multi-Domain Management

Caskin is designed for multi-tenant systems where each domain (tenant/organization)
has completely isolated RBAC — the same user can have different roles in different domains.

### Creating Multiple Domains

```go
package main

import (
    "log"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func main() {
    // Register concrete types (do this once, before caskin.New)
    caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

    dbOption := &caskin.DBOption{DSN: "./multi-domain.db", Type: "sqlite"}
    db, _ := dbOption.NewDB()
    db.AutoMigrate(&example.User{}, &example.Role{}, &example.Object{}, &example.Domain{})

    svc, _ := caskin.New(&caskin.Options{
        Dictionary: &caskin.DictionaryOption{Dsn: "caskin.toml"},
        DB:         dbOption,
    })

    // Create two independent domains
    domainA := &example.Domain{Name: "org-a"}
    domainB := &example.Domain{Name: "org-b"}
    svc.CreateDomain(domainA)
    svc.ResetDomain(domainA) // creates initial roles/objects from caskin.toml
    svc.ResetFeature(domainA)

    svc.CreateDomain(domainB)
    svc.ResetDomain(domainB)
    svc.ResetFeature(domainB)

    // Create a superadmin and a regular user
    superadmin := &example.User{Email: "super@example.com"}
    alice := &example.User{Email: "alice@example.com"}
    svc.CreateUser(superadmin)
    svc.CreateUser(alice)
    svc.AddSuperadmin(superadmin)

    // Assign alice as admin in org-a, member in org-b
    rolesA, _ := svc.GetRole(superadmin, domainA) // [admin, member] from caskin.toml
    rolesB, _ := svc.GetRole(superadmin, domainB)

    svc.ModifyUserRolePerRole(superadmin, domainA, rolesA[0], []*caskin.UserRolePair{
        {User: alice, Role: rolesA[0]},
    })
    svc.ModifyUserRolePerRole(superadmin, domainB, rolesB[1], []*caskin.UserRolePair{
        {User: alice, Role: rolesB[1]},
    })

    // Query alice's roles — different per domain
    pairsA, _ := svc.GetUserRoleByUser(superadmin, domainA, alice)
    pairsB, _ := svc.GetUserRoleByUser(superadmin, domainB, alice)
    log.Printf("alice in org-a: %d role(s)", len(pairsA)) // 1 (admin)
    log.Printf("alice in org-b: %d role(s)", len(pairsB)) // 1 (member)
}
```

### Listing Domains for a User

```go
// Find all domains alice belongs to
domains, _ := svc.GetDomainByUser(alice)
for _, d := range domains {
    log.Printf("alice is in domain: %v", d)
}
```

### Listing Users in a Domain

```go
// List all users in org-a
users, _ := svc.GetUserByDomain(domainA)
for _, u := range users {
    log.Printf("user in org-a: %v", u)
}
```

---

## Role Hierarchy (Inheritance)

Caskin supports role inheritance within a domain: a child role inherits all permissions
of its parent role. This is useful for building permission tiers (e.g., `viewer` → `editor` → `owner`).

### Setting Up Role Inheritance

```go
package main

import (
    "log"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
    "github.com/awatercolorpen/caskin/playground"
)

func main() {
    playground.DictionaryDsn = "configs/caskin.toml"

    // Use playground for a quick fully-bootstrapped environment
    stage, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
    if err != nil {
        log.Fatal(err)
    }
    svc := stage.Service
    superadmin := stage.Superadmin
    domain := stage.Domain

    // Create two extra roles: "editor" and "viewer"
    editor := &example.Role{Name: "editor"}
    viewer := &example.Role{Name: "viewer"}
    svc.CreateRole(superadmin, domain, editor)
    svc.CreateRole(superadmin, domain, viewer)

    // editor inherits from viewer (editor gets all of viewer's permissions, plus more)
    svc.AddRoleG(superadmin, domain, editor, viewer)

    // Now assign alice as viewer and bob as editor
    alice := &example.User{Email: "alice@example.com"}
    bob := &example.User{Email: "bob@example.com"}
    svc.CreateUser(alice)
    svc.CreateUser(bob)

    svc.ModifyUserRolePerRole(superadmin, domain, viewer, []*caskin.UserRolePair{
        {User: alice, Role: viewer},
    })
    svc.ModifyUserRolePerRole(superadmin, domain, editor, []*caskin.UserRolePair{
        {User: bob, Role: editor},
    })

    // List role inheritance
    pairs, _ := svc.GetUserRole(superadmin, domain)
    for _, p := range pairs {
        log.Printf("user %v → role %v", p.User, p.Role)
    }

    // Remove inheritance
    svc.RemoveRoleG(superadmin, domain, editor, viewer)
}
```

> **Note:** Role inheritance is scoped to a single domain. A role in `org-a` cannot
> inherit from a role in `org-b`.

---

## Permission Checks

Caskin provides two levels of permission checks:

1. **Service-level** — high-level checks via `IService` methods (e.g., `GetObject` automatically
   filters by what the caller can access).
2. **Direct enforcement** — low-level `caskin.Check` helper for checking one specific resource.

### Service-Level Filtering

The simplest way: use service methods that require an operator (`User`) and domain.
If the caller lacks permission, the method returns an error or an empty result.

```go
// alice can only see objects she has at least "read" permission on
objects, err := svc.GetObject(alice, domain, caskin.Read)
if err != nil {
    log.Printf("error: %v", err)
}
for _, obj := range objects {
    log.Printf("alice can see: %v", obj)
}
```

### Direct Permission Check with `caskin.Check`

Use `caskin.Check` when you have a concrete enforcer and need to test a single resource:

```go
// Obtain the low-level enforcer — requires holding the concrete *server type.
// In most application code, prefer service-level methods instead.
//
// If you need direct enforcement in tests, use the playground's Service and
// type-assert to access internal helpers, or design your app to use service methods.

// Example: filter a slice of objects for a given user+domain+action
enforcer := ... // obtained from internal server (see architecture docs)
allowed := caskin.Filter(enforcer, alice, domain, caskin.Read, objects)
log.Printf("alice has read access to %d object(s)", len(allowed))
```

### Using `ICurrentService` for Scoped Checks

`ICurrentService` binds a specific operator and domain for the lifetime of a request:

```go
// Bind alice + domain for this request context
current := svc.SetCurrent(alice, domain)

// All calls on current are scoped to alice+domain — no need to pass them each time
myRoles, _ := current.GetCurrentRole()
myObjects, _ := current.GetCurrentObject()
myPolicies, _ := current.GetCurrentPolicy()

log.Printf("alice's roles in domain: %d", len(myRoles))
log.Printf("alice's visible objects: %d", len(myObjects))
```

### Checking Policies by Role

```go
// List all policies for a specific role
roles, _ := svc.GetRole(superadmin, domain)
adminRole := roles[0]

policies, _ := svc.GetPolicyByRole(superadmin, domain, adminRole)
for _, p := range policies {
    log.Printf("admin policy: role=%v object=%v action=%v", p.Role, p.Object, p.Action)
}
```

### Modifying Policies

```go
// Grant the "member" role read access to an object
objects, _ := svc.GetObject(superadmin, domain, caskin.Manage)
articleObject := objects[0] // find the object you want

memberRole := roles[1] // "member" role

newPolicies := []*caskin.Policy{
    {Role: memberRole, Object: articleObject, Action: caskin.Read},
}
svc.ModifyPolicyPerRole(superadmin, domain, memberRole, newPolicies)
```

---

## Frontend / Backend Permission Separation

Caskin supports separate permission models for **backend APIs** (server-side access control)
and **frontend UI** (visibility of menus, buttons, pages). Both are defined in the dictionary
config (`caskin.toml`) and managed through the `IFeatureService` interface.

### Dictionary Setup

Define backend and frontend entries in your `caskin.toml`:

```toml
feature = [
    {name = "article"},
    {name = "user-management"},
]

# Backend: API-level permissions (method + path)
backend = [
    {path = "api/articles",      method = "GET"},
    {path = "api/articles",      method = "POST"},
    {path = "api/articles/:id",  method = "PUT"},
    {path = "api/users",         method = "GET"},
]

# Frontend: UI-level permissions (name + type)
frontend = [
    {name = "article-list",   type = "menu"},
    {name = "article-create", type = "button"},
    {name = "user-list",      type = "menu"},
]

# Packages group backend + frontend into a logical permission unit
package = [
    {
        key      = "article-read",
        backend  = [["api/articles", "GET"]],
        frontend = [["article-list", "menu"]],
    },
    {
        key      = "article-write",
        backend  = [["api/articles", "POST"], ["api/articles/:id", "PUT"]],
        frontend = [["article-create", "button"]],
    },
    {
        key      = "user-management",
        backend  = [["api/users", "GET"]],
        frontend = [["user-list", "menu"]],
    },
]
```

### Querying Backend and Frontend Permissions

```go
// After ResetFeature, query what backend APIs a user can access
current := svc.SetCurrent(alice, domain)

backends, _ := current.GetCurrentBackend()
for _, b := range backends {
    log.Printf("alice can call: %s %s", b.GetMethod(), b.GetPath())
}

frontends, _ := current.GetCurrentFrontend()
for _, f := range frontends {
    log.Printf("alice can see: %s (%s)", f.GetName(), f.GetFrontendType())
}
```

### Using Backend Permissions as a Middleware

```go
// In an HTTP middleware, check if the current user can call this endpoint:
func AuthMiddleware(svc caskin.IService, user caskin.User, domain caskin.Domain) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            current := svc.SetCurrent(user, domain)
            backends, err := current.GetCurrentBackend()
            if err != nil {
                http.Error(w, "unauthorized", http.StatusUnauthorized)
                return
            }
            allowed := false
            for _, b := range backends {
                if b.GetPath() == r.URL.Path && b.GetMethod() == r.Method {
                    allowed = true
                    break
                }
            }
            if !allowed {
                http.Error(w, "forbidden", http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

### Resetting Features After Dictionary Changes

When you update `caskin.toml`, call `ResetFeature` to sync the changes into the database:

```go
// Resync the feature/backend/frontend definitions for all domains
domains, _ := svc.GetDomain()
for _, d := range domains {
    if err := svc.ResetFeature(d); err != nil {
        log.Printf("failed to reset feature for domain %v: %v", d, err)
    }
}
```

---

## Next Steps

- [API Reference](./api-reference.md) — full method reference with parameters and return values
- [Configuration](./configuration.md) — all config options for DB, dictionary, and Redis watcher
- [Architecture](./architecture.md) — how caskin is structured internally (for contributors)
