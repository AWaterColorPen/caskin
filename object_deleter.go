package caskin

type objectDeleter struct {
	visited  map[any]bool
	children ObjectChildrenGetFunc
	delete   ObjectDeleteFunc
}

func (t *objectDeleter) Run(current Object, domain Domain) error {
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

func NewObjectDeleter(children ObjectChildrenGetFunc, delete ObjectDeleteFunc) *objectDeleter {
	return &objectDeleter{
		visited:  map[any]bool{},
		children: children,
		delete:   delete,
	}
}

type ObjectChildrenGetFunc = func(Object, Domain) []Object
type ObjectDeleteFunc = func(Object, Domain) error
