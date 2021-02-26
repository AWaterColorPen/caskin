package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// Filter filter source permission by u, d, action
func Filter(e ienforcer, u User, d Domain, action Action, fn func() Object, source interface{}) interface{} {
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e, u, d, action, fn, v.(entry))
	}).ToSlice(&source)
	return source
}

// Filter check entry permission by u, d, action
func Check(e ienforcer, u User, d Domain, action Action, fn func() Object, one entry) bool {
	if !one.IsObject() {
		return true
	}

	o := fn()
	_ = o.Decode(one.GetObject())
	ok, _ := e.Enforce(u, o, d, action)
	return ok
}

// Diff diff source, target list to get add, remove list
func Diff(source, target []interface{}) (add, remove []interface{}) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

func isValid(e entry) error {
	if e == nil {
		return ErrNil
	}

	if e.GetID() == 0 {
		return ErrEmptyID
	}

	return nil
}


type Users []User

func (u Users) ID() []uint64 {
	var id []uint64
	for _, v := range u {
		id = append(id, v.GetID())
	}
	return id
}

type Roles []Role

func (r Roles) ID() []uint64 {
	var id []uint64
	for _, v := range r {
		id = append(id, v.GetID())
	}
	return id
}

type Objects []Object

func (o Objects) ID() []uint64 {
	var id []uint64
	for _, v := range o {
		id = append(id, v.GetID())
	}
	return id
}

type Domains []Domain

func (d Domains) ID() []uint64 {
	var id []uint64
	for _, v := range d {
		id = append(id, v.GetID())
	}
	return id
}

func getIDList(source interface{}) []uint64 {
	var id []uint64
	linq.From(source).Where(func(v interface{}) bool {
		_, ok := v.(entry)
		return ok
	}).Select(func(v interface{}) interface{} {
		return v.(entry).GetID()
	}).ToSlice(&id)
	return id
}


func getIDMap(source interface{}) map[uint64]entry {
	m := map[uint64]entry{}
	linq.From(source).Where(func(v interface{}) bool {
		_, ok := v.(entry)
		return ok
	}).ForEach(func(v interface{}) {
		u := v.(entry)
		m[u.GetID()] = u
	})
	return m
}