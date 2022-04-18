package caskin

type objectUpdater struct {
	parentGet ObjectParentGetFunc
	parentAdd ObjectParentAddFunc
	parentDel ObjectParentDelFunc
}

func (t *objectUpdater) Run(item Object, domain Domain) error {
	var source, target []uint64
	if item.GetParentID() != 0 {
		target = append(target, item.GetParentID())
	}
	parents := t.parentGet(item, domain)
	for _, v := range parents {
		source = append(source, v.GetID())
	}

	add, remove := Diff(source, target)
	for _, v := range add {
		parent := newByE(item)
		parent.SetID(v)
		if err := t.parentAdd(item, parent, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		parent := newByE(item)
		parent.SetID(v)
		if err := t.parentDel(item, parent, domain); err != nil {
			return err
		}
	}
	return nil
}

func NewObjectUpdater(
	parentGet ObjectParentGetFunc,
	parentAdd ObjectParentAddFunc,
	parentDel ObjectParentDelFunc) *objectUpdater {
	return &objectUpdater{
		parentGet: parentGet,
		parentAdd: parentAdd,
		parentDel: parentDel,
	}
}

type ObjectParentGetFunc = func(Object, Domain) []Object
type ObjectParentAddFunc = func(Object, Object, Domain) error
type ObjectParentDelFunc = func(Object, Object, Domain) error
