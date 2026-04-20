# Getting Started with Caskin

This guide walks you through setting up caskin from scratch — from installing the library
to running real permission checks. All code examples are self-contained and runnable.

## Table of Contents

- [Overview](#overview)
- [Core Concepts](#core-concepts)
- [Step 1: Install](#step-1-install)
- [Step 2: Implement the Four Core Types](#step-2-implement-the-four-core-types)
- [Step 3: Create a Dictionary Config](#step-3-create-a-dictionary-config)
- [Step 4: Initialize the Service](#step-4-initialize-the-service)
- [Step 5: Bootstrap a Domain](#step-5-bootstrap-a-domain)
- [Step 6: Manage Users and Roles](#step-6-manage-users-and-roles)
- [Step 7: Permission Checks](#step-7-permission-checks)
- [Step 8: Use the Playground for Testing](#step-8-use-the-playground-for-testing)
- [Complete Example](#complete-example)
- [Next Steps](#next-steps)

---

## Overview

Caskin is a **multi-domain RBAC** (Role-Based Access Control) library for Go, built on top of
[casbin](https://github.com/casbin/casbin). It manages _authorization business logic_ — creating
and managing users, roles, objects, and domains with fine-grained permission control across
multiple tenants.

**When to use caskin:**

- You need per-tenant/per-organization permission isolation (multi-domain)
- You need hierarchical roles (role A inherits from role B)
- You want to control both backend API access and frontend UI visibility
- You need a dictionary-driven setup where permissions are defined in config files

---

## Core Concepts

| Concept    | Description |
|------------|-------------|
| **User**   | A real person who can be granted roles in domains |
| **Role**   | A named group of permissions within a domain; can be inherited |
| **Object** | A resource or permission node (e.g., an API endpoint, UI menu, or data resource) |
| **Domain** | An isolated tenant/organization; RBAC is fully isolated per domain |
| **Policy** | A binding between a Role and an Object with an allowed Action (read/write/manage) |

A user can belong to multiple domains, with different roles in each.

---

## Step 1: Install

```bash
go get github.com/awatercolorpen/caskin
```

Caskin requires Go 1.21+.

---

## Step 2: Implement the Four Core Types

Caskin works with your own domain models — you define the structs, caskin provides the interfaces.
Each type must implement a small interface: `GetID()`/`SetID()`, `Encode()`/`Decode()`, and a few
domain/object-specific methods.

The easiest way to start is to copy the `example` package implementations:

```go
// user.go
package myapp

import (
    "fmt"
    "time"
    "gorm.io/gorm"
)

type User struct {
    ID        uint64         `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Email     string         `gorm:"unique"`
}

func (u *User) GetID() uint64       { return u.ID }
func (u *User) SetID(id uint64)     { u.ID = id }
func (u *User) Encode() string      { return fmt.Sprintf("user_%d", u.ID) }
func (u *User) Decode(s string) error {
    _, err := fmt.Sscanf(s, "user_%d", &u.ID)
    return err
}
```

```go
// role.go
type Role struct {
    ID       uint64         `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Name     string         `gorm:"uniqueIndex:idx_role"`
    DomainID uint64         `gorm:"uniqueIndex:idx_role"`
    ObjectID uint64         `gorm:"uniqueIndex:idx_role"`
}

func (r *Role) GetID() uint64          { return r.ID }
func (r *Role) SetID(id uint64)        { r.ID = id }
func (r *Role) Encode() string         { return fmt.Sprintf("role_%d", r.ID) }
func (r *Role) Decode(s string) error  { _, err := fmt.Sscanf(s, "role_%d", &r.ID); return err }
func (r *Role) GetObjectID() uint64    { return r.ObjectID }
func (r *Role) SetObjectID(id uint64)  { r.ObjectID = id }
func (r *Role) GetDomainID() uint64    { return r.DomainID }
func (r *Role) SetDomainID(id uint64)  { r.DomainID = id }
```

```go
// object.go
type Object struct {
    ID       uint64         `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Name     string         `gorm:"uniqueIndex:idx_object"`
    Type     string
    DomainID uint64         `gorm:"uniqueIndex:idx_object"`
    ParentID uint64         `gorm:"uniqueIndex:idx_object"`
}

func (o *Object) GetID() uint64           { return o.ID }
func (o *Object) SetID(id uint64)         { o.ID = id }
func (o *Object) Encode() string          { return fmt.Sprintf("object_%d", o.ID) }
func (o *Object) Decode(s string) error   { _, err := fmt.Sscanf(s, "object_%d", &o.ID); return err }
func (o *Object) GetParentID() uint64     { return o.ParentID }
func (o *Object) SetParentID(id uint64)   { o.ParentID = id }
func (o *Object) GetDomainID() uint64     { return o.DomainID }
func (o *Object) SetDomainID(id uint64)   { o.DomainID = id }
func (o *Object) GetObjectType() string   { return o.Type }
```

```go
// domain.go
type Domain struct {
    ID        uint64         `gorm:"primaryKey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
    Name      string         `gorm:"unique"`
}

func (d *Domain) GetID() uint64          { return d.ID }
func (d *Domain) SetID(id uint64)        { d.ID = id }
func (d *Domain) Encode() string         { return fmt.Sprintf("domain_%d", d.ID) }
func (d *Domain) Decode(s string) error  { _, err := fmt.Sscanf(s, "domain_%d", &d.ID); return err }
```

> **Tip:** You can import `github.com/awatercolorpen/caskin/example` directly if you want
> pre-built implementations for prototyping or tests.

---

## Step 3: Create a Dictionary Config

The dictionary config defines what **features** your system exposes, what **objects** are created
when a domain is bootstrapped, and what initial **roles and policies** each domain starts with.

Create `caskin.toml` in your project root:

```toml
# Features your system exposes (logical groupings)
feature = [
    {name = "article"},
    {name = "user-management"},
]

# Backend API endpoints (path + HTTP method)
backend = [
    {path = "api/articles",     method = "GET"},
    {path = "api/articles",     method = "POST"},
    {path = "api/articles/:id", method = "PUT"},
    {path = "api/articles/:id", method = "DELETE"},
    {path = "api/users",        method = "GET"},
    {path = "api/users",        method = "POST"},
]

# Frontend menu/UI items (name + type)
frontend = [
    {name = "articles",         type = "menu"},
    {name = "article-editor",   type = "button"},
    {name = "users",            type = "menu"},
]

# Packages bundle backend + frontend into a logical permission unit
package = [
    {
        key      = "article-read",
        backend  = [["api/articles", "GET"]],
        frontend = [["articles", "menu"]],
    },
    {
        key      = "article-write",
        backend  = [["api/articles", "POST"], ["api/articles/:id", "PUT"], ["api/articles/:id", "DELETE"]],
        frontend = [["article-editor", "button"]],
    },
    {
        key      = "user-management",
        backend  = [["api/users", "GET"], ["api/users", "POST"]],
        frontend = [["users", "menu"]],
    },
]

# Objects created automatically when ResetDomain() is called
creator_object = [
    {name = "role_root", type = "role"},
]

# Roles created automatically when ResetDomain() is called
creator_role = [
    {name = "admin"},
    {name = "editor"},
    {name = "viewer"},
]

# Policies assigned to those initial roles
creator_policy = [
    # admin: full access to roles and all features
    {role = "admin",  object = "role_root",                                   action = ["read", "write", "manage"]},
    {role = "admin",  object = "github.com/awatercolorpen/caskin::article",   action = ["read"]},
    # editor: can write articles
    {role = "editor", object = "role_root",                                   action = ["read"]},
    {role = "editor", object = "github.com/awatercolorpen/caskin::article",   action = ["read"]},
    # viewer: read-only
    {role = "viewer", object = "role_root",                                   action = ["read"]},
]
```

**Dictionary config cheat sheet:**

| Field | Purpose |
|-------|---------|
| `feature` | Logical groupings shown in the feature list |
| `backend` | API endpoints for `AuthBackend()` checks |
| `frontend` | UI items returned by `AuthFrontend()` |
| `package` | Bundles backend + frontend into a named permission package |
| `creator_object` | Objects auto-created on `ResetDomain()` |
| `creator_role` | Roles auto-created on `ResetDomain()` |
| `creator_policy` | Policies auto-assigned to those initial roles |

---

## Step 4: Initialize the Service

```go
package main

import (
    "log"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
    "gorm.io/gorm"
)

func main() {
    // Register your concrete types — call this once at startup, before caskin.New()
    caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

    // Configure the database (SQLite for local dev)
    dbOption := &caskin.DBOption{
        DSN:  "./myapp.db",
        Type: "sqlite",
    }

    // Auto-migrate your tables
    db, err := dbOption.NewDB()
    if err != nil {
        log.Fatal("open db:", err)
    }
    if err := db.AutoMigrate(
        &example.User{},
        &example.Role{},
        &example.Object{},
        &example.Domain{},
    ); err != nil {
        log.Fatal("migrate:", err)
    }

    // Build the caskin service
    svc, err := caskin.New(&caskin.Options{
        Dictionary: &caskin.DictionaryOption{Dsn: "caskin.toml"},
        DB:         dbOption,
    })
    if err != nil {
        log.Fatal("caskin.New:", err)
    }

    _ = svc // use svc below
}
```

> **Production tip:** Swap `"sqlite"` for `"mysql"` or `"postgres"` by changing `DBOption.Type`
> and providing the appropriate DSN. No other code changes needed.

---

## Step 5: Bootstrap a Domain

Before any permissions can be checked, you need to:
1. Create a domain
2. Reset it (creates the initial objects and roles from `caskin.toml`)
3. Reset its features (registers backend/frontend definitions)
4. Create a superadmin user

```go
// Create and bootstrap a domain
domain := &example.Domain{Name: "acme-corp"}
if err := svc.CreateDomain(domain); err != nil {
    log.Fatal(err)
}
if err := svc.ResetDomain(domain); err != nil {   // creates role_root, admin, editor, viewer roles
    log.Fatal(err)
}
if err := svc.ResetFeature(domain); err != nil {  // registers feature/backend/frontend objects
    log.Fatal(err)
}

// Create a superadmin — superadmins bypass domain permission checks
superadmin := &example.User{Email: "superadmin@example.com"}
if err := svc.CreateUser(superadmin); err != nil {
    log.Fatal(err)
}
if err := svc.AddSuperadmin(superadmin); err != nil {
    log.Fatal(err)
}

log.Printf("Domain %q bootstrapped (id=%d)", domain.Name, domain.ID)
log.Printf("Superadmin created (id=%d)", superadmin.ID)
```

---

## Step 6: Manage Users and Roles

### Create regular users

```go
alice := &example.User{Email: "alice@example.com"}
bob   := &example.User{Email: "bob@example.com"}

for _, u := range []*example.User{alice, bob} {
    if err := svc.CreateUser(u); err != nil {
        log.Fatal(err)
    }
}
```

### List roles in a domain

```go
roles, err := svc.GetRole(superadmin, domain)
if err != nil {
    log.Fatal(err)
}
for _, r := range roles {
    log.Printf("role: %+v", r)
}
// roles[0] = admin, roles[1] = editor, roles[2] = viewer (order matches creator_role)
```

### Assign roles to users

```go
adminRole  := roles[0] // admin
editorRole := roles[1] // editor

// Give alice the admin role
if err := svc.ModifyUserRolePerUser(superadmin, domain, alice, []*caskin.UserRolePair{
    {User: alice, Role: adminRole},
}); err != nil {
    log.Fatal(err)
}

// Give bob the editor role
if err := svc.ModifyUserRolePerUser(superadmin, domain, bob, []*caskin.UserRolePair{
    {User: bob, Role: editorRole},
}); err != nil {
    log.Fatal(err)
}
```

### Inspect user-role assignments

```go
pairs, err := svc.GetUserRole(superadmin, domain)
if err != nil {
    log.Fatal(err)
}
for _, p := range pairs {
    log.Printf("user %v → role %v", p.User, p.Role)
}
```

---

## Step 7: Permission Checks

### Check backend API access

```go
// Does bob have access to POST api/articles?
err := svc.AuthBackend(bob, domain, &caskin.Backend{
    Path:   "api/articles",
    Method: "POST",
})
if err != nil {
    log.Printf("bob cannot POST articles: %v", err)
} else {
    log.Println("bob can POST articles ✓")
}
```

### Get frontend UI items for a user

```go
// Which frontend items can bob see?
items := svc.AuthFrontend(bob, domain)
for _, item := range items {
    log.Printf("frontend item: %s (%s)", item.Name, item.Type)
}
```

### Get features visible to a user

```go
features, err := svc.GetFeature(bob, domain)
if err != nil {
    log.Fatal(err)
}
for _, f := range features {
    log.Printf("feature: %s", f.Name)
}
```

### Use SetCurrent for scoped operations

When you always operate in the same user+domain context (e.g., within an HTTP request handler),
use `SetCurrent` to avoid passing user and domain on every call:

```go
// Create a scoped service bound to bob + domain
bobSvc := svc.SetCurrent(bob, domain)

// Now all calls use bob's identity automatically
myRoles, err := bobSvc.GetCurrentRole()
if err != nil {
    log.Fatal(err)
}
for _, r := range myRoles {
    log.Printf("bob's role: %v", r)
}
```

---

## Step 8: Use the Playground for Testing

The `playground` package creates a fully initialized in-memory environment — perfect for unit
tests and prototyping:

```go
import (
    "testing"
    "github.com/awatercolorpen/caskin/playground"
)

func TestMyFeature(t *testing.T) {
    // Creates a SQLite-backed environment in a temp dir
    stage, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
    if err != nil {
        t.Fatal(err)
    }

    svc    := stage.Service     // fully initialized caskin.IService
    admin  := stage.Admin       // pre-created admin user
    member := stage.Member      // pre-created member user
    domain := stage.Domain      // pre-created domain "school-1"

    // Run your assertions
    roles, err := svc.GetRole(admin, domain)
    if err != nil {
        t.Fatal(err)
    }
    if len(roles) == 0 {
        t.Error("expected at least one role")
    }

    _ = member // use member for permission-denied scenarios
}
```

**Pre-created entities in the playground:**

| Field | Value |
|-------|-------|
| `stage.Superadmin` | `superadmin@qq.com` — has superadmin privileges |
| `stage.Admin` | `teacher@qq.com` — has "admin" role in the domain |
| `stage.Member` | `student@qq.com` — has "member" role in the domain |
| `stage.Domain` | `school-1` — fully bootstrapped domain |

---

## Complete Example

Here is a single self-contained program combining all of the above:

```go
package main

import (
    "log"
    "os"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
)

func main() {
    // 1. Register concrete types
    caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

    // 2. Set up SQLite storage
    dbDir, _ := os.MkdirTemp("", "caskin-demo")
    dbOption := &caskin.DBOption{DSN: dbDir + "/demo.db", Type: "sqlite"}
    db, err := dbOption.NewDB()
    must(err)
    must(db.AutoMigrate(&example.User{}, &example.Role{}, &example.Object{}, &example.Domain{}))

    // 3. Create the service (requires caskin.toml in the working directory)
    svc, err := caskin.New(&caskin.Options{
        Dictionary: &caskin.DictionaryOption{Dsn: "caskin.toml"},
        DB:         dbOption,
    })
    must(err)

    // 4. Bootstrap a domain
    domain := &example.Domain{Name: "demo-org"}
    must(svc.CreateDomain(domain))
    must(svc.ResetDomain(domain))
    must(svc.ResetFeature(domain))

    // 5. Create superadmin
    superadmin := &example.User{Email: "root@example.com"}
    must(svc.CreateUser(superadmin))
    must(svc.AddSuperadmin(superadmin))

    // 6. Create regular users
    alice := &example.User{Email: "alice@example.com"}
    bob   := &example.User{Email: "bob@example.com"}
    must(svc.CreateUser(alice))
    must(svc.CreateUser(bob))

    // 7. Assign roles
    roles, err := svc.GetRole(superadmin, domain)
    must(err)
    must(svc.ModifyUserRolePerUser(superadmin, domain, alice, []*caskin.UserRolePair{{User: alice, Role: roles[0]}}))
    must(svc.ModifyUserRolePerUser(superadmin, domain, bob,   []*caskin.UserRolePair{{User: bob,   Role: roles[1]}}))

    // 8. Check permissions
    pairs, err := svc.GetUserRole(superadmin, domain)
    must(err)
    for _, p := range pairs {
        log.Printf("assignment: user=%v role=%v", p.User, p.Role)
    }

    frontendItems := svc.AuthFrontend(alice, domain)
    log.Printf("alice sees %d frontend items", len(frontendItems))

    log.Println("Done!")
}

func must(err error) {
    if err != nil {
        log.Fatal(err)
    }
}
```

---

## Next Steps

| Topic | Where to look |
|-------|---------------|
| All config options (DB, dictionary, watcher) | [Configuration](./configuration.md) |
| Full service method reference | [API Reference](./api.md) |
| Architecture and design decisions | [Design](./design.md) |
| Redis watcher for multi-instance sync | [Configuration → Watcher](./configuration.md#watcher) |
| Real multi-domain scenario walkthrough | [Examples](../example/) |
