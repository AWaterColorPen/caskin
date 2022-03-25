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

func MergedInheritanceRelations(relations ...InheritanceRelations) InheritanceRelations {
	m := InheritanceRelations{}
	for _, graph := range relations {
		for node, adjacency := range graph {
			if _, ok := m[node]; !ok {
				m[node] = InheritanceRelation{}
			}
			for _, v := range adjacency {
				m[node] = append(m[node], v)
			}
		}
	}

	for node, adjacency := range m {
		t := InheritanceRelation{}
		linq.From(adjacency).Distinct().ToSlice(&t)
		m[node] = t
	}
	return SortedInheritanceRelations(m)
}

// TopSort root first
func TopSort[T comparable](graph map[T][]T) []T {
	inDegree := map[T]int{}
	for k := range graph {
		inDegree[k] = 0
	}
	for _, node := range graph {
		for _, v := range node {
			inDegree[v]++
		}
	}

	var queue []T
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}
	for i := 0; i < len(queue); i++ {
		node := queue[i]
		for _, v := range graph[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return queue
}

// Filter do filter source permission by u, d, action
func Filter(e IEnforcer, u User, d Domain, action Action, source interface{}) []interface{} {
	var result []interface{}
	linq.From(source).Where(func(v interface{}) bool {
		return CheckObjectData(e, u, d, v.(ObjectData), action)
	}).ToSlice(&result)
	return result
}

// CheckObjectData check object_data permission by u, d, action
func CheckObjectData(e IEnforcer, u User, d Domain, one ObjectData, action Action) bool {
	o := one.GetObject()
	ok, _ := e.Enforce(u, o, d, action)
	return ok
}

// CheckObject check object permission by u, d, action
func CheckObject(e IEnforcer, u User, d Domain, one Object, action Action) bool {
	ok, _ := e.Enforce(u, one, d, action)
	return ok
}

// Diff do diff source, target list to get add, remove list
func Diff(source, target []any) (add, remove []any) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

// DiffPolicy diff policy source, target list to get add, remove list
func DiffPolicy(source, target []*Policy) (add, remove []*Policy) {
	sourceMap := make(map[any]*Policy)
	targetMap := make(map[any]*Policy)
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
