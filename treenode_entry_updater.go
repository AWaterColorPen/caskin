package caskin

type treeNodeEntryUpdater struct {
	newEntry  func() TreeNodeEntry
	parentGet TreeNodeEntryParentGetFunc
	parentAdd TreeNodeEntryParentAddFunc
	parentDel TreeNodeEntryParentDelFunc
}

func (t *treeNodeEntryUpdater) Run(item TreeNodeEntry, domain Domain) error {
	var source, target []interface{}
	if item.GetParentID() != 0 {
		target = append(target, item.GetParentID())
	}
	parents := t.parentGet(item, domain)
	for _, v := range parents {
		source = append(source, v.GetID())
	}

	add, remove := Diff(source, target)
	for _, v := range add {
		parent := t.newEntry()
		parent.SetID(v.(uint64))
		if err := t.parentAdd(item, parent, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		parent := t.newEntry()
		parent.SetID(v.(uint64))
		if err := t.parentDel(item, parent, domain); err != nil {
			return err
		}
	}
	return nil
}

func NewTreeNodeEntryUpdater(newEntry func() TreeNodeEntry,
	parentGet TreeNodeEntryParentGetFunc,
	parentAdd TreeNodeEntryParentAddFunc,
	parentDel TreeNodeEntryParentDelFunc) *treeNodeEntryUpdater {
	return &treeNodeEntryUpdater{
		newEntry:  newEntry,
		parentGet: parentGet,
		parentAdd: parentAdd,
		parentDel: parentDel,
	}
}

type TreeNodeEntryParentGetFunc = func(TreeNodeEntry, Domain) []TreeNodeEntry
type TreeNodeEntryParentAddFunc = func(TreeNodeEntry, TreeNodeEntry, Domain) error
type TreeNodeEntryParentDelFunc = func(TreeNodeEntry, TreeNodeEntry, Domain) error
