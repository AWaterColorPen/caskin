package caskin

import (
	"fmt"
	"reflect"
)

var defaultRegister = &builtinRegister{}
var defaultFactory = (*builtinFactory)(defaultRegister)

func DefaultRegister() Register {
	return defaultRegister
}

func DefaultFactory() Factory {
	return defaultFactory
}

type Register interface {
	Register(User, Role, Object, Domain)
}

type builtinRegister struct {
	user   []User
	role   []Role
	object []Object
	domain []Domain
}

type Factory interface {
	User(string) (User, error)
	Role(string) (Role, error)
	Object(string) (Object, error)
	Domain(string) (Domain, error)
	NewUser() User
	NewRole() Role
	NewObject() Object
	NewDomain() Domain
}

func (b *builtinRegister) Register(user User, role Role, object Object, domain Domain) {
	b.user = append(b.user, user)
	b.role = append(b.role, role)
	b.object = append(b.object, object)
	b.domain = append(b.domain, domain)
	// builtin
	b.object = append(b.object, &NamedObject{})
}

type builtinFactory builtinRegister

func (b *builtinFactory) User(code string) (User, error) {
	return decode(code, b.user)
}

func (b *builtinFactory) Role(code string) (Role, error) {
	return decode(code, b.role)
}

func (b *builtinFactory) Object(code string) (Object, error) {
	return decode(code, b.object)
}

func (b *builtinFactory) Domain(code string) (Domain, error) {
	return decode(code, b.domain)
}

func (b *builtinFactory) NewUser() User {
	return newByE(b.user[0])
}

func (b *builtinFactory) NewRole() Role {
	return newByE(b.role[0])
}

func (b *builtinFactory) NewObject() Object {
	return newByE(b.object[0])
}

func (b *builtinFactory) NewDomain() Domain {
	return newByE(b.domain[0])
}

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

func newByE[E any](e E) E {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Pointer {
		return *new(E)
	}
	k := reflect.Indirect(v)
	return reflect.New(k.Type()).Interface().(E)
}

func newByT[T any]() T {
	t := *new(T)
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Pointer {
		return t
	}
	return reflect.New(v.Type().Elem()).Interface().(T)
}
