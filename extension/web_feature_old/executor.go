package web_feature_old

import (
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e               *caskin.Executor
	root            *Root
	objectFactory   caskin.ObjectFactory
	operationDomain caskin.Domain
	modelText       string
}

func (e *Executor) setBackendRoot(object caskin.Object) {
	e.root.SetBackendRoot(object)
}

func (e *Executor) setFeatureRoot(object caskin.Object) {
	e.root.SetFeatureRoot(object)
}

func (e *Executor) setFrontendRoot(object caskin.Object) {
	e.root.SetFrontendRoot(object)
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

func (e *Executor) allWebFeatureRelation(domain caskin.Domain) caskin.InheritanceRelations {
	queue := []caskin.Object{e.root.Feature, e.root.Frontend, e.root.Backend}
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
