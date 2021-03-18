package caskin

import "github.com/ahmetb/go-linq/v3"

// GetUsers
// get all user in current domain
func (e *Executor) GetUserInDomain() ([]User, error) {
	_, currentDomain, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	us := e.Enforcer.GetUsersInDomain(currentDomain)
	uid := getIDList(us)
	linq.From(uid).Distinct().ToSlice(&uid)
	return e.DB.GetUserByID(uid)
}

// GetDomainByUser
// get user's all domain
func (e *Executor) GetDomainByUser() ([]Domain, error) {
	if e.IsSuperadminCheck() == nil {
		domain, err := e.GetAllDomain()
		if err != nil {
			return nil, err
		}

		domain = append(domain, e.options.GetSuperadminDomain())
		return domain, nil
	}

	currentUser, _, err := e.provider.Get()
	if err != nil {
		return nil, err
	}

	ds := e.Enforcer.GetDomainsIncludeUser(currentUser)
	did := getIDList(ds)
	linq.From(did).Distinct().ToSlice(&did)
	return e.DB.GetDomainByID(did)
}
