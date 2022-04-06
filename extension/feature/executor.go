package feature

import (
	"github.com/awatercolorpen/caskin"
)

type Executor struct {
	e          *caskin.baseService
	Dictionary Dictionary
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
