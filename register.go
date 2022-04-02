package caskin

import (
	"fmt"
	"reflect"
)

type Register interface {
	Register([]User, []Role, []Object, []Domain)
	Factory() Factory
}

type defaultRegister struct {
	user   []User
	role   []Role
	object []Object
	domain []Domain
}

func (d *defaultRegister) Register(user []User, role []Role, object []Object, domain []Domain) {
	d.user, d.role, d.object, d.domain = user, role, object, domain
}

func (d *defaultRegister) Factory() Factory {
	f := defaultFactory(*d)
	return &f
}

type Factory interface {
	User(string) (User, error)
	Role(string) (Role, error)
	Object(string) (Object, error)
	Domain(string) (Domain, error)
}

type defaultFactory defaultRegister

func (d *defaultFactory) User(code string) (User, error) {
	return decode(code, d.user)
}

func (d *defaultFactory) Role(code string) (Role, error) {
	return decode(code, d.role)
}

func (d *defaultFactory) Object(code string) (Object, error) {
	return decode(code, d.object)
}

func (d *defaultFactory) Domain(code string) (Domain, error) {
	return decode(code, d.domain)
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
