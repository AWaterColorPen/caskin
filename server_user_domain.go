package caskin

import "github.com/ahmetb/go-linq/v3"

// GetUserByDomain
// get all user in domain
// 1. no permission checking
func (s *server) GetUserByDomain(domain Domain) ([]User, error) {
	us := s.Enforcer.GetUsersInDomain(domain)
	uid := ID(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	return s.DB.GetUserByID(uid)
}

// GetDomainByUser
// get user's all domain
// 1. no permission checking
func (s *server) GetDomainByUser(user User) ([]Domain, error) {
	if domain, err := s.getDomainBySuperadmin(user); err == nil {
		return domain, nil
	}
	ds := s.Enforcer.GetDomainsIncludeUser(user)
	did := ID(ds)
	linq.From(did).Distinct().ToSlice(&did)
	return s.DB.GetDomainByID(did)
}

func (s *server) getDomainBySuperadmin(user User) ([]Domain, error) {
	// TODO don't use SuperadminCheck
	if err := s.SuperadminCheck(user); err != nil {
		return nil, err
	}
	domain, err := s.GetDomain()
	if err != nil {
		return nil, err
	}
	domain = append(domain, GetSuperadminDomain())
	return domain, nil
}
