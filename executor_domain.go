package caskin

// CreateDomain if there does not exist the domain, then create a new one
// 1. create a new domain into metadata database
// 2. initialize the new domain
func (e *executor) CreateDomain(domain Domain) error {
	return e.createOrRecoverDomain(domain, e.mdb.CreateDomain)
}

// RecoverDomain if there exist the domain but soft deleted, then recover it
// 1. soft delete one domain at metadata database
// 2. re initialize the recovering domain
func (e *executor) RecoverDomain(domain Domain) error {
	return e.createOrRecoverDomain(domain, e.mdb.RecoverDomain)
}

// DeleteDomain if user has domain's write permission
// 1. delete all user's g in the domain
// 2. don't delete any role's g or object's g2 in the domain
// 3. soft delete one domain in metadata database
func (e *executor) DeleteDomain(domain Domain) error {
	fn := func(domain Domain) error {
		if err := e.e.RemoveUsersInDomain(domain); err != nil {
			return err
		}
		return e.mdb.DeleteDomainByID(domain.GetID())
	}

	return e.writeDomain(domain, fn)
}

// UpdateDomain if there exist the domain and user has domain's write permission
// 1. just update domain's properties
func (e *executor) UpdateDomain(domain Domain) error {
	return e.writeDomain(domain, e.mdb.UpdateDomain)
}

// ReInitializeDomain if there exist the domain and user has domain's write permission
// 1. just re initialize the domain
func (e *executor) ReInitializeDomain(domain Domain) error {
	return e.writeDomain(domain, e.initializeDomain)
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

func (e *executor) createOrRecoverDomain(domain Domain, fn func(Domain) error) error {
	if err := e.mdb.TakeDomain(domain); err == nil {
		return ErrAlreadyExists
	}

	if err := fn(domain); err != nil {
		return err
	}

	return e.initializeDomain(domain)
}

func (e *executor) writeDomain(domain Domain, fn func(Domain) error) error {
	if err := isValid(domain); err != nil {
		return err
	}

	if err := e.mdb.TakeDomain(domain); err != nil {
		return ErrNotExists
	}

	if err := e.check(Write, domain); err != nil {
		return err
	}

	return fn(domain)
}

// initializeDomain it is reentrant to initialize a new domain
// 1. get roles, objects, policies form DomainCreator
// 2. upsert roles, objects into metadata database
// 3. add policies as p into casbin
func (e *executor) initializeDomain(domain Domain) error {
	roles, objects, policies := e.option.DomainCreator(domain)
	for _, v := range roles {
		if err := e.mdb.UpsertRole(v); err != nil {
			return err
		}
	}

	for _, v := range objects {
		if err := e.mdb.UpsertObject(v); err != nil {
			return err
		}
	}

	for _, v := range policies {
		if err := e.e.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
			return err
		}
	}

	return nil
}