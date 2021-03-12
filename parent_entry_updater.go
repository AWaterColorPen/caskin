package caskin

type parentEntryUpdater struct {
	parentGetFn parentGetFn
	parentAddFn parentAddFn
	parentDelFn parentDelFn
	newEntry    func() treeNodeEntry
}

func (p *parentEntryUpdater) update(item treeNodeEntry, domain Domain) error {
	var source, target []interface{}
	if item.GetParentID() != 0 {
		target = append(target, item.GetParentID())
	}
	parents := p.parentGetFn(item, domain)
	for _, v := range parents {
		source = append(source, v.GetID())
	}

	add, remove := Diff(source, target)
	for _, v := range add {
		parent := p.newEntry()
		parent.SetID(v.(uint64))
		if err := p.parentAddFn(item, parent, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		parent := p.newEntry()
		parent.SetID(v.(uint64))
		if err := p.parentDelFn(item, parent, domain); err != nil {
			return err
		}
	}
	return nil
}

type parentGetFn = func(treeNodeEntry, Domain) []treeNodeEntry
type parentAddFn = func(treeNodeEntry, treeNodeEntry, Domain) error
type parentDelFn = func(treeNodeEntry, treeNodeEntry, Domain) error
