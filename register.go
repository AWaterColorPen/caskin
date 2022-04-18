package caskin

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

var defaultFactory Factory

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

func Register[U User, R Role, O Object, D Domain]() {
	b := &builtinRegister[U, R, O, D]{}
	b.user = append(b.user, b.NewUser())
	b.role = append(b.role, b.NewRole())
	b.object = append(b.object, b.NewObject())
	b.domain = append(b.domain, b.NewDomain())
	// builtin
	b.object = append(b.object, &NamedObject{})
	defaultFactory = b
}

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
