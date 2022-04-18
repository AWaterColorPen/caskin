package caskin

import (
	"encoding/json"
	"sort"

	"github.com/ahmetb/go-linq/v3"
	"golang.org/x/exp/constraints"
	"golang.org/x/exp/slices"
)

// InheritanceEdge x is node, y is adjacency
type InheritanceEdge[T constraints.Ordered] struct {
	U T `json:"u"`
	V T `json:"v"`
}

func (i *InheritanceEdge[T]) Encode(u, v T) string {
	i.U, i.V = u, v
	b, _ := json.Marshal(i)
	return string(b)
}

func (i *InheritanceEdge[T]) Decode(in string) error {
	return json.Unmarshal([]byte(in), i)
}

type EdgeSorter[T constraints.Ordered] map[T]int

func (e EdgeSorter[T]) RootFirstSort(edges []*InheritanceEdge[T]) {
	sort.Slice(edges, func(i, j int) bool {
		return e[edges[i].U] < e[edges[j].U]
	})
}

func (e EdgeSorter[T]) LeafFirstSort(edges []*InheritanceEdge[T]) {
	sort.Slice(edges, func(i, j int) bool {
		return e[edges[i].V] > e[edges[j].V]
	})
}

func NewEdgeSorter[T constraints.Ordered](order []T) EdgeSorter[T] {
	sorter := map[T]int{}
	for i, v := range order {
		sorter[v] = i
	}
	return sorter
}

type InheritanceGraph[T constraints.Ordered] map[T][]T

func (g InheritanceGraph[T]) Sort() InheritanceGraph[T] {
	var keys []T
	for k := range g {
		keys = append(keys, k)
	}

	slices.Sort(keys)
	m := InheritanceGraph[T]{}
	for _, k := range keys {
		m[k] = g[k]
		slices.Sort(m[k])
	}
	return m
}

func (g InheritanceGraph[T]) TopSort() []T {
	inDegree := map[T]int{}
	for k := range g {
		inDegree[k] = 0
	}
	for _, node := range g {
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
		for _, v := range g[node] {
			inDegree[v]--
			if inDegree[v] == 0 {
				queue = append(queue, v)
			}
		}
	}
	return queue
}

func MergeInheritanceGraph[T constraints.Ordered](graphs ...InheritanceGraph[T]) InheritanceGraph[T] {
	m := InheritanceGraph[T]{}
	for _, graph := range graphs {
		for node, adjacency := range graph {
			if _, ok := m[node]; !ok {
				m[node] = []T{}
			}
			for _, v := range adjacency {
				m[node] = append(m[node], v)
			}
		}
	}

	for node, adjacency := range m {
		var t []T
		linq.From(adjacency).Distinct().ToSlice(&t)
		m[node] = t
	}
	return m.Sort()
}
