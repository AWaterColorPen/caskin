package caskin

import "fmt"

type objectDirectory struct {
	Tree     map[uint64][]*Directory
	Node     map[uint64]*Directory
	MaxDepth uint64

	visit map[uint64]bool
}

func (t *objectDirectory) GetNode(id uint64) (*Directory, error) {
	node, ok := t.Node[id]
	if !ok {
		return nil, fmt.Errorf("found no directory %v", id)
	}
	return node, nil
}

func (t *objectDirectory) Search(prefix uint64, searchType DirectorySearchType) []*Directory {
	var out []*Directory
	for _, v := range t.Tree[prefix] {
		out = append(out, v)
	}
	switch searchType {
	case DirectorySearchAll:
		for i := 0; i < len(out); i++ {
			node := out[i]
			for _, v := range t.Tree[node.GetID()] {
				out = append(out, v)
			}
		}
	default:
	}
	return out
}

func (t *objectDirectory) SubTree(prefix uint64) (out []*Directory) {
	node, err := t.GetNode(prefix)
	if err != nil {
		return nil
	}
	out = append(out, node)
	children := t.Search(prefix, DirectorySearchAll)
	out = append(out, children...)
	return out
}

func (t *objectDirectory) CheckDepth(from, to uint64) error {
	node1, err := t.GetNode(from)
	if err != nil {
		return err
	}
	node2, err := t.GetNode(to)
	if err != nil {
		return err
	}
	if node1.Distance+node2.Depth >= DirectoryMaxDepth {
		return fmt.Errorf("max depth is too large")
	}
	return nil
}

func (t *objectDirectory) cleanVisit() {
	t.visit = map[uint64]bool{}
}

func (t *objectDirectory) dfs(current *Directory, depth uint64) error {
	if t.visit[current.GetID()] {
		return ErrParentToDescendant
	}
	t.visit[current.GetID()] = true

	current.Depth = depth
	current.Distance = 1
	current.TopDirectoryCount = uint64(len(t.Tree[current.GetID()]))
	current.AllDirectoryCount = current.TopDirectoryCount
	current.AllItemCount = current.TopItemCount
	for _, child := range t.Tree[current.GetID()] {
		if err := t.dfs(child, depth+1); err != nil {
			return err
		}

		current.AllDirectoryCount += child.AllDirectoryCount
		current.AllItemCount += child.AllItemCount
		if current.Distance < child.Distance+1 {
			current.Distance = child.Distance + 1
		}
	}
	return nil
}

func (t *objectDirectory) build() error {
	t.cleanVisit()
	inDegree := map[uint64]int{}
	for k := range t.Node {
		inDegree[k] = 0
	}
	for k, node := range t.Tree {
		if k == 0 {
			continue
		}
		for _, v := range node {
			inDegree[v.GetID()]++
		}
	}

	var queue []uint64
	for k, v := range inDegree {
		if v == 0 {
			queue = append(queue, k)
		}
	}

	for _, v := range queue {
		current := t.Node[v]
		if err := t.dfs(current, 1); err != nil {
			return err
		}
	}

	t.MaxDepth = uint64(0)
	for _, v := range t.Node {
		if t.MaxDepth < v.Depth {
			t.MaxDepth = v.Depth
		}
	}
	return nil
}

func NewObjectDirectory(in []*Directory) (*objectDirectory, error) {
	t := &objectDirectory{
		Tree: map[uint64][]*Directory{},
		Node: IDMap(in),
	}
	for _, v := range in {
		t.Tree[v.GetParentID()] = append(t.Tree[v.GetParentID()], v)
	}
	return t, t.build()
}
