package caskin

// CreateDomain
// if there does not exist the domain
// then create a new one without permission checking
// 1. create a new domain into metadata database
func (e *Executor) CreateDomain(domain Domain) error {
	if err := e.DBCreateCheck(domain); err != nil {
		return err
	}
	return e.DB.Create(domain)
}

// RecoverDomain
// if there exist the domain but soft deleted
// then recover it without permission checking
// 1. recover the soft delete one domain at metadata database
func (e *Executor) RecoverDomain(domain Domain) error {
	if err := e.DBRecoverCheck(domain); err != nil {
		return err
	}
	return e.DB.Recover(domain)
}

// DeleteDomain
// if there exist the domain
// soft delete the domain without permission checking
// 1. delete all user's g in the domain
// 2. don't delete any role's g or object's g2 in the domain
// 3. soft delete one domain in metadata database
func (e *Executor) DeleteDomain(domain Domain) error {
	if err := e.IDInterfaceDeleteCheck(domain); err != nil {
		return err
	}
	if err := e.Enforcer.RemoveUsersInDomain(domain); err != nil {
		return err
	}
	return e.DB.DeleteByID(domain, domain.GetID())
}

// UpdateDomain
// if there exist the domain
// Run domain without permission checking
// 1. just Run domain's properties
func (e *Executor) UpdateDomain(domain Domain) error {
	tmp := e.factory.NewDomain()
	if err := e.IDInterfaceUpdateCheck(domain, tmp); err != nil {
		return err
	}
	return e.DB.Update(domain)
}

// ReInitializeDomain
// if there exist the domain
// re initialize the domain without permission checking
// 1. just re initialize the domain
func (e *Executor) ReInitializeDomain(domain Domain) error {
	tmp := e.factory.NewDomain()
	if err := e.IDInterfaceUpdateCheck(domain, tmp); err != nil {
		return err
	}
	return e.initializeDomain(domain)
}

// GetAllDomain
// get all domain without permission checking
func (e *Executor) GetAllDomain() ([]Domain, error) {
	return e.DB.GetAllDomain()
}

// initializeDomain
// it is reentrant to initialize a new domain
// 1. get roles, objects, policies form DomainCreator
// 2. upsert roles, objects into metadata database
// 3. add policies as p into casbin
func (e *Executor) initializeDomain(domain Domain) error {
	creator := e.options.DomainCreator(domain)
	roles, objects := creator.BuildCreator()
	for _, v := range objects {
		if err := e.DB.Upsert(v); err != nil {
			return err
		}
	}
	for _, v := range roles {
		if err := e.DB.Upsert(v); err != nil {
			return err
		}
	}

	creator.SetRelation()
	for _, v := range roles {
		if err := e.DB.Upsert(v); err != nil {
			return err
		}
	}
	for _, v := range objects {
		if err := e.DB.Upsert(v); err != nil {
			return err
		}
	}

	policies := creator.GetPolicy()
	for _, v := range policies {
		if err := e.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}
