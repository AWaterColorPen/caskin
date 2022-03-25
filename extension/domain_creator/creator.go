package domain_creator

import (
	"github.com/awatercolorpen/caskin"
)

type Creator struct {
	snapshot       SnapshotFunc
	factory        caskin.EntryFactory
	snapshotObject []*DomainCreatorObject
	snapshotRole   []*DomainCreatorRole
	snapshotPolicy []*DomainCreatorPolicy

	domain  caskin.Domain
	objects []caskin.Object
	roles   []caskin.Role

	oIndex map[uint64]uint64
	rIndex map[uint64]uint64
}

func (c *Creator) BuildCreator() ([]caskin.Role, []caskin.Object) {
	c.snapshotObject, c.snapshotRole, c.snapshotPolicy = c.snapshot()
	for _, v := range c.snapshotObject {
		o := c.factory.NewObject()
		o.SetName(v.Name)
		o.SetObjectType(v.Type)
		o.SetObjectID(v.AbsoluteObjectID)
		o.SetDomainID(c.domain.GetID())
		o.SetParentID(v.AbsoluteParentID)
		c.objects = append(c.objects, o)
	}
	for _, v := range c.snapshotRole {
		r := c.factory.NewRole()
		r.SetName(v.Name)
		r.SetObjectID(v.AbsoluteObjectID)
		r.SetDomainID(c.domain.GetID())
		r.SetParentID(v.AbsoluteParentID)
		c.roles = append(c.roles, r)
	}
	return c.roles, c.objects
}

func (c *Creator) SetRelation() {
	c.buildIndex()
	for i, v := range c.objects {
		o := c.snapshotObject[i]
		if o.RelativeObjectID != 0 {
			v.SetObjectID(c.oIndex[o.RelativeObjectID])
		}
		if o.RelativeParentID != 0 {
			v.SetParentID(c.oIndex[o.RelativeParentID])
		}
	}
	for i, v := range c.roles {
		r := c.snapshotRole[i]
		if r.RelativeObjectID != 0 {
			v.SetObjectID(c.rIndex[r.RelativeObjectID])
		}
		if r.RelativeParentID != 0 {
			v.SetParentID(c.rIndex[r.RelativeParentID])
		}
	}
}

func (c *Creator) GetPolicy() []*caskin.Policy {
	var policy []*caskin.Policy
	for _, v := range c.snapshotPolicy {
		r := c.factory.NewRole()
		o := c.factory.NewObject()

		if v.RelativeRoleID != 0 {
			r.SetID(c.rIndex[v.RelativeRoleID])
		} else {
			r.SetID(v.AbsoluteRoleID)
		}
		if v.RelativeObjectID != 0 {
			o.SetID(c.oIndex[v.RelativeObjectID])
		} else {
			o.SetID(v.AbsoluteObjectID)
		}

		policy = append(policy, &caskin.Policy{
			Role: r, Object: o, Domain: c.domain, Action: v.Action,
		})
	}
	return policy
}

func (c *Creator) GetRoles() []caskin.Role {
	return c.roles
}

func (c *Creator) GetObjects() []caskin.Object {
	return c.objects
}

func (c *Creator) buildIndex() {
	c.oIndex = map[uint64]uint64{}
	// for i, v := range c.snapshotObject {
	// c.oIndex[v.ID] = c.objects[i].GetID()
	// }
	c.rIndex = map[uint64]uint64{}
	for i, v := range c.snapshotRole {
		c.rIndex[v.ID] = c.roles[i].GetID()
	}
}
