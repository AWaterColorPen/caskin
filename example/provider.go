package example

import "github.com/awatercolorpen/caskin"

type Provider struct {
	User   caskin.User
	Domain caskin.Domain
}

func (p Provider) Get() (caskin.User, caskin.Domain, error) {
	return p.User, p.Domain, nil
}
