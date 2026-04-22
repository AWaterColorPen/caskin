# API Reference

This document describes every public type, interface, and method in caskin. For a hands-on introduction, see [Getting Started](./getting-started.md).

---

## Table of Contents

- [Core Types](#core-types)
  - [User](#user)
  - [Domain](#domain)
  - [Role](#role)
  - [Object](#object)
  - [ObjectData](#objectdata)
  - [Policy](#policy)
  - [UserRolePair](#userrolepair)
  - [Directory](#directory)
  - [DirectoryRequest](#directoryrequest)
  - [DirectoryResponse](#directoryresponse)
- [Constants and Actions](#constants-and-actions)
  - [Actions](#actions)
  - [Object Types](#object-types)
  - [Directory Search Types](#directory-search-types)
- [Service Interfaces](#service-interfaces)
  - [IService](#iservice)
  - [IBaseService](#ibaseservice)
  - [IFeatureService](#ifeatureservice)
  - [IDirectoryService](#idirectoryservice)
  - [ICurrentService](#icurrentservice)
- [Constructor](#constructor)
  - [New](#new)
  - [Options](#options)
- [Utility Functions](#utility-functions)
  - [Check](#check)
  - [Filter](#filter)
  - [ID and IDMap](#id-and-idmap)

---

## Core Types

### User

```go
type User interface {
    idInterface    // GetID() uint64; SetID(uint64)
    codeInterface  // Encode() string
}
```

Represents an actor in the system (a real person, a service account, etc.).

**Required methods:**

| Method | Description |
|--------|-------------|
| `GetID() uint64` | Returns the unique integer ID |
| `SetID(uint64)` | Sets the integer ID |
| `Encode() string` | Returns a stable string used as casbin subject (e.g. `"user:42"`) |

**Usage:** Users are global — they are not tied to any domain. A user participates in a domain by having at least one role assigned in that domain.

---

### Domain

```go
type Domain interface {
    idInterface    // GetID() uint64; SetID(uint64)
    codeInterface  // Encode() string
}
```

Represents an isolated permission scope such as a tenant or organisation. All role and policy operations are scoped to a domain.

**Required methods:**

| Method | Description |
|--------|-------------|
| `GetID() uint64` | Returns the unique integer ID |
| `SetID(uint64)` | Sets the integer ID |
| `Encode() string` | Returns a stable string used as casbin domain token |

---

### Role

```go
type Role interface {
    ObjectData
    codeInterface  // Encode() string
}
```

A named permission bundle inside a domain. Roles are themselves protected resources (`ObjectData`) stored under the built-in `ObjectTypeRole` object tree. This means:

- A user needs `Read` permission on the role object to list or assign it.
- A user needs `Write` permission to create/update/delete it.
- A user needs `Manage` permission to modify its policies.

Roles support hierarchical inheritance via [`AddRoleG`](#addroleg--removeroleg).

---

### Object

```go
type Object interface {
    idInterface       // GetID() uint64; SetID(uint64)
    codeInterface     // Encode() string
    parentInterface   // GetParentID() uint64; SetParentID(uint64)
    domainInterface   // GetDomainID() uint64; SetDomainID(uint64)
    GetObjectType() string
}
```

Represents a resource or permission node in the hierarchy. Objects are arranged in a tree — a policy granted on a parent is automatically inherited by all child objects.

**Required methods:**

| Method | Description |
|--------|-------------|
| `GetID() / SetID()` | Integer ID |
| `GetParentID() / SetParentID()` | Parent node in the object tree (`0` = root) |
| `GetDomainID() / SetDomainID()` | Owning domain |
| `GetObjectType() string` | Logical type tag (e.g. `"menu"`, `"api"`, `"role"`) |
| `Encode() string` | Stable casbin resource string |

**Constraint:** Object hierarchy is limited to a depth of 10 levels. Attempting to exceed this returns an error.

---

### ObjectData

```go
type ObjectData interface {
    idInterface      // GetID() uint64; SetID(uint64)
    domainInterface  // GetDomainID() uint64; SetDomainID(uint64)
    GetObjectID() uint64
    SetObjectID(uint64)
}
```

A domain-specific data record protected by an `Object`. For example, a document entity might be linked to a "documents" object node. All permission checks (`Read`, `Write`, `Manage`) are evaluated against the linked object.

**Required methods:**

| Method | Description |
|--------|-------------|
| `GetObjectID() uint64` | ID of the protecting `Object` |
| `SetObjectID(uint64)` | Sets the protecting `Object` ID |

---

### Policy

```go
type Policy struct {
    Role   Role   `json:"role"`
    Object Object `json:"object"`
    Domain Domain `json:"domain"`
    Action Action `json:"action"`
}
```

A 4-tuple that grants a role permission to perform an action on an object within a domain. This is the fundamental unit of access control in caskin.

**Method:**

| Method | Description |
|--------|-------------|
| `Key() string` | Returns a stable string identifier for set-difference operations |

**Example:** _"In domain `school-1`, role `teacher` can `read` object `course-management`."_

---

### UserRolePair

```go
type UserRolePair struct {
    User User `json:"user"`
    Role Role `json:"role"`
}
```

Binds a user to a role within a domain. Used in [`AddUserRole`](#addusurrole--removeusurrole) and related methods.

---

### Directory

```go
type Directory struct {
    Object
    AllDirectoryCount  uint64 `json:"all_directory_count"`
    AllItemCount       uint64 `json:"all_item_count"`
    TopDirectoryCount  uint64 `json:"top_directory_count"`
    TopItemCount       uint64 `json:"top_item_count"`
}
```

Decorates an `Object` with aggregate counts for the subtree. Returned by directory listing operations.

| Field | Description |
|-------|-------------|
| `AllDirectoryCount` | Total directory nodes in the entire subtree |
| `AllItemCount` | Total leaf items (`ObjectData`) in the entire subtree |
| `TopDirectoryCount` | Direct child directories only |
| `TopItemCount` | Direct child items only |

---

### DirectoryRequest

```go
type DirectoryRequest struct {
    To              uint64
    ID              []uint64
    Type            string
    Policy          string
    SearchType      string
    CountDirectory  func([]uint64) (map[uint64]uint64, error)
    ActionDirectory func([]uint64) error
}
```

Parameter bag for all directory operations. Fields are selectively used depending on which method is called.

| Field | Type | Description |
|-------|------|-------------|
| `To` | `uint64` | Target parent object ID for move operations |
| `ID` | `[]uint64` | Object IDs to operate on |
| `Type` | `string` | `ObjectType` filter |
| `SearchType` | `string` | `DirectorySearchAll` or `DirectorySearchTop` |
| `CountDirectory` | `func` | Optional callback to count items per directory |
| `ActionDirectory` | `func` | Optional callback for side effects on a directory set |

---

### DirectoryResponse

```go
type DirectoryResponse struct {
    DoneDirectoryCount uint64 `json:"done_directory_count,omitempty"`
    DoneItemCount      uint64 `json:"done_item_count,omitempty"`
    ToDoDirectoryCount uint64 `json:"to_do_directory_count,omitempty"`
    ToDoItemCount      uint64 `json:"to_do_item_count,omitempty"`
}
```

Summarises the outcome of a directory move or copy.

| Field | Description |
|-------|-------------|
| `DoneDirectoryCount` | Directories already at the destination (skipped) |
| `DoneItemCount` | Items already at the destination (skipped) |
| `ToDoDirectoryCount` | Directories successfully moved/copied |
| `ToDoItemCount` | Items successfully moved/copied |

---

## Constants and Actions

### Actions

| Constant | Value | Description |
|----------|-------|-------------|
| `Read` | `"read"` | Permission to view/query a resource |
| `Write` | `"write"` | Permission to create or modify a resource |
| `Manage` | `"manage"` | Administrative permission; typically implies Read and Write |

Custom actions are supported — any string value can be used as an action in a `Policy`.

---

### Object Types

| Constant | Value | Description |
|----------|-------|-------------|
| `ObjectTypeRole` | `"role"` | Built-in type for the `Role` object tree |

Your application defines its own object types (e.g. `"menu"`, `"api"`, `"document"`) by registering them with the dictionary. See [Configuration](./configuration.md).

---

### Directory Search Types

| Constant | Value | Description |
|----------|-------|-------------|
| `DirectorySearchAll` | `"all"` | Traverse the full subtree (all descendants) |
| `DirectorySearchTop` | `"top"` | Return only direct children |

---

## Service Interfaces

### IService

```go
type IService interface {
    IBaseService
    IDirectoryService
    IFeatureService
    ICurrentService
}
```

The top-level service interface returned by [`New`](#new). Composes all four sub-interfaces.

---

### IBaseService

Core CRUD operations for users, domains, objects, roles, and their relationships.

#### Superadmin

```go
AddSuperadmin(User) error
DeleteSuperadmin(User) error
GetSuperadmin() ([]User, error)
```

Superadmins bypass all permission checks. They are stored in a reserved casbin domain (`superdomain`) and are independent of any regular domain.

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `AddSuperadmin` | `user User` | `error` | Grants superadmin status to a user |
| `DeleteSuperadmin` | `user User` | `error` | Revokes superadmin status |
| `GetSuperadmin` | — | `[]User, error` | Lists all superadmin users |

---

#### User Management

```go
CreateUser(User) error
RecoverUser(User) error
DeleteUser(User) error
UpdateUser(User) error
```

No permission checks are applied — these operations are administrative. Users are global (not domain-scoped).

| Method | Parameters | Description |
|--------|-----------|-------------|
| `CreateUser` | `user User` | Creates a user. Returns `ErrAlreadyExists` if ID is taken |
| `RecoverUser` | `user User` | Restores a soft-deleted user. Returns `ErrAlreadyExists` if not deleted, `ErrNotExists` if never created |
| `DeleteUser` | `user User` | Soft-deletes a user, removing all their role assignments across every domain |
| `UpdateUser` | `user User` | Updates a user's properties |

---

#### Domain Management

```go
CreateDomain(Domain) error
RecoverDomain(Domain) error
DeleteDomain(Domain) error
UpdateDomain(Domain) error
GetDomain() ([]Domain, error)
ResetDomain(Domain) error
```

No permission checks — domain management is administrative.

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `CreateDomain` | `domain Domain` | `error` | Creates a domain |
| `RecoverDomain` | `domain Domain` | `error` | Restores a soft-deleted domain |
| `DeleteDomain` | `domain Domain` | `error` | Soft-deletes a domain and removes all user-role assignments within it. Does **not** remove role inheritance (g) or object policies (g2) |
| `UpdateDomain` | `domain Domain` | `error` | Updates domain properties |
| `GetDomain` | — | `[]Domain, error` | Lists all domains |
| `ResetDomain` | `domain Domain` | `error` | Re-initialises a domain (e.g. re-creates built-in role/object trees) |

---

#### Object Management

```go
CreateObject(user User, domain Domain, object Object) error
RecoverObject(user User, domain Domain, object Object) error
DeleteObject(user User, domain Domain, object Object) error
UpdateObject(user User, domain Domain, object Object) error
GetObject(user User, domain Domain, action Action, types ...ObjectType) ([]Object, error)
GetObjectHierarchyLevel(user User, domain Domain, object Object) (int, error)
```

Object operations are **permission-checked**: the calling user must have `Manage` permission on the parent to create or modify objects.

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `CreateObject` | `user, domain, object` | `error` | Creates an object under its parent. User must have `Manage` on the parent |
| `RecoverObject` | `user, domain, object` | `error` | Restores a soft-deleted object |
| `DeleteObject` | `user, domain, object` | `error` | Soft-deletes an object and all its policies |
| `UpdateObject` | `user, domain, object` | `error` | Updates object properties. Cannot change `ObjectType`. Cannot reparent to a descendant |
| `GetObject` | `user, domain, action, types...` | `[]Object, error` | Lists objects the user has the given `action` on, filtered by optional types |
| `GetObjectHierarchyLevel` | `user, domain, object` | `int, error` | Returns how many levels deep the object is in its tree |

**Update constraints:**
- `ObjectType` is immutable — changing it returns `ErrCantChangeObjectType`
- Cannot set a descendant as the new parent — returns `ErrParentToDescendant`
- Maximum hierarchy depth is 10 levels

---

#### Role Management

```go
CreateRole(user User, domain Domain, role Role) error
RecoverRole(user User, domain Domain, role Role) error
DeleteRole(user User, domain Domain, role Role) error
UpdateRole(user User, domain Domain, role Role) error
GetRole(user User, domain Domain) ([]Role, error)
```

Roles are stored as `ObjectData` under the built-in `ObjectTypeRole` tree. Permission is checked against the role's linked object.

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `CreateRole` | `user, domain, role` | `error` | Creates a role. User must have `Write` on the role object |
| `RecoverRole` | `user, domain, role` | `error` | Restores a soft-deleted role |
| `DeleteRole` | `user, domain, role` | `error` | Removes the role, all its `g` assignments, and its `p` policies |
| `UpdateRole` | `user, domain, role` | `error` | Updates role properties |
| `GetRole` | `user, domain` | `[]Role, error` | Lists roles the user can `Read` in the domain |

---

#### User–Role Assignments

```go
AddUserRole(user User, domain Domain, pairs []*UserRolePair) error
RemoveUserRole(user User, domain Domain, pairs []*UserRolePair) error
AddRoleG(user User, domain Domain, child Role, parent Role) error
RemoveRoleG(user User, domain Domain, child Role, parent Role) error
```

| Method | Parameters | Description |
|--------|-----------|-------------|
| `AddUserRole` | `user, domain, pairs` | Assigns roles to users. Caller must have `Write` on each role |
| `RemoveUserRole` | `user, domain, pairs` | Removes role assignments |
| `AddRoleG` | `user, domain, child, parent` | Makes `child` inherit all permissions of `parent` |
| `RemoveRoleG` | `user, domain, child, parent` | Removes the inheritance edge |

---

#### User–Role Query Methods

```go
GetUserByDomain(domain Domain) ([]User, error)
GetDomainByUser(user User) ([]Domain, error)
GetUserRole(user User, domain Domain) ([]*UserRolePair, error)
GetUserRoleByUser(user User, domain Domain, target User) ([]*UserRolePair, error)
GetUserRoleByRole(user User, domain Domain, target Role) ([]*UserRolePair, error)
ModifyUserRolePerUser(user User, domain Domain, target User, pairs []*UserRolePair) error
ModifyUserRolePerRole(user User, domain Domain, target Role, pairs []*UserRolePair) error
```

| Method | Returns | Description |
|--------|---------|-------------|
| `GetUserByDomain` | `[]User` | All users that have at least one role in the domain |
| `GetDomainByUser` | `[]Domain` | All domains the user belongs to |
| `GetUserRole` | `[]*UserRolePair` | All user–role assignments visible to the calling user in the domain |
| `GetUserRoleByUser` | `[]*UserRolePair` | Assignments for a specific user |
| `GetUserRoleByRole` | `[]*UserRolePair` | Assignments for a specific role |
| `ModifyUserRolePerUser` | `error` | Replaces all role assignments for a user (diff-based, minimal casbin writes) |
| `ModifyUserRolePerRole` | `error` | Replaces all user assignments for a role |

---

#### Policy Management

```go
GetPolicy(user User, domain Domain) ([]*Policy, error)
GetPolicyByRole(user User, domain Domain, role Role) ([]*Policy, error)
ModifyPolicyPerRole(user User, domain Domain, role Role, policies []*Policy) error
```

Policies bind roles to objects with an action. Only policies for objects the calling user has `Manage` permission on are visible/modifiable.

| Method | Description |
|--------|-------------|
| `GetPolicy` | Lists all (role, object, action) triples visible to the calling user |
| `GetPolicyByRole` | Lists policies for a specific role |
| `ModifyPolicyPerRole` | Replaces policies for a role (diff-based: adds new, removes stale) |

---

#### ObjectData Management

```go
CreateObjectData(user User, domain Domain, data ObjectData, objType ObjectType) error
RecoverObjectData(user User, domain Domain, data ObjectData) error
DeleteObjectData(user User, domain Domain, data ObjectData) error
UpdateObjectData(user User, domain Domain, data ObjectData, objType ObjectType) error
```

Generic object-data CRUD. The concrete GORM model is resolved via `objType`.

| Method | Permission Required | Description |
|--------|-------------------|-------------|
| `CreateObjectData` | `Write` on linked object | Creates a new record |
| `RecoverObjectData` | `Write` | Restores a soft-deleted record |
| `DeleteObjectData` | `Write` | Soft-deletes the record |
| `UpdateObjectData` | `Write` | Updates the record |

**Check methods** (same signatures but return `error` without mutating state):

```go
CheckCreateObjectData(user, domain, data, objType) error
CheckRecoverObjectData(user, domain, data) error
CheckDeleteObjectData(user, domain, data) error
CheckWriteObjectData(user, domain, data, objType) error
CheckUpdateObjectData(user, domain, data, objType) error
CheckModifyObjectData(user, domain, data) error
CheckGetObjectData(user, domain, data) error
```

Use these to validate permissions before performing an action, or to surface informative errors in your API layer.

---

### IFeatureService

Feature operations are for managing frontend/backend access permissions defined in the dictionary configuration.

```go
AuthBackend(user User, domain Domain, backend *Backend) error
AuthFrontend(user User, domain Domain) []*Frontend
GetFeature(user User, domain Domain) ([]*Feature, error)
GetFeaturePolicy(user User, domain Domain) ([]*Policy, error)
GetFeaturePolicyByRole(user User, domain Domain, role Role) ([]*Policy, error)
ModifyFeaturePolicyPerRole(user User, domain Domain, role Role, policies []*Policy) error
ResetFeature(domain Domain) error
```

| Method | Parameters | Returns | Description |
|--------|-----------|---------|-------------|
| `AuthBackend` | `user, domain, backend` | `error` | Returns `nil` if the user has `Read` on the backend's object. Returns `ErrNoBackendPermission` otherwise. Use this as a middleware guard |
| `AuthFrontend` | `user, domain` | `[]*Frontend` | Returns the list of frontend items the user can access. Use this to build permission-filtered navigation menus |
| `GetFeature` | `user, domain` | `[]*Feature, error` | Returns all features (both backend and frontend) the user can `Read` |
| `GetFeaturePolicy` | `user, domain` | `[]*Policy, error` | Returns feature policies for all roles the user can see |
| `GetFeaturePolicyByRole` | `user, domain, role` | `[]*Policy, error` | Returns feature policies for a specific role |
| `ModifyFeaturePolicyPerRole` | `user, domain, role, policies` | `error` | Replaces feature policies for a role. Only `Read` action is allowed on features; the action field in input is ignored |
| `ResetFeature` | `domain` | `error` | Syncs feature objects from the dictionary into the domain's object tree |

**Feature vs Policy distinction:**

- `GetPolicy` / `ModifyPolicyPerRole` — manage permissions on **data objects** (things with `ObjectData` records in your database)
- `GetFeaturePolicy` / `ModifyFeaturePolicyPerRole` — manage permissions on **feature objects** (virtual objects defined in the dictionary TOML, with no database rows)

---

### IDirectoryService

Directory operations treat the object tree as a file-system-like structure with folders and items.

```go
CreateDirectory(user User, domain Domain, object Object) error
UpdateDirectory(user User, domain Domain, object Object) error
DeleteDirectory(user User, domain Domain, req *DirectoryRequest) error
GetDirectory(user User, domain Domain, req *DirectoryRequest) ([]*Directory, error)
MoveDirectory(user User, domain Domain, req *DirectoryRequest) (*DirectoryResponse, error)
MoveItem(user User, domain Domain, data ObjectData, req *DirectoryRequest) (*DirectoryResponse, error)
CopyItem(user User, domain Domain, data ObjectData, req *DirectoryRequest) (*DirectoryResponse, error)
```

| Method | Description |
|--------|-------------|
| `CreateDirectory` | Creates a new directory node (an `Object` used as a folder) |
| `UpdateDirectory` | Renames or re-parents a directory |
| `DeleteDirectory` | Removes a directory and its full subtree based on `req.ID` and `req.Type` |
| `GetDirectory` | Lists directories matching `req`, enriched with item counts if `req.CountDirectory` is provided |
| `MoveDirectory` | Moves directories in `req.ID` to the target `req.To` |
| `MoveItem` | Moves an `ObjectData` item to the target directory in `req.To` |
| `CopyItem` | Copies an `ObjectData` item to the target directory in `req.To` |

---

### ICurrentService

Binds a user and domain to a service instance so individual method calls don't need to repeat them.

```go
SetCurrent(user User, domain Domain) IService
```

Returns a new `IService` instance with `user` and `domain` pre-filled. All `WithCurrent` methods on the returned instance use these bound values.

**ObjectData methods with current context:**

```go
CreateObjectDataWithCurrent(data ObjectData, objType ObjectType) error
RecoverObjectDataWithCurrent(data ObjectData) error
DeleteObjectDataWithCurrent(data ObjectData) error
UpdateObjectDataWithCurrent(data ObjectData, objType ObjectType) error

CheckCreateObjectDataWithCurrent(data ObjectData, objType ObjectType) error
CheckRecoverObjectDataWithCurrent(data ObjectData) error
CheckDeleteObjectDataWithCurrent(data ObjectData) error
CheckWriteObjectDataWithCurrent(data ObjectData, objType ObjectType) error
CheckUpdateObjectDataWithCurrent(data ObjectData, objType ObjectType) error
CheckModifyObjectDataWithCurrent(data ObjectData) error
CheckGetObjectDataWithCurrent(data ObjectData) error
```

**Example:**

```go
// Without SetCurrent — verbose
svc.CreateObjectData(currentUser, currentDomain, doc, "document")

// With SetCurrent — cleaner for handler code
cs := svc.SetCurrent(currentUser, currentDomain)
cs.CreateObjectDataWithCurrent(doc, "document")
```

---

## Constructor

### New

```go
func New(options *Options, opts ...Option) (IService, error)
```

Creates and returns a fully initialised `IService`.

**Parameters:**

| Parameter | Type | Description |
|-----------|------|-------------|
| `options` | `*Options` | Required configuration (database, dictionary, etc.) |
| `opts` | `...Option` | Optional functional option overrides |

**Returns:** `(IService, error)` — error if database connection, casbin model, or dictionary initialisation fails.

**Example:**

```go
svc, err := caskin.New(&caskin.Options{
    DB: &caskin.DBOption{
        DSN:    "user:pass@tcp(localhost:3306)/mydb?charset=utf8mb4&parseTime=True",
        Driver: "mysql",
    },
    Dictionary: &caskin.DictionaryOption{
        File: "configs/dictionary.toml",
    },
})
if err != nil {
    log.Fatal(err)
}
```

---

### Options

```go
type Options struct {
    DefaultSuperadminDomainName string          `json:"default_superadmin_domain_name"`
    DefaultSuperadminRoleName   string          `json:"default_superadmin_role_name"`
    Dictionary                  *DictionaryOption
    DB                          *DBOption
    Watcher                     *WatcherOption
}
```

| Field | Default | Description |
|-------|---------|-------------|
| `DefaultSuperadminDomainName` | `"superadmin_domain"` | Override the reserved superadmin domain name |
| `DefaultSuperadminRoleName` | `"superadmin_role"` | Override the reserved superadmin role name |
| `Dictionary` | nil (empty in-memory) | Dictionary configuration for features and creators |
| `DB` | — | **Required.** Database connection settings |
| `Watcher` | nil (disabled) | Optional Redis watcher for distributed policy synchronisation |

---

## Utility Functions

### Check

```go
func Check[T any](e IEnforcer, u User, d Domain, one T, action Action) bool
```

Low-level permission check. Accepts either an `Object` or an `ObjectData`. If `T` is `ObjectData`, the check is forwarded to its linked `Object`.

**Example:**

```go
if caskin.Check(enforcer, user, domain, myDocument, caskin.Write) {
    // user can write to myDocument
}
```

---

### Filter

```go
func Filter[T any](e IEnforcer, u User, d Domain, action Action, source []T) []T
```

Filters a slice, returning only items the user has the given `action` on.

**Example:**

```go
// Keep only objects the user can manage
manageable := caskin.Filter(enforcer, user, domain, caskin.Manage, allObjects)
```

---

### ID and IDMap

```go
func ID[E idInterface](in []E) []uint64
func IDMap[E idInterface](in []E) map[uint64]E
```

Convenience helpers for working with slices of entities.

| Function | Description |
|----------|-------------|
| `ID(slice)` | Extracts `[]uint64` IDs from a slice |
| `IDMap(slice)` | Converts a slice to a `map[uint64]E` keyed by ID |

**Example:**

```go
ids := caskin.ID(roles)                // []uint64{1, 2, 3}
roleMap := caskin.IDMap(roles)         // map[uint64]Role{1: role1, ...}
role, ok := roleMap[targetID]
```

---

## Error Reference

| Error | Meaning |
|-------|---------|
| `ErrAlreadyExists` | Attempted to create a record that already exists |
| `ErrNotExists` | Attempted to operate on a record that doesn't exist |
| `ErrNoBackendPermission` | User lacks access to the requested backend endpoint |
| `ErrCantOperateRootObject` | Attempted to modify a root object (parentID = 0) |
| `ErrInValidObjectType` | Parent and child objects have mismatched types |
| `ErrParentToDescendant` | Attempted to set a descendant as the new parent |
| `ErrParentCanNotBeItself` | Attempted to set an object as its own parent |
| `ErrCantChangeObjectType` | Attempted to change an object's type after creation |

---

*For configuration details (DB, Dictionary, Watcher options), see [Configuration](./configuration.md).*  
*For a step-by-step walkthrough, see [Getting Started](./getting-started.md).*
