package feature

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
)

func (e *Executor) Build() error {
	domains, err := e.e.DomainGet()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err := e.buildToOneDomain(v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) buildToOneDomain(domain caskin.Domain) error {
	sourceRelations := e.GetSourceRelation(domain)
	targetRelations := e.GetTargetRelation()
	var source, target []any
	for k, relation := range sourceRelations {
		for _, v := range relation {
			source = append(source, relationEncode(k, v))
		}
	}
	for k, relation := range targetRelations {
		for _, v := range relation {
			target = append(target, relationEncode(k, v))
		}
	}

	var toDelete []caskin.Object
	for k := range sourceRelations {
		if _, ok := targetRelations[k]; !ok {
			toDelete = append(toDelete, k)
		}
	}
	for _, v := range toDelete {
		if err := e.e.Enforcer.RemoveObjectInDomain(v, domain); err != nil {
			return err
		}
	}

	graph := MergedInheritanceRelations(sourceRelations, targetRelations)
	index := caskin.TopSort(graph)
	sorter := caskin.NewEdgeSorter(index)

	add, remove := caskin.Diff(source, target)
	if err := relationsAction(add, sorter.RootFirstSort, domain, e.objectFactory, e.e.Enforcer.AddParentForObjectInDomain); err != nil {
		return err
	}

	if err := relationsAction(remove, sorter.LeafFirstSort, domain, e.objectFactory, e.e.Enforcer.RemoveParentForObjectInDomain); err != nil {
		return err
	}

	return nil
}

func relationEncode(k, v any) string {
	return fmt.Sprintf("%v%v%v", k, caskin.DefaultSeparator, v)
}

func relationDecode(in any) (edge *Edge, err error) {
	format := fmt.Sprintf("%%d%v%%d", caskin.DefaultSeparator)
	edge = &Edge{}
	_, err = fmt.Sscanf(in.(string), format, &edge.X, &edge.Y)
	return
}

func relationsAction(in []any,
	sortFn func([]*Edge),
	domain caskin.Domain,
	factory caskin.ObjectFactory,
	action func(caskin.Object, caskin.Object, caskin.Domain) error) error {
	var edges []*Edge
	for _, v := range in {
		edge, err := relationDecode(v)
		if err != nil {
			return err
		}
		edges = append(edges, edge)
	}

	sortFn(edges)

	for _, v := range edges {
		ox, oy := factory(), factory()
		ox.SetID(v.X)
		oy.SetID(v.Y)
		if err := action(oy, ox, domain); err != nil {
			return err
		}
	}
	return nil
}

type InheritanceRelations = map[caskin.Object][]caskin.Object

func MergedInheritanceRelations(relations ...InheritanceRelations) InheritanceRelations {
	return nil
}
