package caskin

import "sort"

// InheritanceRelation value is sons' id list
type InheritanceRelation = []uint64

// InheritanceRelations key is parent id, value is sons' id list
type InheritanceRelations = map[uint64]InheritanceRelation

// InheritanceEdge x is node, y is adjacency
type InheritanceEdge struct {
	X, Y uint64
}

type EdgeSorter struct {
	mIndex map[uint64]uint64
}

func (e *EdgeSorter) RootFirstSort(edges []*InheritanceEdge) {
	sort.Slice(edges, func(i, j int) bool {
		return e.mIndex[edges[i].X] < e.mIndex[edges[j].X]
	})
}

func (e *EdgeSorter) LeafFirstSort(edges []*InheritanceEdge) {
	sort.Slice(edges, func(i, j int) bool {
		return e.mIndex[edges[i].Y] > e.mIndex[edges[j].Y]
	})
}

func NewEdgeSorter(order []uint64) *EdgeSorter {
	mIndex := map[uint64]uint64{}
	for i, v := range order {
		mIndex[v] = uint64(i)
	}

	return &EdgeSorter{
		mIndex: mIndex,
	}
}
