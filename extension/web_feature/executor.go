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

func (e *Executor) check(object caskin.Object) error {
	o := e.objectFactory()
	o.SetObjectID(object.GetID())
	return e.e.Enforce(o, caskin.Read)
}

func isEmptyObject(object caskin.Object) error {
	if object.GetID() != 0 {
		return caskin.ErrInValidObject
	}
	return nil
}
