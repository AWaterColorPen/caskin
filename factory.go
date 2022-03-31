package caskin

import (
	"reflect"
)

type Factory interface {
	User(string) (User, error)
	Role(string) (Role, error)
	Object(string) (Object, error)
	Domain(string) (Domain, error)
}

func createByE[E any](e E) E {
	v := reflect.ValueOf(e)
	if v.Kind() != reflect.Pointer {
		return *new(E)
	}
	k := reflect.Indirect(v)
	return reflect.New(k.Type()).Interface().(E)
}

func createByT[T any]() T {
	t := *new(T)
	v := reflect.ValueOf(t)
	if v.Kind() != reflect.Pointer {
		return t
	}
	return reflect.New(v.Type().Elem()).Interface().(T)
}

type RoleFactory = func() Role
type ObjectFactory = func() Object
type DomainFactory = func() Domain
