package caskin

type treeNodeDeleter struct {
	visited  map[any]bool
	children TreeNodeChildrenGetFunc
	delete   TreeNodeDeleteFunc
}

func (t *treeNodeDeleter) Run(current treeNode, domain Domain) error {
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

func NewTreeNodeDeleter(children TreeNodeChildrenGetFunc, delete TreeNodeDeleteFunc) *treeNodeDeleter {
	return &treeNodeDeleter{
		visited:  map[any]bool{},
		children: children,
		delete:   delete,
	}
}

type TreeNodeChildrenGetFunc = func(treeNode, Domain) []treeNode
type TreeNodeDeleteFunc = func(treeNode, Domain) error
