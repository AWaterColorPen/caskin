package caskin

import (
	"errors"
	"fmt"
	"reflect"
	"strings"

	"gorm.io/gorm"
)

// defaultFactory holds the singleton Factory set by Register.
var defaultFactory Factory

// Factory is the type registry used by caskin to instantiate and decode
// concrete User, Role, Object, and Domain values from their string
// representations stored in casbin.
//
// Call [Register] once at program startup with your concrete types, then
// use [DefaultFactory] to access the registered factory throughout your
// application.
type Factory interface {
	// User decodes a casbin subject string into a User.
	User(string) (User, error)
	// Role decodes a casbin role string into a Role.
	Role(string) (Role, error)
	// Object decodes a casbin object string into an Object.
	Object(string) (Object, error)
	// Domain decodes a casbin domain string into a Domain.
	Domain(string) (Domain, error)
	// NewUser returns a zero-value User of the registered concrete type.
	NewUser() User
	// NewRole returns a zero-value Role of the registered concrete type.
	NewRole() Role
	// NewObject returns a zero-value Object of the registered concrete type.
	NewObject() Object
	// NewDomain returns a zero-value Domain of the registered concrete type.
	NewDomain() Domain
	// MetadataDB wraps the given GORM database with the registered type
	// information and returns a [MetaDB] implementation.
	MetadataDB(db *gorm.DB) MetaDB
}

// RegisterOption is a functional option for [Register] that customises the
// type factory. Use [WithObject], [WithRole], [WithUser], and [WithDomain]
// to add extra candidate types to the decode chain.
type RegisterOption func(*registerConfig)

// registerConfig holds the additional candidates supplied via RegisterOption.
type registerConfig struct {
	extraUsers   []User
	extraRoles   []Role
	extraObjects []Object
	extraDomains []Domain
	noBuiltins   bool
}

// WithObject adds an extra [Object] candidate to the factory's decode chain.
// This allows registering additional Object implementations (e.g. custom
// object types) without hardcoding them inside [Register].
//
// Example:
//
//	caskin.Register[*MyUser, *MyRole, *MyObject, *MyDomain](
//	    caskin.WithObject(&MySpecialObject{}),
//	)
func WithObject(o Object) RegisterOption {
	return func(c *registerConfig) {
		c.extraObjects = append(c.extraObjects, o)
	}
}

// WithRole adds an extra [Role] candidate to the factory's decode chain.
func WithRole(r Role) RegisterOption {
	return func(c *registerConfig) {
		c.extraRoles = append(c.extraRoles, r)
	}
}

// WithUser adds an extra [User] candidate to the factory's decode chain.
func WithUser(u User) RegisterOption {
	return func(c *registerConfig) {
		c.extraUsers = append(c.extraUsers, u)
	}
}

// WithDomain adds an extra [Domain] candidate to the factory's decode chain.
func WithDomain(d Domain) RegisterOption {
	return func(c *registerConfig) {
		c.extraDomains = append(c.extraDomains, d)
	}
}

// WithoutBuiltins disables the automatic registration of built-in candidates
// (e.g. [NamedObject]). Use this when you want full control over the decode
// chain and don't need the default fallback types.
func WithoutBuiltins() RegisterOption {
	return func(c *registerConfig) {
		c.noBuiltins = true
	}
}

// Register wires up the global type factory with the caller's concrete
// implementations of [User], [Role], [Object], and [Domain]. It must be called
// exactly once before creating any caskin service via [New].
//
// All four type parameters must be pointer types that implement their
// respective interfaces:
//
//	caskin.Register[*MyUser, *MyRole, *MyObject, *MyDomain]()
//
// Optional [RegisterOption] values can be passed to add extra candidate types:
//
//	caskin.Register[*MyUser, *MyRole, *MyObject, *MyDomain](
//	    caskin.WithObject(&MyNamedObject{}),
//	)
func Register[U User, R Role, O Object, D Domain](opts ...RegisterOption) {
	cfg := &registerConfig{}
	for _, opt := range opts {
		opt(cfg)
	}

	b := &builtinRegister[U, R, O, D]{}
	b.user = append(b.user, b.NewUser())
	b.user = append(b.user, cfg.extraUsers...)
	b.role = append(b.role, b.NewRole())
	b.role = append(b.role, cfg.extraRoles...)
	b.object = append(b.object, b.NewObject())
	b.object = append(b.object, cfg.extraObjects...)
	b.domain = append(b.domain, b.NewDomain())
	b.domain = append(b.domain, cfg.extraDomains...)

	// Built-in candidates (e.g. NamedObject as fallback decoder for objects).
	if !cfg.noBuiltins {
		b.object = append(b.object, &NamedObject{})
	}

	defaultFactory = b
}

// DefaultFactory returns the global [Factory] set by the most recent call to
// [Register]. It panics (via nil-pointer) if [Register] has not been called.
func DefaultFactory() Factory {
	return defaultFactory
}

type builtinRegister[U User, R Role, O Object, D Domain] struct {
	user   []User
	role   []Role
	object []Object
	domain []Domain
}

func (b *builtinRegister[U, R, O, D]) User(code string) (User, error) {
	return decode(code, b.user)
}

func (b *builtinRegister[U, R, O, D]) Role(code string) (Role, error) {
	return decode(code, b.role)
}

func (b *builtinRegister[U, R, O, D]) Object(code string) (Object, error) {
	return decode(code, b.object)
}

func (b *builtinRegister[U, R, O, D]) Domain(code string) (Domain, error) {
	return decode(code, b.domain)
}

func (b *builtinRegister[U, R, O, D]) NewUser() User {
	return newByT[U]()
}

func (b *builtinRegister[U, R, O, D]) NewRole() Role {
	return newByT[R]()
}

func (b *builtinRegister[U, R, O, D]) NewObject() Object {
	return newByT[O]()
}

func (b *builtinRegister[U, R, O, D]) NewDomain() Domain {
	return newByT[D]()
}

func (b *builtinRegister[U, R, O, D]) MetadataDB(db *gorm.DB) MetaDB {
	return &builtinMetadataDB[U, R, O, D]{DB: db}
}

// decode tries each candidate by creating a fresh copy and calling Decode.
// The first candidate that decodes without error is returned. If all
// candidates fail, the returned error wraps every individual decode error
// for diagnostic purposes.
func decode[T codeInterface](code string, candidate []T) (T, error) {
	var errs []error
	for _, v := range candidate {
		e := newByE(v)
		if err := e.Decode(code); err == nil {
			return e, nil
		} else {
			errs = append(errs, fmt.Errorf("%T: %w", v, err))
		}
	}
	var zero T
	joinedErr := errors.Join(errs...)
	typeName := "unknown"
	if len(candidate) > 0 {
		typeName = strings.TrimPrefix(fmt.Sprintf("%T", candidate[0]), "*")
	}
	return zero, fmt.Errorf("no registered factory for %v (type %s): %w", code, typeName, joinedErr)
}

// newByE creates a new zero-value instance of the same concrete type as e.
// If e is a pointer, it allocates a new value of the pointed-to type.
func newByE[E any](e E) E {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Pointer {
		return *new(E)
	}
	k := reflect.Indirect(v)
	return reflect.New(k.Type()).Interface().(E)
}

// newByT creates a new zero-value instance of type T.
// If T is a pointer type, it allocates a new value of the pointed-to type.
func newByT[T any]() T {
	t := *new(T)
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Pointer {
		return t
	}
	return reflect.New(v.Type().Elem()).Interface().(T)
}
