package caskin

// DomainCreate
// if there does not exist the domain then create a new one
// 1. no permission checking
// 2. create a new domain into metadata database
func (s *server) CreateDomain(domain Domain) error {
	if err := s.DBCreateCheck(domain); err != nil {
		return err
	}
	return s.DB.Create(domain)
}

// DomainRecover
// if there exist the domain but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one domain at metadata database
func (s *server) RecoverDomain(domain Domain) error {
	if err := s.DBRecoverCheck(domain); err != nil {
		return err
	}
	return s.DB.Recover(domain)
}

// DomainDelete
// if there exist the domain soft delete the domain
// 1. no permission checking
// 2. delete all user's g in the domain
// 3. don't delete any role's g or object's g2 in the domain
// 4. soft delete one domain in metadata database
func (s *server) DeleteDomain(domain Domain) error {
	if err := s.IDInterfaceDeleteCheck(domain); err != nil {
		return err
	}
	if err := s.Enforcer.RemoveUsersInDomain(domain); err != nil {
		return err
	}
	return s.DB.DeleteByID(domain, domain.GetID())
}

// DomainUpdate
// if there exist the domain update domain
// 1. no permission checking
// 2. just update domain's properties
func (s *server) UpdateDomain(domain Domain) error {
	old := newByE(domain)
	if err := s.IDInterfaceUpdateCheck(domain, old); err != nil {
		return err
	}
	return s.DB.Update(domain)
}

// DomainReset
// if there exist the domain reset the domain
// 1. no permission checking
// 2. just reinitialize the domain
func (s *server) ResetDomain(domain Domain) error {
	old := newByE(domain)
	if err := s.IDInterfaceUpdateCheck(domain, old); err != nil {
		return err
	}
	return s.initializeDomain(domain)
}

// DomainGet
// get all domain
// 1. no permission checking
func (s *server) GetDomain() ([]Domain, error) {
	return s.DB.GetAllDomain()
}

func (s *server) initializeDomain(domain Domain) error {
	// TODO
	// creator := s.options.DomainCreator(domain)
	// roles, objects := creator.BuildCreator()
	// for _, v := range objects {
	// 	if err := s.dbRoleUpdateOrObjectWhenInitializeDomain(v); err != nil {
	// 		return err
	// 	}
	// }
	//
	// creator.SetRelation()
	// for _, v := range roles {
	// 	if err := s.dbRoleUpdateOrObjectWhenInitializeDomain(v); err != nil {
	// 		return err
	// 	}
	// }
	//
	// policies := creator.GetPolicy()
	// for _, v := range policies {
	// 	if err := s.Enforcer.AddPolicyInDomain(v.Role, v.Object, v.Domain, v.Action); err != nil {
	// 		return err
	// 	}
	// }

	return nil
}

func (s *server) dbRoleUpdateOrObjectWhenInitializeDomain(item roleOrObject) error {
	tmp := newByE(item)
	tmp.SetName(item.GetName())
	tmp.SetDomainID(item.GetDomainID())
	switch s.DB.UpsertType(tmp) {
	case UpsertTypeCreate:
		return s.DB.Create(item)
	case UpsertTypeRecover:
		if err := s.DB.Recover(tmp); err != nil {
			return err
		}
		item.SetID(tmp.GetID())
		return s.DB.Update(item)
	case UpsertTypeUpdate:
		item.SetID(tmp.GetID())
		return s.DB.Update(item)
	default:
		return nil
	}
}

type roleOrObject interface {
	idInterface
	nameInterface
	domainInterface
}
