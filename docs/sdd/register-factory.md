# System Design Document: Register Factory

**Author**: slyao (via long-haul agent)  
**Created**: 2026-05-28  
**Status**: Draft  
**Cross-references**: [Object Hierarchy SDD](object-hierarchy.md), [Rules Engine SDD](rules-engine.md)

---

## 1. Overview

The Register Factory is caskin's **type registration and decoding subsystem**. It bridges the gap between casbin's string-based entity model (subjects, objects, domains are all strings) and caskin's strongly-typed Go structs.

At application startup, the caller registers concrete types via the generic `Register[U, R, O, D]()` function. At runtime, the factory decodes casbin strings back into typed structs and constructs zero-value instances for database operations.

---

## 2. Goals & Non-Goals

### Goals

- Provide a compile-time-safe API for registering application-specific entity types
- Decode casbin string tokens into concrete Go structs at runtime
- Supply zero-value constructors for use by MetaDB (GORM layer)
- Keep the public API surface minimal (one generic call at startup)

### Non-Goals

- Supporting multiple concurrent type registries (multi-tenant in-process)
- Runtime hot-swapping of registered types
- Validating business logic of decoded entities (that's the rules engine's job)

---

## 3. Architecture Context

```
┌──────────────────────────────────────────────────────────────┐
│                     Application Startup                        │
│   caskin.Register[*MyUser, *MyRole, *MyObject, *MyDomain]()   │
└──────────────────────────────┬───────────────────────────────┘
                               │ writes singleton
                               ▼
┌──────────────────────────────────────────────────────────────┐
│                   defaultFactory (global)                      │
│         builtinRegister[U, R, O, D] implements Factory         │
└──────┬──────────────┬──────────────────┬─────────────────────┘
       │              │                  │
       ▼              ▼                  ▼
  Enforcer         MetaDB           Server layer
  (Decode)      (NewUser/etc)     (DefaultFactory())
```

The factory sits at the center of caskin's type system, consumed by:
- **Enforcer** — decodes casbin policy strings into typed structs during permission checks
- **MetaDB** — uses `NewUser()`/`NewRole()`/etc. to construct GORM model instances
- **Server layer** — accesses the factory via `DefaultFactory()` for ad-hoc operations

---

## 4. Interface Design

### 4.1 Factory Interface

```go
type Factory interface {
    User(string) (User, error)
    Role(string) (Role, error)
    Object(string) (Object, error)
    Domain(string) (Domain, error)
    NewUser() User
    NewRole() Role
    NewObject() Object
    NewDomain() Domain
    MetadataDB(db *gorm.DB) MetaDB
}
```

The interface is intentionally symmetric: four `Decode` methods + four `New` methods + one `MetadataDB` constructor. Each entity category (User, Role, Object, Domain) maps 1:1 to a casbin model position.

### 4.2 codeInterface (Encoding Contract)

```go
type codeInterface interface {
    Encode() string
    Decode(string) error
}
```

Every registered type must implement `codeInterface`. This is the fundamental contract that makes string↔struct conversion possible. The `Encode()` output is stored in casbin; `Decode()` reverses it.

---

## 5. Registration Protocol

### 5.1 Compile-Time Safety

```go
func Register[U User, R Role, O Object, D Domain]()
```

The generic constraints ensure that:
- `U` satisfies `User` (= `idInterface` + `codeInterface`)
- `R` satisfies `Role` (= `ObjectData` + `codeInterface`)
- `O` satisfies `Object` (= `idInterface` + `codeInterface` + `parentInterface` + `domainInterface` + `GetObjectType()`)
- `D` satisfies `Domain` (= `idInterface` + `codeInterface`)

A mismatched type produces a compile error, not a runtime panic.

### 5.2 Runtime Initialization

When `Register[U, R, O, D]()` executes:

1. Constructs `builtinRegister[U, R, O, D]{}` 
2. Populates each candidate list with one zero-value instance of the registered type
3. **Appends `&NamedObject{}` to the object candidates** — this is a built-in fallback for simple named objects
4. Assigns the constructed register to the package-level `defaultFactory`

### 5.3 Single-Call Invariant

`Register` writes directly to `defaultFactory` without protection. It is designed to be called **exactly once** at program startup, before any concurrent access. Calling it more than once silently replaces the previous factory.

---

## 6. Decode Strategy

### 6.1 Candidate Iteration

```go
func decode[T codeInterface](code string, candidate []T) (T, error) {
    for _, v := range candidate {
        e := newByE(v)
        if err := e.Decode(code); err == nil {
            return e, nil
        }
    }
    var zero T
    return zero, fmt.Errorf("no register factory for %v", code)
}
```

The decode function:
1. Iterates over all registered candidates for a given entity category
2. For each candidate, creates a fresh zero-value copy via `newByE`
3. Attempts `Decode(code)` on the fresh instance
4. Returns the first successful decode; if all fail, returns a generic error

### 6.2 Candidate Counts

| Category | Typical Candidates | Reason |
|---|---|---|
| User | 1 | Single user type per application |
| Role | 1 | Single role type per application |
| Object | 2 | Application type + built-in `NamedObject` |
| Domain | 1 | Single domain type per application |

The O(n) scan is negligible because n ∈ {1, 2} in practice.

---

## 7. Instance Construction (Reflect Layer)

### 7.1 newByT — Type-Parameter Constructor

```go
func newByT[T any]() T
```

Creates a zero-value instance of `T`. If `T` is a pointer type (the normal case for interface-satisfying types), it allocates via `reflect.New(v.Type().Elem())` and returns the interface cast.

**Usage**: Called by `NewUser()`, `NewRole()`, etc. — the "factory method" side of the Factory.

### 7.2 newByE — Exemplar-Based Constructor

```go
func newByE[E any](e E) E
```

Creates a new zero-value instance of the **same concrete type** as the given exemplar `e`. Uses `reflect.Indirect` + `reflect.New` to handle pointer types.

**Usage**: Called by `decode()` to create a fresh target for each decode attempt, ensuring candidates are never mutated.

### 7.3 Why Reflect?

Go generics (1.18+) cannot express "give me a new instance of the concrete type behind this interface" without reflection. The type parameter gives compile-time safety; reflect gives runtime instantiation. This is a well-known Go pattern for generic factories.

---

## 8. MetadataDB Integration

```go
func (b *builtinRegister[U, R, O, D]) MetadataDB(db *gorm.DB) MetaDB {
    return &builtinMetadataDB[U, R, O, D]{DB: db}
}
```

The factory is the only component that knows all four concrete types at once. This makes it the natural place to construct `MetaDB`, which needs type-parameterized GORM queries (e.g., `db.Model(new(U)).Where(...)`).

The `MetadataDB` method is called once during `caskin.New()` initialization and the resulting `MetaDB` is stored for the lifetime of the service.

---

## 9. Type System Mapping

```
        idInterface                    codeInterface
        ├── GetID() uint64             ├── Encode() string
        └── SetID(uint64)              └── Decode(string) error

User  = idInterface + codeInterface
Domain = idInterface + codeInterface
Role  = ObjectData + codeInterface
           └── ObjectData = idInterface + domainInterface + GetObjectID/SetObjectID
Object = idInterface + codeInterface + parentInterface + domainInterface + GetObjectType()
```

This hierarchy ensures that every entity can:
- Be identified by a numeric ID (for GORM/DB operations)
- Be encoded/decoded to strings (for casbin operations)
- Carry domain-specific metadata appropriate to its role in the permission model

---

## 10. Performance Characteristics

| Operation | Cost | Hot Path? |
|---|---|---|
| `Register[...]()` | O(1) + 4 reflect.New | No (startup only) |
| `factory.User(s)` | 1 reflect.New + 1 Decode | Yes (every permission check) |
| `factory.Object(s)` | 1-2 reflect.New + 1-2 Decode | Yes |
| `factory.NewUser()` | 1 reflect.New | Moderate (DB operations) |

**Observations**:
- `reflect.New` allocates on the heap; in a high-throughput permission-check loop, this creates GC pressure
- The decode path (especially for Objects with 2 candidates) does speculative allocation — the first candidate's allocated instance is discarded if decode fails
- No caching is implemented: the same string decoded 1000 times triggers 1000 reflect.New calls

**Benchmarking guidance**: Profile `decode()` in applications with > 10K permission checks/second. Below that threshold, the overhead is negligible relative to casbin's own evaluation.

---

## 11. Thread Safety Analysis

| Component | Thread-Safe? | Rationale |
|---|---|---|
| `defaultFactory` (write) | ❌ | No sync.Once; concurrent Register calls race |
| `defaultFactory` (read) | ✅ | After startup, only reads occur; safe by Go memory model (happens-before via goroutine start) |
| `builtinRegister` fields | ✅ | Candidate slices are append-only at registration, then read-only |
| `decode()` | ✅ | Creates fresh instances per call; no shared mutable state |
| `newByE` / `newByT` | ✅ | Pure functions using reflect; no global state |

**Summary**: Thread-safe in normal usage (Register once at startup, then concurrent reads). Unsafe only if Register is called concurrently or after service start.

---

## 12. Extension Points

1. **Custom candidates**: Applications can implement multiple Object types that decode different string formats; Register appends NamedObject automatically, but the pattern generalizes
2. **Alternative Factory implementation**: The `Factory` interface allows replacing `builtinRegister` entirely (e.g., for testing or multi-tenant scenarios)
3. **MetaDB swapping**: By implementing a different `MetadataDB()` return, non-GORM storage backends could be supported

---

## 13. Known Limitations

| # | Limitation | Impact | Severity |
|---|---|---|---|
| 1 | Singleton — no multi-registry support | Cannot test in parallel with different type sets | Medium |
| 2 | No de-registration or replacement | Must restart process to change types | Low |
| 3 | Reflect allocation per decode | GC pressure in ultra-high-throughput scenarios | Low |
| 4 | Hardcoded NamedObject append | Cannot opt out of NamedObject as candidate | Low |
| 5 | Generic error message | `"no register factory for %v"` loses per-candidate failure reasons | Low |
| 6 | No validation of Encode/Decode roundtrip | Register doesn't verify that Encode∘Decode = identity | Low |

---

## 12. Future Opportunities

### Short-term (v0.3.x)

- **Error wrapping**: Collect per-candidate decode errors and wrap with `errors.Join` for better diagnostics
- **sync.Once protection**: Guard `defaultFactory` write to prevent accidental double-registration panics
- **Decode cache**: `sync.Map[string, T]` for repeated decodes of the same string (opt-in, size-bounded)

### Medium-term (v0.4.x / casbin v3 migration)

- **Functional options for Register**: `Register[...](WithExtraCandidates(...))` to avoid hardcoded NamedObject
- **Parameterized Factory**: `New(options, WithFactory(f))` to support testing and multi-tenant use
- **Benchmark suite**: Establish decode throughput baselines before casbin/v3 migration

### Long-term

- **Code generation**: Replace reflect with generated constructors for zero-allocation decode
- **Type registry as first-class resource**: Enable runtime inspection (list registered types, their schemas, etc.)

---

## Appendix: File Map

| File | Role |
|---|---|
| `register.go` | Factory interface, Register function, builtinRegister, decode, newByT, newByE |
| `schema.go` | codeInterface, idInterface, User/Role/Object/Domain interface definitions |
| `metadata_db.go` | builtinMetadataDB implementation (consumes Factory types) |
| `server.go` | Service layer consuming DefaultFactory() |
| `casbin.go` | Enforcer consuming Factory for policy decode |
