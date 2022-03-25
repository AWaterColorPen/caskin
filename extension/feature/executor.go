package feature

import (
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e          *caskin.Executor
	Dictionary Dictionary
}

func (e *Executor) AuthBackend(backend *Backend) error {
	value, err := e.Dictionary.GetBackendByKey(backend.GetKey())
	if err != nil || backend == nil {
		return caskin.ErrNoBackendPermission
	}
	if e.e.EnforceObject(value.ToObject(), caskin.Read) != nil {
		return caskin.ErrNoBackendPermission
	}
	return nil
}

func (e *Executor) AuthFrontend() []*Frontend {
	var res []*Frontend
	frontend, _ := e.Dictionary.GetFrontend()
	for _, v := range frontend {
		if e.e.EnforceObject(v.ToObject(), caskin.Read) == nil {
			res = append(res, v)
		}
	}
	return res
}

func (e *Executor) GetFeature() ([]*Feature, error) {
	if err := e.e.IsSuperadminAndSuperdomainCheck(); err != nil {
		return nil, err
	}
	return e.Dictionary.GetFeature()
}

func (e *Executor) GetBackend() ([]*Backend, error) {
	if err := e.e.IsSuperadminAndSuperdomainCheck(); err != nil {
		return nil, err
	}
	return e.Dictionary.GetBackend()
}

func (e *Executor) GetFrontend() ([]*Frontend, error) {
	if err := e.e.IsSuperadminAndSuperdomainCheck(); err != nil {
		return nil, err
	}
	return e.Dictionary.GetFrontend()
}

func (e *Executor) GetSourceRelation(domain caskin.Domain) map[caskin.Object][]caskin.Object {
	var queue []caskin.Object
	inQueue := map[uint64]bool{}
	for _, v := range queue {
		inQueue[v.GetID()] = true
	}

	m := map[caskin.Object][]caskin.Object{}
	for i := 0; i < len(queue); i++ {
		m[queue[i]] = []caskin.Object{}
		ll := e.e.Enforcer.GetChildrenForObjectInDomain(queue[i], domain)
		for _, v := range ll {
			if _, ok := inQueue[v.GetID()]; !ok {
				queue = append(queue, v)
				inQueue[v.GetID()] = true
			}
			m[queue[i]] = append(m[queue[i]], v)
		}
	}

	return m
}

func (e *Executor) GetTargetRelation() map[caskin.Object][]caskin.Object {
	m := map[caskin.Object][]caskin.Object{}
	return m
}
