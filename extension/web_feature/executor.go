package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e              *caskin.Executor
	root           *Root
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

func (e *Executor) frontendAndBackendUpdate(object caskin.Object) error {
	if err := e.e.ObjectTreeNodeUpdateCheck(object, e.objectFactory()); err != nil {
		return err
	}
	if err := e.e.ObjectTreeNodeParentCheck(object); err != nil {
		return err
	}

	provider := e.e.GetCurrentProvider()
	_, domain, _ := provider.Get()
	object.SetDomainID(domain.GetID())
	if err := e.e.DB.Update(object); err != nil {
		return err
	}

	newEntry := func() caskin.TreeNodeEntry { return e.objectFactory() }
	updater := caskin.NewTreeNodeEntryUpdater(newEntry, e.frontendAndBackendParentGetFunc(), e.e.DefaultObjectParentAddFunc(), e.e.DefaultObjectParentDelFunc())
	return updater.Run(object, domain)
}

func (e *Executor) frontendAndBackendParentGetFunc() caskin.TreeNodeEntryChildrenGetFunc {
	feature, _ := e.GetFeature()
	index := initTreeMapFromPair(feature)
	return func(p caskin.TreeNodeEntry, domain caskin.Domain) []caskin.TreeNodeEntry {
		var out []caskin.TreeNodeEntry
		parents := e.e.DefaultObjectParentGetFunc()(p, domain)
		for _, v := range parents {
			if _, ok := index[v.GetID()]; !ok {
				out = append(out, v)
			}
		}
		return out
	}
}

