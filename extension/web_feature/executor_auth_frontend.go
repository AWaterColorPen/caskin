package web_feature

import "github.com/awatercolorpen/caskin"

type CasbinStruct struct {
	M  string     `json:"m"`
	P  [][]string `json:"p"`
	G  [][]string `json:"g"`
	G2 [][]string `json:"g2"`
}

func (e *Executor) AuthFrontend() []*Frontend {
	var res []*Frontend
	frontend, _ := e.getFrontend()
	for _, v := range frontend {
		if e.check(v.Object) == nil {
			res = append(res, v.ObjectCustomizedData.(*Frontend))
		}
	}
	return res
}

func (e *Executor) AuthFrontendCaskinStruct(subject string) (*CasbinStruct, error) {
	casbin := &CasbinStruct{}
	provider := e.e.GetCurrentProvider()
	_, domain, err := provider.Get()
	if err != nil {
		return nil, err
	}

	if e.e.IsSuperadminCheck() == nil && domain.Encode() == e.operationDomain.Encode() {
		casbin.G = append(casbin.G, []string{"g", subject, caskin.SuperadminRole, e.operationDomain.Encode()})
		return casbin, nil
	}

	frontend := e.AuthFrontend()
	for _, v := range frontend {
		casbin.P = append(casbin.P, []string{"p", subject, domain.Encode(), v.GetName(), string(caskin.Read)})
	}
	return casbin, nil
}
