package example

import (
    "github.com/awatercolorpen/caskin"
)

func DomainCreator(domain caskin.Domain) ([]caskin.Role, []caskin.Object, []*caskin.Policy) {
    object1 := &Object{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, DomainID: domain.GetID()}
    object2 := &Object{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, DomainID: domain.GetID()}
    object3 := &Object{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, DomainID: domain.GetID()}
    objects := []caskin.Object{object1, object2, object3}

    role1 := &Role{Name: "admin", ObjectID: object2.GetID(), DomainID: domain.GetID()}
    role2 := &Role{Name: "member", ObjectID: object2.GetID(), DomainID: domain.GetID()}
    roles := []caskin.Role{role1, role2}

    policies := []*caskin.Policy{
        {role1, object1, domain, caskin.Read},
        {role1, object1, domain, caskin.Write},
        {role1, object2, domain, caskin.Read},
        {role1, object2, domain, caskin.Write},
        {role1, object3, domain, caskin.Read},
        {role1, object3, domain, caskin.Write},
        {role2, object3, domain, caskin.Read},
        {role2, object3, domain, caskin.Write},
    }

    return roles, objects, policies
}
