package caskin

import (
	"github.com/ahmetb/go-linq/v3"
	"golang.org/x/exp/constraints"
)

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
func Filter[T any](e IEnforcer, u User, d Domain, action Action, source []T) []T {
	var result []T
	for _, v := range source {
		if Check(e, u, d, v, action) {
			result = append(result, v)
		}
	}
	return result
}

// Check object/object_data permission by u, d, action
func Check[T any](e IEnforcer, u User, d Domain, one T, action Action) bool {
	if o, ok := any(one).(ObjectData); ok {
		ok, _ = e.Enforce(u, o.GetObject(), d, action)
		return ok
	}
	if o, ok := any(one).(Object); ok {
		ok, _ = e.Enforce(u, o, d, action)
		return ok
	}
	return false
}

// Diff do diff source, target list to get add, remove list
func Diff(source, target []any) (add, remove []any) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

func Diff2[T constraints.Ordered](source, target []T) (add, remove []T) {
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
