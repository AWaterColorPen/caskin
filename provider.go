package caskin

// CurrentProvider include current user and domain
type CurrentProvider interface {
	Get() (User, Domain, error)
}

type cachedProvider struct {
	User   User
	Domain Domain
}

func (c *cachedProvider) Get() (User, Domain, error) {
	if c.User == nil || c.Domain == nil {
		return nil, nil, ErrProviderGet
	}
	return c.User, c.Domain, nil
}

func NewCachedProvider(user User, domain Domain) *cachedProvider {
	return &cachedProvider{
		User:   user,
		Domain: domain,
	}
}
