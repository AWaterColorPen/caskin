package caskin

import (
	"sort"

	"github.com/ahmetb/go-linq/v3"
)

func SortedInheritanceRelations(relations InheritanceRelations) InheritanceRelations {
	var keys []uint64
	for k := range relations {
		keys = append(keys, k)
	}

	sort.SliceStable(keys, func(i, j int) bool {
		return keys[i] < keys[j]
	})

	m := InheritanceRelations{}
	for _, k := range keys {
		m[k] = relations[k]
		m[k] = SortedInheritanceRelation(m[k])
	}
	return m
}

func SortedInheritanceRelation(relation InheritanceRelation) InheritanceRelation {
	sort.SliceStable(relation, func(i, j int) bool {
		return relation[i] < relation[j]
	})
	return relation
}

// Filter filter source permission by u, d, action
func Filter(e IEnforcer, u User, d Domain, action Action, source interface{}) []interface{} {
	var result []interface{}
	linq.From(source).Where(func(v interface{}) bool {
		return Check(e, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result
}

// Check check entry permission by u, d, action
func Check(e IEnforcer, u User, d Domain, one ObjectData, action Action) bool {
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
	for _, v := range source {
		sourceMap[v.Key()] = v
	}
	for _, v := range target {
		targetMap[v.Key()] = v
	}

	for _, v := range source {
		if _, ok := targetMap[v.Key()]; !ok {
			remove = append(remove, v)
		}
	}
	for _, v := range target {
		if _, ok := sourceMap[v.Key()]; !ok {
			add = append(add, v)
		}

	}
	return
}

func isValidFamily(data1, data2 ObjectData, take func(interface{}) error) error {
	o1 := data1.GetObject()
	o2 := data2.GetObject()
	if err := take(o1); err != nil {
		return ErrInValidParentObject
	}
	if err := take(o2); err != nil {
		return ErrInValidParentObject
	}
	if o1.GetObjectType() != o2.GetObjectType() {
		return ErrInValidParentObject
	}
	return nil
}

func isValid(item idInterface) error {
	if item == nil {
		return ErrNil
	}

	if item.GetID() == 0 {
		return ErrEmptyID
	}

	return nil
}

func isRoot(node treeNode) bool {
	return node.GetParentID() == 0
}

func isObjectTypeObjectIDBeSelfIDCheck(object Object) error {
	if object.GetObjectType() == ObjectTypeObject &&
		object.GetObject().GetID() != object.GetID() {
		return ErrObjectTypeObjectIDMustBeItselfID
	}
	return nil
}
