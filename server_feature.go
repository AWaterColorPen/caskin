package caskin

func (s *server) GetFeature() ([]*Feature, error) {
	return s.Dictionary.GetFeature()
}

func (s *server) GetBackend() ([]*Backend, error) {
	return s.Dictionary.GetBackend()
}

func (s *server) GetFrontend() ([]*Frontend, error) {
	return s.Dictionary.GetFrontend()
}

func (s *server) AuthBackend(user User, domain Domain, backend *Backend) error {
	value, err := s.Dictionary.GetBackendByKey(backend.Key())
	if err != nil {
		return ErrNoBackendPermission
	}
	if value == nil {
		value = &Backend{}
	}
	if s.CheckObject(user, domain, value.ToObject(), Read) != nil {
		return ErrNoBackendPermission
	}
	return nil
}

func (s *server) AuthFrontend(user User, domain Domain) []*Frontend {
	var out []*Frontend
	frontend, _ := s.Dictionary.GetFrontend()
	for _, v := range frontend {
		if s.CheckObject(user, domain, v.ToObject(), Read) == nil {
			out = append(out, v)
		}
	}
	return out
}

func (s *server) GetFeatureObject(user User, domain Domain) ([]Object, error) {
	var out []Object
	feature, _ := s.Dictionary.GetFeature()
	for _, v := range feature {
		if s.CheckObject(user, domain, v.ToObject(), Read) == nil {
			out = append(out, v.ToObject())
		}
	}
	return out, nil
}

func (s *server) GetFeaturePolicy(user User, domain Domain) ([]*Policy, error) {
	roles, err := s.GetRole(user, domain)
	if err != nil {
		return nil, err
	}
	objects, err := s.GetFeatureObject(user, domain)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	for _, v := range roles {
		policy := s.Enforcer.GetPoliciesForRoleInDomain(v, domain)
		for _, p := range policy {
			if object, ok := om[p.Object.GetID()]; ok {
				list = append(list, &Policy{
					Role:   v,
					Object: object,
					Domain: domain,
					Action: p.Action,
				})
			}
		}
	}
	return list, nil
}

func (s *server) GetFeaturePolicyByRole(user User, domain Domain, byRole Role) ([]*Policy, error) {
	if err := s.ObjectDataGetCheck(user, domain, byRole); err != nil {
		return nil, err
	}
	objects, err := s.GetFeatureObject(user, domain)
	if err != nil {
		return nil, err
	}

	om := IDMap(objects)
	var list []*Policy
	policy := s.Enforcer.GetPoliciesForRoleInDomain(byRole, domain)
	for _, p := range policy {
		if object, ok := om[p.Object.GetID()]; ok {
			list = append(list, &Policy{
				Role:   byRole,
				Object: object,
				Domain: domain,
				Action: p.Action,
			})
		}
	}
	return list, nil
}

func (s *server) ModifyFeaturePolicyPerRole(user User, domain Domain, perRole Role, input []*Policy) error {
	if err := s.ObjectDataModifyCheck(user, domain, perRole); err != nil {
		return err
	}
	if err := isValidPolicyWithRole(input, perRole); err != nil {
		return err
	}

	policy := s.Enforcer.GetPoliciesForRoleInDomain(perRole, domain)
	objects, err := s.GetFeatureObject(user, domain)
	if err != nil {
		return err
	}

	om := IDMap(objects)
	// make source and target role id list
	var source, target []*Policy
	for _, v := range policy {
		if _, ok := om[v.Object.GetID()]; ok {
			source = append(source, v)
		}
	}
	for _, v := range input {
		if _, ok := om[v.Object.GetID()]; ok {
			target = append(target, v)
		}
	}

	// get diff to add and remove
	add, remove := DiffPolicy(source, target)
	for _, v := range add {
		if err = s.Enforcer.AddPolicyInDomain(v.Role, v.Object, domain, Read); err != nil {
			return err
		}
	}
	for _, v := range remove {
		if err = s.Enforcer.RemovePolicyInDomain(v.Role, v.Object, domain, Read); err != nil {
			return err
		}
	}

	return nil
}
