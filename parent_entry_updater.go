package caskin

type parentEntryUpdater struct {
	parentsFn      parentsFn
	addParentFn    addParentFn
	deleteParentFn deleteParentFn
	newEntry       func() parentEntry
}

func (p *parentEntryUpdater) update(item parentEntry, domain Domain) error {
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

type parentsFn = func(parentEntry, Domain) []parentEntry
type addParentFn = func(parentEntry, parentEntry, Domain) error
type deleteParentFn = func(parentEntry, parentEntry, Domain) error

func singleParentsFunc(item parentEntry, newEntry func() parentEntry) parentsFn {
	return func(parentEntry, Domain) []parentEntry {
		parent := newEntry()
		parent.SetID(item.GetParentID())
		return []parentEntry{parent}
	}
}
