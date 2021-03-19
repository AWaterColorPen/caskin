package example

import (
	"github.com/awatercolorpen/caskin"
)

type DomainCreator struct {
	domain  caskin.Domain
	objects []caskin.Object
	roles   []caskin.Role
}

func NewDomainCreator(domain caskin.Domain) caskin.Creator {
	return &DomainCreator{domain: domain}
}

func (d *DomainCreator) BuildCreator() ([]caskin.Role, []caskin.Object) {
	role0 := &Role{Name: "admin", DomainID: d.domain.GetID()}
	role1 := &Role{Name: "member", DomainID: d.domain.GetID()}
	d.roles = []caskin.Role{role0, role1}

	object0 := &Object{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, DomainID: d.domain.GetID()}
	object1 := &Object{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, DomainID: d.domain.GetID()}
	object2 := &Object{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, DomainID: d.domain.GetID()}
	d.objects = []caskin.Object{object0, object1, object2}

	return d.roles, d.objects
}

func (d *DomainCreator) SetRelation() {
	ooId := d.objects[0].GetID()
	for _, object := range d.objects {
		object.SetObjectID(ooId)
	}

	roId := d.objects[1].GetID()
	for _, role := range d.roles {
		role.SetObjectID(roId)
	}
}

func (d *DomainCreator) GetRoles() []caskin.Role {
	return d.roles
}

func (d *DomainCreator) GetObjects() []caskin.Object {
	return d.objects
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
