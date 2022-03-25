package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type Root struct {
	Super caskin.Object `json:"super"`
}

func InitRootObject(domain caskin.Domain) (*Root, error) {
	r := &Root{}
	for _, v := range []caskin.Object{r.Super} {
		v.SetDomainID(domain.GetID())
		v.SetObjectID(r.Super.GetID())
	}
	return r, nil
}
