package caskin

import (
	"cmp"
	"encoding/json"
	"slices"
)

// InheritanceEdge represents a directed edge in an inheritance graph where
// U is the parent node and V is the child node. It is used to serialise
// role/object hierarchy edges for storage and comparison.
type InheritanceEdge[T cmp.Ordered] struct {
	U T `json:"u"`
	V T `json:"v"`
}

// Encode stores the (u, v) pair and returns the JSON string representation.
func (i *InheritanceEdge[T]) Encode(u, v T) string {
	i.U, i.V = u, v
	b, _ := json.Marshal(i)
	return string(b)
}

// Decode parses a JSON string produced by [InheritanceEdge.Encode] back into
// the edge struct.
func (i *InheritanceEdge[T]) Decode(in string) error {
	return json.Unmarshal([]byte(in), i)
}

// EdgeSorter maps node values to their topological sort order. It is used to
// sort slices of [InheritanceEdge] so that root nodes come first
// ([EdgeSorter.RootFirstSort]) or leaf nodes come first
// ([EdgeSorter.LeafFirstSort]).
type EdgeSorter[T cmp.Ordered] map[T]int

// RootFirstSort sorts edges so that edges whose source (U) is closer to the
// root of the graph come first. This is useful for processing parent nodes
// before their children.
func (e EdgeSorter[T]) RootFirstSort(edges []*InheritanceEdge[T]) {
	slices.SortFunc(edges, func(a, b *InheritanceEdge[T]) int {
		return cmp.Compare(e[a.U], e[b.U])
	})
}

// LeafFirstSort sorts edges so that edges whose destination (V) is closer to
// the leaves of the graph come first. This is useful for processing children
// before their parents (e.g. when deleting a subtree).
func (e EdgeSorter[T]) LeafFirstSort(edges []*InheritanceEdge[T]) {
	slices.SortFunc(edges, func(a, b *InheritanceEdge[T]) int {
		return cmp.Compare(e[b.V], e[a.V])
	})
}

// NewEdgeSorter builds an [EdgeSorter] from a topological ordering of node
// values. Nodes that appear earlier in order are considered closer to the root.
func NewEdgeSorter[T cmp.Ordered](order []T) EdgeSorter[T] {
	sorter := map[T]int{}
	for i, v := range order {
		sorter[v] = i
	}
	return sorter
}

// InheritanceGraph is an adjacency list representation of a directed graph
// where each key is a node and the associated slice contains its direct
// children (nodes it inherits into).
type InheritanceGraph[T cmp.Ordered] map[T][]T

// Sort returns a new [InheritanceGraph] where the keys and each adjacency list
// are sorted, making the graph representation deterministic.
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

// TopSort performs a topological sort (Kahn's algorithm / BFS) on the graph
// and returns the nodes in an order where every parent appears before all of
// its children. The graph must be a DAG; cycles will cause nodes to be omitted
// from the result.
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
	// BFS topological sort: queue grows as in-degrees reach zero
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

// distinct returns a new slice with duplicate elements removed, preserving order.
func distinct[T cmp.Ordered](s []T) []T {
	seen := make(map[T]struct{}, len(s))
	out := make([]T, 0, len(s))
	for _, v := range s {
		if _, ok := seen[v]; !ok {
			seen[v] = struct{}{}
			out = append(out, v)
		}
	}
	return out
}

// MergeInheritanceGraph merges multiple [InheritanceGraph] values into one,
// deduplicating adjacency entries and sorting the result. This is useful when
// combining role/object inheritance rules from multiple sources.
func MergeInheritanceGraph[T cmp.Ordered](graphs ...InheritanceGraph[T]) InheritanceGraph[T] {
	m := InheritanceGraph[T]{}
	for _, graph := range graphs {
		for node, adjacency := range graph {
			if _, ok := m[node]; !ok {
				m[node] = []T{}
			}
			m[node] = append(m[node], adjacency...)
		}
	}

	for node, adjacency := range m {
		m[node] = distinct(adjacency)
	}
	return m.Sort()
}
