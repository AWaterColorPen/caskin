package caskin

import "github.com/ahmetb/go-linq/v3"

// UserInDomainGet
// get all user in domain
// 1. no permission checking
func (e *Executor) UserInDomainGet(domain Domain) ([]User, error) {
	us := e.Enforcer.GetUsersInDomain(domain)
	uid := ID(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	return e.DB.GetUserByID(uid)
}

// DomainByUserGet
// get user's all domain
// 1. no permission checking
func (e *Executor) DomainByUserGet(user User) ([]Domain, error) {
	if domain, err := e.getDomainBySuperadmin(user); err == nil {
		return domain, nil
	}
	ds := e.Enforcer.GetDomainsIncludeUser(user)
	did := ID(ds)
	linq.From(did).Distinct().ToSlice(&did)
	return e.DB.GetDomainByID(did)
}

func (e *Executor) getDomainBySuperadmin(user User) ([]Domain, error) {
	if err := e.SuperadminCheck(user); err != nil {
		return nil, err
	}
	domain, err := e.DomainGet()
	if err != nil {
		return nil, err
	}
	domain = append(domain, GetSuperadminDomain())
	return domain, nil
}
