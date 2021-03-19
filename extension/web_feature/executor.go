package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e                         *caskin.Executor
	objectFactory             caskin.ObjectFactory
	operationDomain           caskin.Domain
	enableBackendAPIAuthCache bool
	FeatureRootObject         func() caskin.Object
}

func (e *Executor) get3pair() (feature, frontend, backend []*caskin.CustomizedDataPair, err error) {
	feature, err = e.GetFeature()
	if err != nil {
		return
	}
	frontend, err = e.GetFrontend()
	if err != nil {
		return
	}
	backend, err = e.GetBackend()
	return
}

func (e *Executor) allWebFeatureRelation(domain caskin.Domain) caskin.InheritanceRelations {
	queue := []caskin.Object{GetFeatureRootObject(), GetBackendRootObject(), GetFrontendRootObject()}
	inQueue := map[uint64]bool{}
	for _, v := range queue {
		inQueue[v.GetID()] = true
	}

	m := caskin.InheritanceRelations{}
	for i := 0; i < len(queue); i++ {
		m[queue[i].GetID()] = caskin.InheritanceRelation{}
		ll := e.e.Enforcer.GetChildrenForObjectInDomain(queue[i], domain)
		for _, v := range ll {
			if _, ok := inQueue[v.GetID()]; !ok {
				queue = append(queue, v)
				inQueue[v.GetID()] = true
			}
			m[queue[i].GetID()] = append(m[queue[i].GetID()], v.GetID())
		}
	}

	return m
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

func (e *Executor) filterWithNoError(source interface{}) []interface{} {
	var result []interface{}
	linq.From(source).Where(func(v interface{}) bool {
		return e.check(v.(caskin.Object)) == nil
	}).ToSlice(&result)
	return result
}

func (e *Executor) check(object caskin.Object) error {
	o := e.objectFactory()
	o.SetObjectID(object.GetID())
	return e.e.Enforce(o, caskin.Read)
}
