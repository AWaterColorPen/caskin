package caskin

import "sort"

// InheritanceRelation value is sons' id list
type InheritanceRelation = []uint64

// InheritanceRelations key is parent id, value is sons' id list
type InheritanceRelations = map[uint64]InheritanceRelation

// InheritanceEdge x is node, y is adjacency
type InheritanceEdge[T comparable] struct {
	U, V T
}

type EdgeSorter[T comparable] map[T]int

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

func NewEdgeSorter[T comparable](order []T) EdgeSorter[T] {
	sorter := map[T]int{}
	for i, v := range order {
		sorter[v] = i
	}
	return sorter
}
