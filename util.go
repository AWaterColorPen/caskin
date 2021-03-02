package caskin

import (
	"encoding/json"
	"github.com/ahmetb/go-linq/v3"
)

// Filter filter source permission by u, d, action
func Filter(e ienforcer, u User, d Domain, action Action, source interface{}) interface{} {
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e, u, d, action, v.(entry))
	}).ToSlice(&source)
	return source
}

// Filter check entry permission by u, d, action
func Check(e ienforcer, u User, d Domain, action Action, one entry) bool {
	if !one.IsObject() {
		return true
	}

	o := one.GetObject()
	ok, _ := e.Enforce(u, o, d, action)
	return ok
}

// Diff diff source, target list to get add, remove list
func Diff(source, target []interface{}) (add, remove []interface{}) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

// DiffPolicy diff policy source, target list to get add, remove list
func DiffPolicy(source, target []*Policy) (add, remove []*Policy) {
	var s, t []interface{}
	for _, v := range source {
		b, _ := json.Marshal(v)
		s = append(s, b)
	}
	for _, v := range target {
		b, _ := json.Marshal(v)
		t = append(t, b)
	}

	// get diff to add and remove
	a, r := Diff(s, t)
	for _, v := range a {
		p := &Policy{}
		_ = json.Unmarshal(v.([]byte), p)
		add = append(add, p)
	}
	for _, v := range r {
		p := &Policy{}
		_ = json.Unmarshal(v.([]byte), p)
		remove = append(remove, p)
	}
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
	return getIDList(u)
}

type Roles = []Role

func (r Roles) ID() []uint64 {
	return getIDList(r)
}

type Objects = []Object

func (o Objects) ID() []uint64 {
	return getIDList(o)
}

type Domains = []Domain

func (d Domains) ID() []uint64 {
	return getIDList(d)
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

func getTree(source interface{}) map[uint64]uint64 {
	m := map[uint64]uint64{}
	linq.From(source).Where(func(v interface{}) bool {
		_, ok := v.(parentEntry)
		return ok
	}).ForEach(func(v interface{}) {
		u := v.(parentEntry)
		if u.GetParentID() != 0 {
			m[u.GetID()] = u.GetParentID()
		}
	})
	return m
}
