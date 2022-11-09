# Caskin

[![Go](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml/badge.svg?branch=main)](https://github.com/AWaterColorPen/caskin/actions/workflows/go.yml)

Caskin is a multi-domain rbac library for Golang projects. It develops base on [caskin](https://github.com/casbin/casbin) 

## Introduction

### Example

## Documentation

1. [Configuration](./docs/configuration.md) to configure caskin instance and dictionary.
2. [Design](./docs/design.md) for the details of design.
3. [API](./docs/api.md) for the details of caskin service method.

## Getting Started

### Define the dictionary configuration file

Create a new file for example named `caskin.toml` to define
[feature](./docs/configuration.md#feature),
[backend](./docs/configuration.md#backend),
[frontend](./docs/configuration.md#frontend),
[package](./docs/configuration.md#package),
[creator_object](./docs/configuration.md#creator-object),
[creator_role](./docs/configuration.md#creator-role),
[creator_policy](./docs/configuration.md#creator-policy).

```toml
feature = [
    {name = "feature"},
]

backend = [
    {path = "api/feature", method = "GET"},
    {path = "api/feature", method = "POST"},
]

frontend = [
    {name = "feature", type = "menu"},
]

package = [
    {key = "feature", backend = [["api/feature", "GET"], ["api/feature", "POST"]], frontend = [["feature", "menu"]]},
]

creator_object = [
    {name = "role_root", type = "role"},
]

creator_role = [
    {name = "admin"},
    {name = "member"},
]

creator_policy = [
    {role = "admin", object = "role_root", action = ["read", "write", "manage"]},
    {role = "admin", object = "github.com/awatercolorpen/caskin::feature", action = ["read"]},
    {role = "member", object = "role_root", action = ["read"]},
]
```

### To make use of caskin in golang

Register [user-role-object-domain]() instance. 

It should implement the interface of `caskin.User`, `caskin.Role`, `caskin.Object`, `caskin.Domain` generally.
Or use the example implementation in `github.com/awatercolorpen/caskin/example` for the prototype.

```golang
import "github.com/awatercolorpen/caskin"
import "github.com/awatercolorpen/caskin/example"

// register instance type
caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()
```

Create a new [caskin service instance](./docs/configuration.md#service-configuration).

```golang
import "github.com/awatercolorpen/caskin"

// set db option
dbOption := &caskin.DBOption{
	DSN:  "./sqlite.db", 
	Type: "sqlite",
}

// set dictionary option
dictionaryOption := &caskin.DictionaryOption{
	Dsn: "caskin.toml",
}

// build service option
option := &caskin.Options{
	Dictionary: dictionaryOption, 
	DB:         dbOption,
}

// create a new service instance
service, err := caskin.New(option)
```

Initialize first [domain](), and add first superadmin.

```golang
domain := &example.Domain{Name: "school-1"}
superadmin := &example.User{Email: "superadmin@qq.com"}

// create domain
err := service.CreateDomain(domain)

// reset domain by the creator setting from caskin.toml
err := service.ResetDomain(domain)

// reset domain by the feature setting from caskin.toml
err := service.ResetFeature(domain)

// add a user to caskin
err := service.CreateUser(superadmin)

// set a user as superadmin
err := service.AddSuperadmin(p.Superadmin)
```

### To manage the authorization business

Use the `caskin.Service`'s [API](./docs/api.md) to control on authorization management.

```golang
// authorization business: delete one role
err := service.DeleteRole(operatorUser, workingOnDomain, toDeleteRole))
```

Use the `caskin.CurrentService` interface.

```golang
currentService := service.SetCurrent(operatorUser, workingOnDomain)
```

## License

See the [License File](./LICENSE).