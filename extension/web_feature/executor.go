package web_feature

import "github.com/awatercolorpen/caskin"

type Executor struct {
	e             *caskin.Executor
	objectFactory caskin.ObjectFactory
}

func (e *Executor) operationPermissionCheck() error {
	// TODO get user and domain, check super domain and superadmin
	return nil
}