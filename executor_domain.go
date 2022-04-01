package caskin

// DomainCreate
// if there does not exist the domain then create a new one
// 1. no permission checking
// 2. create a new domain into metadata database
func (e *Executor) DomainCreate(domain Domain) error {
	if err := e.DBCreateCheck(domain); err != nil {
		return err
	}
	return e.DB.Create(domain)
}

// DomainRecover
// if there exist the domain but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one domain at metadata database
func (e *Executor) DomainRecover(domain Domain) error {
	if err := e.DBRecoverCheck(domain); err != nil {
		return err
	}
	return e.DB.Recover(domain)
}

// DomainDelete
// if there exist the domain soft delete the domain
// 1. no permission checking
// 2. delete all user's g in the domain
// 3. don't delete any role's g or object's g2 in the domain
// 4. soft delete one domain in metadata database
func (e *Executor) DomainDelete(domain Domain) error {
	if err := e.IDInterfaceDeleteCheck(domain); err != nil {
		return err
	}
	if err := e.Enforcer.RemoveUsersInDomain(domain); err != nil {
		return err
	}
	return e.DB.DeleteByID(domain, domain.GetID())
}

// DomainUpdate
// if there exist the domain update domain
// 1. no permission checking
// 2. just update domain's properties
func (e *Executor) DomainUpdate(domain Domain) error {
	if err := e.IDInterfaceUpdateCheck(domain); err != nil {
		return err
	}
	return e.DB.Update(domain)
}

// DomainInitialize
// if there exist the domain reinitialize the domain
// 1. no permission checking
// 2. just reinitialize the domain
func (e *Executor) DomainInitialize(domain Domain) error {
	if err := e.IDInterfaceUpdateCheck(domain); err != nil {
		return err
	}
	return e.initializeDomain(domain)
}

// DomainGet
// get all domain
// 1. no permission checking
func (e *Executor) DomainGet() ([]Domain, error) {
	return e.DB.GetAllDomain()
}

func (e *Executor) initializeDomain(domain Domain) error {
	// TODO
	// creator := e.options.DomainCreator(domain)
	// roles, objects := creator.BuildCreator()
	// for _, v := range objects {
	// 	if err := e.dbUpdateRoleOrObjectWhenInitializeDomain(v); err != nil {
	// 		return err
	// 	}
	// }
	//
	// creator.SetRelation()
	// for _, v := range roles {
	// 	if err := e.dbUpdateRoleOrObjectWhenInitializeDomain(v); err != nil {
	// 		return err
	// 	}
	// }
	//
	// policies := creator.GetPolicy()
	// for _, v := range policies {
	// 	if err := e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (e *Executor) dbUpdateRoleOrObjectWhenInitializeDomain(item roleOrObject) error {
	tmp := createByE(item)
	tmp.SetName(item.GetName())
	tmp.SetDomainID(item.GetDomainID())
	switch e.DB.UpsertType(tmp) {
	case UpsertTypeCreate:
		return e.DB.Create(item)
	case UpsertTypeRecover:
		if err := e.DB.Recover(tmp); err != nil {
			return err
		}
		item.SetID(tmp.GetID())
		return e.DB.Update(item)
	case UpsertTypeUpdate:
		item.SetID(tmp.GetID())
		return e.DB.Update(item)
	default:
		return nil
	}
}

type roleOrObject interface {
	nameInterface
	domainInterface
}
