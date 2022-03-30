package caskin

type treeNodeEntryDeleter struct {
	visited  map[any]bool
	children TreeNodeEntryChildrenGetFunc
	delete   TreeNodeEntryDeleteFunc
}

func (t *treeNodeEntryDeleter) Run(current TreeNodeEntry, domain Domain) error {
	if _, ok := t.visited[current.GetID()]; ok {
		return nil
	}

	children := t.children(current, domain)
	for _, v := range children {
		if err := t.Run(v, domain); err != nil {
			return err
		}
	}

	return t.delete(current, domain)
}

func NewTreeNodeEntryDeleter(children TreeNodeEntryChildrenGetFunc, delete TreeNodeEntryDeleteFunc) *treeNodeEntryDeleter {
	return &treeNodeEntryDeleter{
		visited:  map[any]bool{},
		children: children,
		delete:   delete,
	}
}

type TreeNodeEntryChildrenGetFunc = func(TreeNodeEntry, Domain) []TreeNodeEntry
type TreeNodeEntryDeleteFunc = func(TreeNodeEntry, Domain) error
