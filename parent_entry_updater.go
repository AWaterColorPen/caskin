package caskin

type parentEntryUpdater struct {
	parentsFn      parentsFn
	addParentFn    addParentFn
	deleteParentFn deleteParentFn
	newEntry       func() treeNodeEntry
}

func (p *parentEntryUpdater) update(item treeNodeEntry, domain Domain) error {
	var source, target []interface{}
	if item.GetParentID() != 0 {
		target = append(target, item.GetParentID())
	}
	parents := p.parentsFn(item, domain)
	for _, v := range parents {
		source = append(source, v.GetID())
	}

	add, remove := Diff(source, target)
	for _, v := range add {
		parent := p.newEntry()
		parent.SetID(v.(uint64))
		if err := p.addParentFn(item, parent, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		parent := p.newEntry()
		parent.SetID(v.(uint64))
		if err := p.deleteParentFn(item, parent, domain); err != nil {
			return err
		}
	}
	return nil
}

type parentsFn = func(treeNodeEntry, Domain) []treeNodeEntry
type addParentFn = func(treeNodeEntry, treeNodeEntry, Domain) error
type deleteParentFn = func(treeNodeEntry, treeNodeEntry, Domain) error

func singleParentsFunc(item treeNodeEntry, newEntry func() treeNodeEntry) parentsFn {
	return func(treeNodeEntry, Domain) []treeNodeEntry {
		parent := newEntry()
		parent.SetID(item.GetParentID())
		return []treeNodeEntry{parent}
	}
}
