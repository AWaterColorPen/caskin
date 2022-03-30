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

func create[T any](in T) T {
	v := reflect.ValueOf(in)
	k := reflect.Indirect(v)
	b := reflect.New(k.Type()).Interface().(T)
	return b
}

type RoleFactory = func() Role
type ObjectFactory = func() Object
type DomainFactory = func() Domain
