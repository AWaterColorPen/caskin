package caskin

type parentEntryDeleter struct {
	visited    map[interface{}]bool
	childrenFn childrenFn
	deleteFn   deleteFn
}

func (p *parentEntryDeleter) dfs(current treeNodeEntry, domain Domain) error {
	if _, ok := p.visited[current.GetID()]; ok {
		return nil
	}

	children := p.childrenFn(current, domain)
	for _, v := range children {
		if err := p.dfs(v, domain); err != nil {
			return err
		}
	}

	return p.deleteFn(current, domain)
}

func newParentEntryDeleter(childrenFn childrenFn, deleteFn deleteFn) *parentEntryDeleter {
	return &parentEntryDeleter{
		visited:    map[interface{}]bool{},
		childrenFn: childrenFn,
		deleteFn:   deleteFn,
	}
}

type childrenFn = func(treeNodeEntry, Domain) []treeNodeEntry
type deleteFn = func(treeNodeEntry, Domain) error
