package caskin

type cachedProvider struct {
	User   User
	Domain Domain
	err    error
}

func (c *cachedProvider) Get() (User, Domain, error) {
	if c.err != nil {
		return nil, nil, c.err
	}
	if c.User == nil || c.Domain == nil {
		return nil, nil, ErrProviderGet
	}
	return c.User, c.Domain, nil
}

func NewCachedProvider(user User, domain Domain, err error) *cachedProvider {
	return &cachedProvider{
		User:   user,
		Domain: domain,
		err:    err,
	}
}
