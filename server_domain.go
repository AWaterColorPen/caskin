package caskin

// CreateDomain
// if there does not exist the domain then create a new one
// 1. no permission checking
// 2. create a new domain into metadata database
func (s *server) CreateDomain(domain Domain) error {
	if err := s.DBCreateCheck(domain); err != nil {
		return err
	}
	return s.DB.Create(domain)
}

// RecoverDomain
// if there exist the domain but soft deleted then recover it
// 1. no permission checking
// 2. recover the soft delete one domain at metadata database
func (s *server) RecoverDomain(domain Domain) error {
	if err := s.DBRecoverCheck(domain); err != nil {
		return err
	}
	return s.DB.Recover(domain)
}

// DeleteDomain
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

// UpdateDomain
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

// GetDomain
// get all domain
// 1. no permission checking
func (s *server) GetDomain() ([]Domain, error) {
	return s.DB.GetAllDomain()
}

// ResetDomain
// if there exist the domain reset the domain
// 1. no permission checking
// 2. just reset the domain
func (s *server) ResetDomain(domain Domain) error {
	old := newByE(domain)
	if err := s.IDInterfaceUpdateCheck(domain, old); err != nil {
		return err
	}
	return s.resetDomain(domain)
}

func (s *server) resetDomain(domain Domain) error {
	co, _ := s.Dictionary.GetCreatorObject()
	cr, _ := s.Dictionary.GetCreatorRole()
	cp, _ := s.Dictionary.GetCreatorPolicy()

	roleObjectID := uint64(0)
	var object []Object
	for _, v := range co {
		o := v.ToObject()
		o.SetDomainID(domain.GetID())
		if err := s.dbRoleOrObjectWhenResetDomain(o); err != nil {
			return err
		}
		if o.GetObjectType() == ObjectTypeRole {
			roleObjectID = o.GetID()
		}
		object = append(object, o)
	}
	if roleObjectID == 0 {
		return ErrInValidObject
	}

	var role []Role
	for _, v := range cr {
		r := v.ToRole()
		r.SetDomainID(domain.GetID())
		r.SetObjectID(roleObjectID)
		if err := s.dbRoleOrObjectWhenResetDomain(r); err != nil {
			return err
		}
		role = append(role, r)
	}

	for _, v := range cp {
		var o Object
		for i, u := range co {
			if u.Name == v.Object {
				o = object[i]
			}
		}
		var r Role
		for i, u := range cr {
			if u.Name == v.Role {
				r = role[i]
			}
		}
		for _, action := range v.Action {
			if err := s.Enforcer.AddPolicyInDomain(r, o, domain, Action(action)); err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *server) dbRoleOrObjectWhenResetDomain(item any) error {
	// tmp := newByE(item)
	// tmp.SetName(item.GetName())
	// tmp.SetDomainID(item.GetDomainID())
	switch s.DB.UpsertType(item) {
	case UpsertTypeCreate:
		return s.DB.Create(item)
	case UpsertTypeRecover:
		if err := s.DB.Recover(item); err != nil {
			return err
		}
		return s.DB.Update(item)
	case UpsertTypeUpdate:
		return s.DB.Update(item)
	default:
		return nil
	}
}

// type roleOrObject interface {
// 	idInterface
// 	nameInterface
// 	domainInterface
// }
