package web_feature

import (
	"fmt"
	"github.com/awatercolorpen/caskin"
)

func initTreeMapFromPair(pair []*caskin.CustomizedDataPair) Relations {
	m := Relations{}
	for _, v := range pair {
		m[v.Object.GetID()] = Relation{}
	}
	for _, v := range pair {
		if v.Object.GetParentID() != 0 {
			m[v.Object.GetParentID()] = append(m[v.Object.GetParentID()], v.Object.GetID())
		}
	}
	return m
}

func initFeatureRelationMap(feature, frontend, backend Relations, relations Relations) Relations {
	m := Relations{}
	for k := range feature {
		m[k] = Relation{}
	}

	for k, relation := range relations {
		if _, ok := feature[k]; !ok {
			continue
		}
		for _, v := range relation {
			if _, ok := frontend[v]; ok {
				m[k] = append(m[k], v)
			}
			if _, ok := backend[v]; ok {
				m[k] = append(m[k], v)
			}
		}
	}
	return m
}

func initSingleFeatureRelation(feature, frontend, backend Relations, relation Relation) Relation {
	m := Relation{}
	for _, v := range relation {
		if _, ok := feature[v]; ok {
			continue
		}
		if _, ok := frontend[v]; ok {
			m = append(m, v)
		}
		if _, ok := backend[v]; ok {
			m = append(m, v)
		}
	}
	return m
}

func isEmptyObject(object caskin.Object) error {
	if object.GetID() != 0 {
		return caskin.ErrInValidObject
	}
	return nil
}

func isCompatible(m1, m2 Relations) bool {
	for k := range m1 {
		if _, ok := m2[k]; !ok {
			return false
		}
	}
	return true
}

func relationEncode(k, v interface{}) string {
	return fmt.Sprintf("%v%v%v", k, caskin.DefaultSeparator, v)
}

func relationDecode(in interface{}) (x, y uint64, err error) {
	format := fmt.Sprintf("%%d%v%%d", caskin.DefaultSeparator)
	_, err = fmt.Sscanf(in.(string), format, &x, &y)
	return
}

func relationsAction(in []interface{}, domain caskin.Domain, factory func() caskin.Object, action func(caskin.Object, caskin.Object, caskin.Domain) error) error {
	for _, v := range in {
		x, y, err := relationDecode(v)
		if err != nil {
			return err
		}
		ox, oy := factory(), factory()
		ox.SetID(x)
		oy.SetID(y)
		err = action(oy, ox, domain)
		if err != nil {
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
