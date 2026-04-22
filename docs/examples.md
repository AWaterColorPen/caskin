# Usage Examples

This document covers common real-world scenarios for caskin. Each example is
self-contained and uses the types defined in the [`example/`](../example/)
package. All examples assume you have already set up a service via
[Getting Started](./getting-started.md).

---

## Table of Contents

- [Setup Helper](#setup-helper)
- [Scenario 1: Multi-Domain Management](#scenario-1-multi-domain-management)
- [Scenario 2: Role Inheritance](#scenario-2-role-inheritance)
- [Scenario 3: Permission Checks](#scenario-3-permission-checks)
- [Scenario 4: Frontend / Backend Permission Separation](#scenario-4-frontend--backend-permission-separation)

---

## Setup Helper

All examples below share this setup function that creates a working caskin
service backed by an in-memory SQLite database.

```go
package main

import (
    "os"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
    "gorm.io/gorm"
)

// newService creates an in-memory caskin service for demonstration.
// It mirrors the setup in playground/playground.go.
func newService() (caskin.IService, *gorm.DB) {
    dir, _ := os.MkdirTemp("", "caskin-example-*")
    dbOption := &caskin.DBOption{
        DSN:  dir + "/sqlite",
        Type: "sqlite",
    }
    db, _ := dbOption.NewDB()
    _ = db.AutoMigrate(
        &example.User{},
        &example.Role{},
        &example.Object{},
        &example.Domain{},
    )

    // Register must be called before New so the factory knows which
    // concrete types to instantiate.
    caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

    svc, _ := caskin.New(&caskin.Options{
        DB:         dbOption,
        Dictionary: &caskin.DictionaryOption{Dsn: "configs/caskin.toml"},
    })
    return svc, db
}
```

---

## Scenario 1: Multi-Domain Management

caskin is designed for **multi-tenant** systems. Each domain is an isolated
permission scope — users, roles, and policies in one domain do not affect
another.

### Create and bootstrap multiple domains

```go
func multiDomainExample() {
    svc, _ := newService()

    // --- Create two tenants (domains) ---
    engineering := &example.Domain{Name: "engineering"}
    marketing   := &example.Domain{Name: "marketing"}
    _ = svc.CreateDomain(engineering)
    _ = svc.CreateDomain(marketing)

    // Bootstrap each domain: creates the built-in admin/member roles and
    // the root object tree defined by the dictionary.
    _ = svc.ResetDomain(engineering)
    _ = svc.ResetDomain(marketing)
    _ = svc.ResetFeature(engineering)
    _ = svc.ResetFeature(marketing)

    // --- Create a shared superadmin ---
    superadmin := &example.User{Email: "admin@company.com"}
    _ = svc.CreateUser(superadmin)
    _ = svc.AddSuperadmin(superadmin)

    // --- Create domain-specific admins ---
    engAdmin := &example.User{Email: "eng-lead@company.com"}
    mktAdmin := &example.User{Email: "mkt-lead@company.com"}
    for _, u := range []caskin.User{engAdmin, mktAdmin} {
        _ = svc.CreateUser(u)
    }

    // Assign engAdmin as admin of the engineering domain only.
    engRoles, _ := svc.GetRole(superadmin, engineering)
    // engRoles[0] is the admin role created by ResetDomain.
    _ = svc.ModifyUserRolePerRole(superadmin, engineering, engRoles[0],
        []*caskin.UserRolePair{{User: engAdmin, Role: engRoles[0]}},
    )

    // Assign mktAdmin as admin of the marketing domain only.
    mktRoles, _ := svc.GetRole(superadmin, marketing)
    _ = svc.ModifyUserRolePerRole(superadmin, marketing, mktRoles[0],
        []*caskin.UserRolePair{{User: mktAdmin, Role: mktRoles[0]}},
    )

    // --- Verify isolation ---
    // engAdmin cannot see any roles in the marketing domain.
    rolesSeenByEngAdmin, _ := svc.GetRole(engAdmin, marketing)
    fmt.Println("eng admin sees marketing roles:", len(rolesSeenByEngAdmin)) // 0

    // The superadmin can see all domains.
    domains, _ := svc.GetDomain(superadmin)
    fmt.Println("total domains:", len(domains)) // 2
}
```

### Key points

| Behaviour | Detail |
|---|---|
| Domain isolation | Roles, objects, and policies are scoped per domain |
| Superadmin bypass | Superadmins can act across all domains |
| `ResetDomain` | Must be called after `CreateDomain` to initialise the built-in role/object tree |

---

## Scenario 2: Role Inheritance

caskin supports **role inheritance** (also called role hierarchies). A child
role automatically inherits all permissions of its parent role.

```go
func roleInheritanceExample() {
    svc, _ := newService()

    domain := &example.Domain{Name: "app"}
    _ = svc.CreateDomain(domain)
    _ = svc.ResetDomain(domain)
    _ = svc.ResetFeature(domain)

    superadmin := &example.User{Email: "root@example.com"}
    _ = svc.CreateUser(superadmin)
    _ = svc.AddSuperadmin(superadmin)

    // --- Build a three-level role hierarchy ---
    //   viewer  ←  editor  ←  owner
    // (viewer has fewest, owner has most permissions)

    viewer := &example.Role{Name: "viewer", DomainID: domain.GetID()}
    editor := &example.Role{Name: "editor", DomainID: domain.GetID()}
    owner  := &example.Role{Name: "owner",  DomainID: domain.GetID()}

    for _, r := range []caskin.ObjectData{viewer, editor, owner} {
        _ = svc.CreateObjectData(superadmin, domain, r, caskin.ObjectTypeRole)
    }

    // editor inherits from viewer  (editor >= viewer)
    _ = svc.AddRoleG(superadmin, domain, editor, viewer)

    // owner inherits from editor   (owner >= editor >= viewer)
    _ = svc.AddRoleG(superadmin, domain, owner, editor)

    // --- Grant baseline permissions to viewer on a resource object ---
    objects, _ := svc.GetObject(superadmin, domain, caskin.Read)
    // Use the first non-root object as the demo resource.
    var resource caskin.Object
    for _, o := range objects {
        if o.GetParentID() != 0 {
            resource = o
            break
        }
    }

    // viewer: read-only
    _ = svc.ModifyPolicyPerRole(superadmin, domain, viewer,
        []*caskin.Policy{{Role: viewer, Object: resource, Domain: domain, Action: caskin.Read}},
    )
    // editor: also write
    _ = svc.ModifyPolicyPerRole(superadmin, domain, editor,
        []*caskin.Policy{{Role: editor, Object: resource, Domain: domain, Action: caskin.Write}},
    )
    // owner: also manage
    _ = svc.ModifyPolicyPerRole(superadmin, domain, owner,
        []*caskin.Policy{{Role: owner, Object: resource, Domain: domain, Action: caskin.Manage}},
    )

    // --- Assign users to roles ---
    alice := &example.User{Email: "alice@example.com"}
    bob   := &example.User{Email: "bob@example.com"}
    carol := &example.User{Email: "carol@example.com"}
    for _, u := range []caskin.User{alice, bob, carol} {
        _ = svc.CreateUser(u)
    }

    _ = svc.ModifyUserRolePerRole(superadmin, domain, viewer,
        []*caskin.UserRolePair{{User: alice, Role: viewer}})
    _ = svc.ModifyUserRolePerRole(superadmin, domain, editor,
        []*caskin.UserRolePair{{User: bob, Role: editor}})
    _ = svc.ModifyUserRolePerRole(superadmin, domain, owner,
        []*caskin.UserRolePair{{User: carol, Role: owner}})

    // --- Verify that Bob (editor) inherits viewer permissions ---
    // Bob should be able to read (inherited) and write (direct).
    // IService exposes CheckObject which returns nil on success.
    canRead  := svc.CheckObject(bob, domain, resource, caskin.Read) == nil
    canWrite := svc.CheckObject(bob, domain, resource, caskin.Write) == nil
    fmt.Println("bob can read:", canRead)   // true (inherited from viewer)
    fmt.Println("bob can write:", canWrite) // true (direct on editor)

    // Alice (viewer) cannot write.
    aliceWrite := svc.CheckObject(alice, domain, resource, caskin.Write) == nil
    fmt.Println("alice can write:", aliceWrite) // false

    // --- Remove the editor → viewer link at runtime ---
    _ = svc.RemoveRoleG(superadmin, domain, editor, viewer)
    // Now Bob no longer inherits viewer's read permission via that path.
}
```

### Key points

| API | Description |
|---|---|
| `AddRoleG(user, domain, from, to)` | `from` inherits all permissions of `to` |
| `RemoveRoleG(user, domain, from, to)` | Remove the inheritance link at runtime |
| Transitivity | Inheritance is transitive: owner → editor → viewer |

---

## Scenario 3: Permission Checks

caskin exposes two layers for checking permissions:

1. **`IService.CheckObject`** — service-level check; returns a typed error and
   respects the caller's own permission scope.
2. **`ICurrentService.Check*WithCurrent`** — middleware pattern; binds user and
   domain once via `SetCurrent` then checks without re-passing them.

```go
func permissionCheckExample() {
    svc, _ := newService()

    domain := &example.Domain{Name: "wiki"}
    _ = svc.CreateDomain(domain)
    _ = svc.ResetDomain(domain)
    _ = svc.ResetFeature(domain)

    superadmin := &example.User{Email: "root@example.com"}
    _ = svc.CreateUser(superadmin)
    _ = svc.AddSuperadmin(superadmin)

    editor := &example.Role{Name: "editor", DomainID: domain.GetID()}
    _ = svc.CreateObjectData(superadmin, domain, editor, caskin.ObjectTypeRole)

    alice := &example.User{Email: "alice@example.com"}
    _ = svc.CreateUser(alice)
    _ = svc.ModifyUserRolePerRole(superadmin, domain, editor,
        []*caskin.UserRolePair{{User: alice, Role: editor}})

    // Grab a real Object to check against.
    objects, _ := svc.GetObject(superadmin, domain, caskin.Read)
    var article caskin.Object
    for _, o := range objects {
        if o.GetParentID() != 0 {
            article = o
            break
        }
    }

    // Grant editor the write permission on article.
    _ = svc.ModifyPolicyPerRole(superadmin, domain, editor,
        []*caskin.Policy{{Role: editor, Object: article, Domain: domain, Action: caskin.Write}},
    )

    // --- Method 1: service-level boolean check ---
    ok := svc.CheckObject(alice, domain, article, caskin.Write) == nil
    fmt.Println("alice can write:", ok) // true

    // --- Method 2: service-level check (respects caller's own permissions) ---
    err := svc.CheckObject(alice, domain, article, caskin.Write)
    fmt.Println("alice write error:", err) // <nil>

    err = svc.CheckObject(alice, domain, article, caskin.Manage)
    fmt.Println("alice manage error:", err) // "no manage permission"

    // --- Using ICurrentService for middleware-style checks ---
    // Bind the current user + domain once (e.g. in an HTTP middleware) and
    // then call the Check* methods without passing user/domain on every call.
    current := svc.SetCurrent(alice, domain)
    err = current.CheckModifyObjectDataWithCurrent(editor)
    fmt.Println("alice modify editor (current):", err) // <nil> — alice is editor
}
```

### Choosing the right check

| Scenario | Recommended API |
|---|---|
| HTTP middleware / auth gate | `ICurrentService.Check*WithCurrent` after `SetCurrent` |
| Business logic, needs typed error | `IService.CheckObject` / `CheckObjectData` |
| Testing policy directly (with concrete `*server`) | `caskin.Check(enforcer, ...)` (package-level, not on `IService`) |

---

## Scenario 4: Frontend / Backend Permission Separation

A common pattern is to expose **different object trees** to frontend (UI
buttons/pages) and backend (API endpoints). caskin models this naturally
because each `Object` can have a custom type, and you can organise objects into
separate sub-trees.

```go
// FrontendObject represents a UI element (e.g. a menu item or button).
type FrontendObject struct {
    example.Object
    // Extra fields meaningful to the frontend, e.g. RouteKey or ComponentName.
    RouteKey string `gorm:"column:route_key"`
}

// BackendObject represents an API endpoint permission.
type BackendObject struct {
    example.Object
    // Extra fields meaningful to the backend, e.g. HTTP method and path.
    Method string `gorm:"column:method"`
    Path   string `gorm:"column:path"`
}

func frontendBackendSeparationExample() {
    svc, _ := newService()

    domain := &example.Domain{Name: "saas-app"}
    _ = svc.CreateDomain(domain)
    _ = svc.ResetDomain(domain)
    _ = svc.ResetFeature(domain)

    superadmin := &example.User{Email: "root@example.com"}
    _ = svc.CreateUser(superadmin)
    _ = svc.AddSuperadmin(superadmin)

    // Retrieve the root objects created by ResetFeature.
    // By convention, the first root-level object is where you hang your
    // custom sub-trees.
    rootObjects, _ := svc.GetObject(superadmin, domain, caskin.Read)
    var root caskin.Object
    for _, o := range rootObjects {
        if o.GetParentID() == 0 {
            root = o
            break
        }
    }

    // --- Create two top-level "namespace" objects ---
    // All frontend objects live under "ui-root".
    // All backend objects live under "api-root".
    uiRoot := &example.Object{
        Name:     "ui-root",
        Type:     "ui",
        ParentID: root.GetID(),
        DomainID: domain.GetID(),
    }
    apiRoot := &example.Object{
        Name:     "api-root",
        Type:     "api",
        ParentID: root.GetID(),
        DomainID: domain.GetID(),
    }
    _ = svc.CreateObject(superadmin, domain, uiRoot)
    _ = svc.CreateObject(superadmin, domain, apiRoot)

    // --- Add child objects for specific UI pages and API endpoints ---
    dashboardPage := &example.Object{
        Name: "/dashboard", Type: "ui", ParentID: uiRoot.GetID(), DomainID: domain.GetID(),
    }
    settingsPage := &example.Object{
        Name: "/settings", Type: "ui", ParentID: uiRoot.GetID(), DomainID: domain.GetID(),
    }
    apiUsers := &example.Object{
        Name: "GET /api/users", Type: "api", ParentID: apiRoot.GetID(), DomainID: domain.GetID(),
    }
    apiUsersWrite := &example.Object{
        Name: "POST /api/users", Type: "api", ParentID: apiRoot.GetID(), DomainID: domain.GetID(),
    }
    for _, o := range []caskin.Object{dashboardPage, settingsPage, apiUsers, apiUsersWrite} {
        _ = svc.CreateObject(superadmin, domain, o)
    }

    // --- Create roles with different permission scopes ---
    readonly := &example.Role{Name: "readonly", DomainID: domain.GetID()}
    fullAccess := &example.Role{Name: "full-access", DomainID: domain.GetID()}
    for _, r := range []caskin.ObjectData{readonly, fullAccess} {
        _ = svc.CreateObjectData(superadmin, domain, r, caskin.ObjectTypeRole)
    }

    // readonly: can see the dashboard and call GET /api/users
    _ = svc.ModifyPolicyPerRole(superadmin, domain, readonly, []*caskin.Policy{
        {Role: readonly, Object: dashboardPage, Domain: domain, Action: caskin.Read},
        {Role: readonly, Object: apiUsers, Domain: domain, Action: caskin.Read},
    })

    // full-access: everything including settings and write APIs
    _ = svc.ModifyPolicyPerRole(superadmin, domain, fullAccess, []*caskin.Policy{
        {Role: fullAccess, Object: dashboardPage, Domain: domain, Action: caskin.Read},
        {Role: fullAccess, Object: settingsPage, Domain: domain, Action: caskin.Manage},
        {Role: fullAccess, Object: apiUsers, Domain: domain, Action: caskin.Read},
        {Role: fullAccess, Object: apiUsersWrite, Domain: domain, Action: caskin.Write},
    })

    // --- Assign users ---
    alice := &example.User{Email: "alice@example.com"} // read-only viewer
    bob   := &example.User{Email: "bob@example.com"}   // full-access admin
    for _, u := range []caskin.User{alice, bob} {
        _ = svc.CreateUser(u)
    }
    _ = svc.ModifyUserRolePerRole(superadmin, domain, readonly,
        []*caskin.UserRolePair{{User: alice, Role: readonly}})
    _ = svc.ModifyUserRolePerRole(superadmin, domain, fullAccess,
        []*caskin.UserRolePair{{User: bob, Role: fullAccess}})

    // --- Simulate what the frontend queries at login ---
    // "Which UI pages can Alice see?"
    aliceUIObjects := filterByType(mustGetObjects(svc, alice, domain, caskin.Read), "ui")
    fmt.Println("alice UI pages:", names(aliceUIObjects))
    // [/dashboard]  — /settings is not in her policy

    // --- Simulate what an API gateway checks per request ---
    // "Can Alice call POST /api/users?"
    canPost := svc.CheckObject(alice, domain, apiUsersWrite, caskin.Write) == nil
    fmt.Println("alice can POST /api/users:", canPost) // false

    canPost = svc.CheckObject(bob, domain, apiUsersWrite, caskin.Write) == nil
    fmt.Println("bob can POST /api/users:", canPost) // true
}

// --- helpers used above ---

func mustGetObjects(svc caskin.IService, u caskin.User, d caskin.Domain, a caskin.Action) []caskin.Object {
    objs, _ := svc.GetObject(u, d, a)
    return objs
}

func filterByType(objs []caskin.Object, ty caskin.ObjectType) []caskin.Object {
    var out []caskin.Object
    for _, o := range objs {
        if o.GetObjectType() == ty {
            out = append(out, o)
        }
    }
    return out
}

func names(objs []caskin.Object) []string {
    out := make([]string, len(objs))
    for i, o := range objs {
        out[i] = o.GetName()
    }
    return out
}
```

### Recommended middleware pattern

```go
// HTTP middleware example (framework-agnostic pseudocode).
func PermissionMiddleware(svc caskin.IService) Middleware {
    return func(ctx Context, next Handler) {
        user   := ctx.CurrentUser().(caskin.User)
        domain := ctx.CurrentDomain().(caskin.Domain)
        object := lookupObjectForRoute(ctx.Route()).(caskin.Object)

        if svc.CheckObject(user, domain, object, caskin.Read) != nil {
            ctx.Abort(http.StatusForbidden)
            return
        }
        next(ctx)
    }
}
```

---

## See Also

- [Getting Started](./getting-started.md) — step-by-step setup from scratch
- [API Reference](./api-reference.md) — full type and method documentation
- [`playground/`](../playground/) — runnable integration test environment
- [`example/`](../example/) — sample implementations of the four core types
