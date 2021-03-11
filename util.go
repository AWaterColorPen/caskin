package caskin

import (
	"github.com/ahmetb/go-linq/v3"
)

// Filter filter source permission by u, d, action
func Filter(e ienforcer, u User, d Domain, action Action, source interface{}) []interface{} {
	var result []interface{}
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result
}

// Check check entry permission by u, d, action
func Check(e ienforcer, u User, d Domain, one ObjectData, action Action) bool {
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
	sourceMap := make(map[interface{}]*Policy)
	targetMap := make(map[interface{}]*Policy)
	var s, t []interface{}
	for _, v := range source {
		key := v.Key()
		s = append(s, key)
		sourceMap[key] = v
	}
	for _, v := range target {
		key := v.Key()
		t = append(t, key)
		targetMap[key] = v
	}

	// get diff to add and remove
	a, r := Diff(s, t)
	for _, v := range a {
		if p, ok := sourceMap[v]; ok {
			add = append(add, p)
		}
		if p, ok := targetMap[v]; ok {
			add = append(add, p)
		}
	}
	for _, v := range r {
		if p, ok := sourceMap[v]; ok {
			remove = append(remove, p)
		}
		if p, ok := targetMap[v]; ok {
			remove = append(remove, p)
		}
	}
	return
}

func isValidFamily(data1, data2 ObjectData, take func(interface{}) error) error {
	o1 := data1.GetObject()
	o2 := data2.GetObject()
	if err := take(o1); err != nil {
		return ErrNotValidParentObject
	}
	if err := take(o2); err != nil {
		return ErrNotValidParentObject
	}
	if o1.GetObjectType() != o2.GetObjectType() {
		return ErrNotValidParentObject
	}
	return nil
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
