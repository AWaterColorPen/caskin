package caskin

type objectDirectory struct {
	Tree  map[uint64][]*Directory
	Node  map[uint64]*Directory
	visit map[uint64]bool
}

func (o *objectDirectory) Search(prefix uint64, searchType DirectorySearchType) []*Directory {
	var out []*Directory
	if prefix == 0 {
		for _, v := range o.Node {
			if _, ok := o.Node[v.GetParentID()]; !ok {
				out = append(out, v)
			}
		}
	} else {
		for _, v := range o.Tree[prefix] {
			out = append(out, v)
		}
	}

	switch searchType {
	case DirectorySearchAll:
		for i := 0; i < len(out); i++ {
			node := out[i]
			for _, v := range o.Tree[node.GetID()] {
				out = append(out, v)
			}
		}
	default:
	}
	return out
}

func (o *objectDirectory) dfs(current *Directory) {
	o.visit[current.GetID()] = true
	current.TopDirectoryCount = uint64(len(o.Tree[current.GetID()]))
	current.AllDirectoryCount = current.TopDirectoryCount
	current.AllItemCount = current.TopItemCount
	for _, child := range o.Tree[current.GetID()] {
		if !o.visit[child.GetID()] {
			o.dfs(child)
		}
		current.AllDirectoryCount += child.AllDirectoryCount
		current.AllItemCount += child.AllItemCount
	}
}

func (o *objectDirectory) build() {
	for _, v := range o.Node {
		if _, ok := o.Node[v.GetParentID()]; !ok {
			o.dfs(v)
		}
	}
}

func NewObjectDirectory(in []*Directory) *objectDirectory {
	t := &objectDirectory{
		Tree:  map[uint64][]*Directory{},
		Node:  IDMap(in),
		visit: map[uint64]bool{},
	}
	for _, v := range in {
		t.Tree[v.GetParentID()] = append(t.Tree[v.GetParentID()], v)
	}
	t.build()
	return t
}
