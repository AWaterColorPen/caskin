package example

import "github.com/awatercolorpen/caskin"

type Provider struct {
	User   caskin.User
	Domain caskin.Domain
}

func (p *Provider) Get() (caskin.User, caskin.Domain, error) {
	if p.User == nil || p.Domain == nil {
		return nil, nil, caskin.ErrProviderGet
	}
	return p.User, p.Domain, nil
}
