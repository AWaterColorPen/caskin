package caskin

type treeNodeEntryUpdater[T TreeNodeEntry] struct {
	parentGet TreeNodeEntryParentGetFunc
	parentAdd TreeNodeEntryParentAddFunc
	parentDel TreeNodeEntryParentDelFunc
}

func (t *treeNodeEntryUpdater[T]) Run(item TreeNodeEntry, domain Domain) error {
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

func NewTreeNodeEntryUpdater[T TreeNodeEntry](
	parentGet TreeNodeEntryParentGetFunc,
	parentAdd TreeNodeEntryParentAddFunc,
	parentDel TreeNodeEntryParentDelFunc) *treeNodeEntryUpdater[T] {
	return &treeNodeEntryUpdater[T]{
		parentGet: parentGet,
		parentAdd: parentAdd,
		parentDel: parentDel,
	}
}

type TreeNodeEntryParentGetFunc = func(TreeNodeEntry, Domain) []TreeNodeEntry
type TreeNodeEntryParentAddFunc = func(TreeNodeEntry, TreeNodeEntry, Domain) error
type TreeNodeEntryParentDelFunc = func(TreeNodeEntry, TreeNodeEntry, Domain) error
