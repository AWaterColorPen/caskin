package caskin

type executor struct {
	e        ienforcer
	mdb      MetaDB
	provider CurrentUserProvider
	factory  EntryFactory
	option   *Option
}

// GetAllDomain if user has domain's read permission
// 1. get all domain
func (e *executor) GetAllDomain() ([]Domain, error) {
	domains, err := e.mdb.GetAllDomain()
	if err != nil {
		return nil, err
	}

	out, err := e.filter(Read, domains)
	if err != nil {
		return nil, err
	}

	return out.([]Domain), nil
}

// DeleteDomain if user has domain's write permission
// 1. soft delete one domain at metadata database
// 2. delete all user's g in the domain
// 3. don't delete any role's g or object's g2 in the domain
func (e *executor) DeleteDomain(domain Domain) error {
	if err := isValid(domain); err != nil {
		return err
	}

	if err := e.mdb.TakeDomain(domain); err != nil {
		return err
	}

	if err := e.check(Write, domain); err != nil {
		return err
	}

	if err := e.e.RemoveUsersInDomain(domain); err != nil {
		return err
	}

	return e.mdb.DeleteDomainByID(domain.GetID())
}

func (e *executor) filter(action Action, source interface{}) (interface{}, error) {
	u, d, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	return Filter(e.e, u, d, action, e.factory.NewObject, source), nil
}

func (e *executor) check(action Action, one entry) error {
	u, d, err := e.provider.Get()
	if err != nil {
		return err
	}

	if ok := Check(e.e, u, d, action, e.factory.NewObject, one); !ok {
		switch action {
		case Read:
			return ErrNoReadPermission
		case Write:
			return ErrNoWritePermission
		default:
		}
	}

	return nil
}