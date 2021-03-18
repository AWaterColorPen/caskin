package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e                         *caskin.Executor
	objectFactory             caskin.ObjectFactory
	operationDomain           caskin.Domain
	enableBackendAPIAuthCache bool
	FeatureRootObject         func() caskin.Object
}

func (e *Executor) operationPermissionCheck() error {
	provider := e.e.GetCurrentProvider()
	_, domain, err := provider.Get()
	if err != nil {
		return err
	}
	if domain.Encode() != e.operationDomain.Encode() {
		return caskin.ErrCanOnlyAllowAtValidDomain
	}
	return nil
}
