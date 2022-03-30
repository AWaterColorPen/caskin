package web_feature

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
)

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

func isValidAction(action caskin.Action) error {
	if action != caskin.Read {
		return caskin.ErrInValidAction
	}
	return nil
}

func isPolicyValidAction(policy *caskin.Policy) error {
	return isValidAction(policy.Action)
}

func isPolicyListValidAction(list []*caskin.Policy) error {
	for _, v := range list {
		if err := isPolicyValidAction(v); err != nil {
			return err
		}
	}
	return nil
}
