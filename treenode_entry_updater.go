package caskin

type treeNodeUpdater[T treeNode] struct {
	parentGet TreeNodeParentGetFunc
	parentAdd TreeNodeParentAddFunc
	parentDel TreeNodeParentDelFunc
}

func (t *treeNodeUpdater[T]) Run(item treeNode, domain Domain) error {
	var source, target []any
	if item.GetParentID() != 0 {
		target = append(target, item.GetParentID())
	}
	parents := t.parentGet(item, domain)
	for _, v := range parents {
		source = append(source, v.GetID())
	}

	add, remove := Diff(source, target)
	for _, v := range add {
		parent := createByT[T]()
		parent.SetID(v.(uint64))
		if err := t.parentAdd(item, parent, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		parent := createByT[T]()
		parent.SetID(v.(uint64))
		if err := t.parentDel(item, parent, domain); err != nil {
			return err
		}
	}
	return nil
}

func NewTreeNodeUpdater[T treeNode](
	parentGet TreeNodeParentGetFunc,
	parentAdd TreeNodeParentAddFunc,
	parentDel TreeNodeParentDelFunc) *treeNodeUpdater[T] {
	return &treeNodeUpdater[T]{
		parentGet: parentGet,
		parentAdd: parentAdd,
		parentDel: parentDel,
	}
}

type TreeNodeParentGetFunc = func(treeNode, Domain) []treeNode
type TreeNodeParentAddFunc = func(treeNode, treeNode, Domain) error
type TreeNodeParentDelFunc = func(treeNode, treeNode, Domain) error
