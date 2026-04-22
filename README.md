# Caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)
[![Go Reference](https://pkg.go.dev/badge/github.com/awatercolorpen/caskin.svg)](https://pkg.go.dev/github.com/awatercolorpen/caskin)

Caskin is a multi-domain RBAC (Role-Based Access Control) authorization library for Go, built on top of [casbin](https://github.com/casbin/casbin).

It focuses on **managing authorization business** — create and manage users, roles, objects, and domains with fine-grained permission control across multiple tenants.

## Features

- **Multi-domain / multi-tenant** — isolated RBAC per domain
- **Role hierarchy** — role inheritance within a domain
- **Dictionary-driven** — define objects, roles, and policies via TOML config
- **Backend & frontend permissions** — separate API-level and UI-level permission packages
- **Pluggable storage** — supports SQLite, MySQL, PostgreSQL, SQL Server
- **Optional Redis watcher** — sync enforcer across multiple instances

## Quick Start

### 1. Install

```bash
go get github.com/awatercolorpen/caskin
```

### 2. Define a dictionary config

Create `caskin.toml` to declare your features, backend APIs, frontend menus, and initial roles/policies:

```toml
# Features your system exposes
feature = [
    {name = "article"},
]

# Backend API endpoints
backend = [
    {path = "api/article", method = "GET"},
    {path = "api/article", method = "POST"},
]

# Frontend menu/UI items
frontend = [
    {name = "article", type = "menu"},
]

# Packages bundle backend + frontend into a logical permission unit
package = [
    {key = "article", backend = [["api/article", "GET"], ["api/article", "POST"]], frontend = [["article", "menu"]]},
]

# Initial objects created when a domain is reset
creator_object = [
    {name = "role_root", type = "role"},
]

# Initial roles created when a domain is reset
creator_role = [
    {name = "admin"},
    {name = "member"},
]

# Initial policies assigned to those roles
creator_policy = [
    {role = "admin", object = "role_root", action = ["read", "write", "manage"]},
    {role = "admin", object = "github.com/awatercolorpen/caskin::article", action = ["read"]},
    {role = "member", object = "role_root", action = ["read"]},
]
```

### 3. Implement the four core interfaces

Caskin requires four types: **User**, **Role**, **Object**, and **Domain**. You can embed them in your own models, or use the ready-made implementations in the `example` package as a starting point.

```go
import "github.com/awatercolorpen/caskin/example"

// example.User  implements caskin.User
// example.Role  implements caskin.Role
// example.Object implements caskin.Object
// example.Domain implements caskin.Domain
```

Each type must satisfy its interface. For example, `caskin.User` requires `GetID()`, `SetID()`, `Encode()`, and `Decode()`.

### 4. Create a service instance

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
    // 1. Register your concrete types (generics-based, called once at startup)
    caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

    // 2. Configure storage (SQLite for local dev; swap for MySQL/Postgres in production)
    dbOption := &caskin.DBOption{
        DSN:  "./caskin.db",
        Type: "sqlite",
    }

    // 3. Auto-migrate your tables
    db, err := dbOption.NewDB()
    if err != nil {
        log.Fatal(err)
    }
    db.AutoMigrate(&example.User{}, &example.Role{}, &example.Object{}, &example.Domain{})

    // 4. Build the service
    service, err := caskin.New(&caskin.Options{
        Dictionary: &caskin.DictionaryOption{Dsn: "caskin.toml"},
        DB:         dbOption,
    })
    if err != nil {
        log.Fatal(err)
    }

    // 5. Bootstrap: create a domain and a superadmin
    domain := &example.Domain{Name: "my-org"}
    superadmin := &example.User{Email: "admin@example.com"}

    _ = service.CreateDomain(domain)
    _ = service.ResetDomain(domain)   // creates initial objects + roles from caskin.toml
    _ = service.ResetFeature(domain)  // registers features + backend/frontend definitions

    _ = service.CreateUser(superadmin)
    _ = service.AddSuperadmin(superadmin)

    // 6. Permission check
    roles, _ := service.GetRole(superadmin, domain)
    log.Printf("domain %q has %d roles", domain.Name, len(roles))
}
```

### 5. Manage authorization

Use the operator (`superadmin`) and the target domain to call any service method:

```go
// Create a regular user
user := &example.User{Email: "alice@example.com"}
service.CreateUser(user)

// Assign a role to the user
roles, _ := service.GetRole(superadmin, domain)
adminRole := roles[0] // first role is "admin" by default
service.ModifyUserRolePerRole(superadmin, domain, adminRole, []*caskin.UserRolePair{
    {User: user, Role: adminRole},
})

// Check what roles the user has
pairs, _ := service.GetUserRole(superadmin, domain)
for _, p := range pairs {
    log.Printf("user %v → role %v", p.User, p.Role)
}

// Use a scoped CurrentService for a specific operator+domain context
current := service.SetCurrent(user, domain)
myRoles, _ := current.GetCurrentRole()
myObjects, _ := current.GetCurrentObject()
```

### 6. Use the example package directly (for prototyping / tests)

The `playground` package sets up a fully initialized in-memory environment for quick testing:

```go
import "github.com/awatercolorpen/caskin/playground"

stage, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
// stage.Service  — the caskin service
// stage.Superadmin, stage.Admin, stage.Member — pre-created users
// stage.Domain   — pre-created domain with roles and policies
```

## Documentation

| Doc | Description |
|-----|-------------|
| [Getting Started](./docs/getting-started.md) | Step-by-step guide with complete runnable examples |
| [Configuration](./docs/configuration.md) | All config options: DB, dictionary, watcher |
| [Design](./docs/design.md) | Architecture and design decisions |
| [API Reference](./docs/api-reference.md) | Full type and method reference with parameters, return values, and usage notes |

## Storage backends

| Driver | DBOption.Type | DSN example |
|--------|---------------|-------------|
| SQLite | `"sqlite"` | `./caskin.db` |
| MySQL | `"mysql"` | `user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True` |
| PostgreSQL | `"postgres"` | `host=localhost user=postgres password=pass dbname=caskin port=5432` |

## Optional: Redis watcher (multi-instance sync)

```go
service, err := caskin.New(&caskin.Options{
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

## License

[MIT](./LICENSE)
