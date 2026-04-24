# Contributing to Caskin

Thank you for your interest in contributing to caskin! This guide covers
everything you need to get started: setting up your environment, the contribution
workflow, coding conventions, and how to run the tests.

## Table of Contents

- [Prerequisites](#prerequisites)
- [Setting Up the Development Environment](#setting-up-the-development-environment)
- [Contribution Workflow](#contribution-workflow)
- [Coding Conventions](#coding-conventions)
- [Running the Tests](#running-the-tests)
- [Adding a New Feature](#adding-a-new-feature)
- [Common Pitfalls](#common-pitfalls)
- [Commit Message Format](#commit-message-format)
- [Reporting Issues](#reporting-issues)

---

## Prerequisites

- **Go 1.24+** — caskin targets the latest stable Go release
- **Git**
- **SQLite** (available on most systems; used by default in tests)

Optional (for integration tests):
- Docker (for MySQL / PostgreSQL tests via test containers, if added)
- Redis (for watcher tests)

---

## Setting Up the Development Environment

```bash
# 1. Fork the repository on GitHub and clone your fork
git clone https://github.com/<your-username>/caskin.git
cd caskin

# 2. Install dependencies
go mod download

# 3. Verify everything builds and tests pass
go build ./...
go test ./...
```

All tests use in-memory SQLite via `t.TempDir()`, so no external services are required
for the default test suite.

---

## Contribution Workflow

caskin follows a **feature-branch + pull request** workflow:

1. **Create a feature branch** from `main`:
   ```bash
   git checkout -b feature/your-feature-name
   ```

2. **Make your changes** in small, focused commits (see
   [Commit Message Format](#commit-message-format)).

3. **Run the full test suite** before pushing:
   ```bash
   go test ./...
   ```

4. **Push your branch** and open a Pull Request against `main`.

5. **Address review feedback** — push additional commits to the same branch.
   The PR will be merged once approved.

### Branch Naming

| Type of change | Branch prefix |
|----------------|---------------|
| New feature | `feature/` |
| Bug fix | `fix/` |
| Documentation | `docs/` |
| Refactoring | `refactor/` |
| Dependency updates | `deps/` |

---

## Coding Conventions

### Go Style

- Follow the [Effective Go](https://go.dev/doc/effective_go) guidelines and standard
  Go formatting (`gofmt` / `goimports`).
- All exported symbols must have godoc comments. Follow the existing style in `schema.go`
  and `service.go` as a reference.
- Use `slices` and `cmp` from the standard library rather than third-party utilities
  (caskin removed `go-linq` and `golang.org/x/exp` in Phase 1).

### Interface Contracts

When implementing or extending caskin's core interfaces, keep these rules in mind:

- **`caskin.Register[U, R, O, D]()`** must always be called before `caskin.New`.
  Every test setup function must include this call.
- **`GetObjectType()`** is the correct method name on `caskin.Object`; not `GetType()`.
- **`caskin.New`** takes a `*caskin.Options` struct — there are no functional-option
  (`WithXxx`) variants.
- **`IService`** does not expose `GetEnforcer()`. Use service-level methods
  (`GetObject`, `GetCurrentBackend`, etc.) or the `caskin.Check` / `caskin.Filter`
  helpers when you need permission checks.

### Error Handling

- Return errors from all service methods; never swallow them silently.
- Use the sentinel errors in `error.go` where applicable (e.g., `caskin.ErrNotFound`).
- New error types should be added to `error.go` and follow the existing naming pattern.

### No Breaking Changes Without a Discussion

- Changes to any exported interface (`IService`, `User`, `Role`, `Object`, `Domain`,
  `Factory`) are **breaking changes** and require prior discussion in a GitHub issue.
- Internal struct fields and unexported methods can be changed freely.
- New methods on interfaces must be added in a backward-compatible way (new optional
  interface, separate method on `server`, etc.) unless a major version bump is planned.

---

## Running the Tests

### Full suite (SQLite, no external dependencies)

```bash
go test ./...
```

### Verbose output

```bash
go test -v ./...
```

### A single package

```bash
go test -v github.com/awatercolorpen/caskin/playground
```

### With race detection

```bash
go test -race ./...
```

### Test coverage

```bash
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

### Using the playground in your own tests

The `playground` package provides a fully bootstrapped in-memory environment:

```go
import "github.com/awatercolorpen/caskin/playground"

func TestSomething(t *testing.T) {
    playground.DictionaryDsn = "../configs/caskin.toml"
    stage, err := playground.NewPlaygroundWithSqlitePath(t.TempDir())
    if err != nil {
        t.Fatal(err)
    }
    // stage.Service    — IService
    // stage.Superadmin — *example.User (superadmin)
    // stage.Admin      — *example.User (admin in stage.Domain)
    // stage.Member     — *example.User (member in stage.Domain)
    // stage.Domain     — *example.Domain
}
```

---

## Adding a New Feature

Before writing code, please:

1. **Open a GitHub issue** describing the feature and your proposed approach.
2. Wait for a maintainer to acknowledge the issue (or express concerns).
3. If the feature touches a core interface, discuss the API design in the issue.

For the implementation:

- Add or update the relevant interface in `service.go` if needed.
- Implement the feature in a focused source file (e.g., `server_your_feature.go`).
- Add tests alongside the implementation (e.g., `server_your_feature_test.go`).
- Update the dictionary model (`dictionary_model.go`) if the feature requires new
  TOML config keys.
- Update `docs/api-reference.md` with the new methods.
- Update `docs/use-cases.md` if the feature enables a common usage pattern.

### Phase gating

caskin uses a phased modernization roadmap
(`docs/superpowers/specs/2026-03-10-caskin-modernization.md`).
Large changes should be proposed as a Phase item and coordinated with maintainers
before implementation.

---

## Common Pitfalls

These are the most common mistakes found during code review. Check them before
submitting a PR:

1. **Missing `caskin.Register[...]()` call** — any test setup function or `main`
   that calls `caskin.New` must call `caskin.Register` first, or you'll get a
   runtime panic.

2. **Using `GetType()` instead of `GetObjectType()`** — the `caskin.Object` interface
   uses `GetObjectType()`. Check `schema.go` when unsure.

3. **Functional options in examples** — `caskin.New` does not have `WithDB` or
   `WithDictionary` helpers. Always use the `*Options` struct.

4. **Calling `IService.GetEnforcer()`** — `IService` does not expose the enforcer.
   Use service-level helpers or `caskin.Check`.

5. **Assuming domain-isolation doesn't apply** — roles, objects, and policies are
   always scoped to a domain. There is no "global" role except superadmin.

---

## Commit Message Format

caskin uses [Conventional Commits](https://www.conventionalcommits.org/):

```
<type>(<scope>): <short summary>

[optional body]
[optional footer]
```

Common types:

| Type | When to use |
|------|-------------|
| `feat` | New feature |
| `fix` | Bug fix |
| `docs` | Documentation only |
| `refactor` | Code change that is neither a feat nor a fix |
| `test` | Adding or correcting tests |
| `chore` | Build, CI, dependency updates |

Examples:

```
feat(directory): add GetObjectHierarchyLevel method
fix(watcher): handle nil WatcherOption gracefully
docs: add use-cases guide with multi-domain examples
refactor: replace sort.Slice with slices.SortFunc
```

---

## Reporting Issues

- **Bug reports** — include the Go version, caskin version, a minimal reproducer,
  and the observed vs. expected behavior.
- **Feature requests** — describe the problem you're trying to solve and why
  existing APIs don't cover it.
- **Security issues** — please do **not** open a public issue. Email the maintainer
  directly (see the GitHub profile).
