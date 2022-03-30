package caskin

import "github.com/ahmetb/go-linq/v3"

// GetUserInDomain
// get all user in domain
func (e *Executor) GetUserInDomain(domain Domain) ([]User, error) {
	us := e.Enforcer.GetUsersInDomain(domain)
	uid := ID(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	return e.DB.GetUserByID(uid)
}

// GetDomainByUser
// get user's all domain
func (e *Executor) GetDomainByUser(user User) ([]Domain, error) {
	if domain, err := e.getDomainBySuperadmin(user); err == nil {
		return domain, nil
	}
	ds := e.Enforcer.GetDomainsIncludeUser(user)
	did := ID(ds)
	linq.From(did).Distinct().ToSlice(&did)
	return e.DB.GetDomainByID(did)
}

func (e *Executor) getDomainBySuperadmin(user User) ([]Domain, error) {
	if err := e.IsSuperadminCheck(user); err != nil {
		return nil, err
	}
	domain, err := e.GetAllDomain()
	if err != nil {
		return nil, err
	}
	domain = append(domain, e.options.GetSuperadminDomain())
	return domain, nil
}
