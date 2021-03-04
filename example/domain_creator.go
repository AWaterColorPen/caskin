package example

import (
	"github.com/awatercolorpen/caskin"
)

//func DomainCreator(domain caskin.Domain) ([]caskin.Role, []caskin.Object, []*caskin.Policy) {
//	object1 := &Object{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, DomainID: domain.GetID()}
//	object2 := &Object{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, DomainID: domain.GetID()}
//	object3 := &Object{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, DomainID: domain.GetID()}
//	objects := []caskin.Object{object1, object2, object3}
//
//	object1.ObjectID = object1.GetID()
//	object2.ObjectID = object1.GetID()
//	object2.ObjectID = object1.GetID()
//
//	if object1.Type == caskin.ObjectTypeRole {
//
//	}
//	role1 := &Role{Name: "admin", ObjectID: object2.GetID(), DomainID: domain.GetID()}
//	role2 := &Role{Name: "member", ObjectID: object2.GetID(), DomainID: domain.GetID()}
//	roles := []caskin.Role{role1, role2}
//
//	policies := []*caskin.Policy{
//		{role1, object1, domain, caskin.Read},
//		{role1, object1, domain, caskin.Write},
//		{role1, object2, domain, caskin.Read},
//		{role1, object2, domain, caskin.Write},
//		{role1, object3, domain, caskin.Read},
//		{role1, object3, domain, caskin.Write},
//		{role2, object3, domain, caskin.Read},
//		{role2, object3, domain, caskin.Write},
//	}
//
//	return roles, objects, policies
//}

type DomainCreator struct {
	domain  caskin.Domain
	objects caskin.Objects
	roles   caskin.Roles
}

func NewDomainCreator(domain caskin.Domain) *DomainCreator {
	return &DomainCreator{domain: domain}
}

func (d *DomainCreator) BuildCreator() (caskin.Roles, caskin.Objects) {

	role0 := &Role{Name: "admin", DomainID: d.domain.GetID()}
	role1 := &Role{Name: "member", DomainID: d.domain.GetID()}
	d.roles = []caskin.Role{role0, role1}

	object0 := &Object{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, DomainID: d.domain.GetID()}
	object1 := &Object{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, DomainID: d.domain.GetID()}
	object2 := &Object{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, DomainID: d.domain.GetID()}
	d.objects = []caskin.Object{object0, object1, object2}

	return d.roles, d.objects
}

func (d *DomainCreator) Set() {
	ooId := d.objects[0].GetID()
	for _, object := range d.objects {
		object.SetObjectId(ooId)
	}

	roId := d.objects[1].GetID()
	for _, role := range d.roles {
		role.SetObjectId(roId)
	}
}

func (d *DomainCreator) GetPolicy() []*caskin.Policy {
	return []*caskin.Policy{
		{d.roles[0], d.objects[0], d.domain, caskin.Read},
		{d.roles[0], d.objects[0], d.domain, caskin.Write},
		{d.roles[0], d.objects[1], d.domain, caskin.Read},
		{d.roles[0], d.objects[1], d.domain, caskin.Write},
		{d.roles[0], d.objects[2], d.domain, caskin.Read},
		{d.roles[0], d.objects[2], d.domain, caskin.Write},
		{d.roles[1], d.objects[2], d.domain, caskin.Read},
		{d.roles[1], d.objects[2], d.domain, caskin.Write},
	}
}

func (d *DomainCreator) GetRoles() caskin.Roles {
	return d.roles
}

func (d *DomainCreator) GetObjects() caskin.Objects {
	return d.objects
}
