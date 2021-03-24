package db

import (
	"github.com/awatercolorpen/caskin"
)

type Creator struct {
	snapshot     SnapshotFunc
	factory      caskin.EntryFactory
	snapshotObject []*DomainCreatorObject
	snapshotRole   []*DomainCreatorRole
	snapshotPolicy []*DomainCreatorPolicy

	domain  caskin.Domain
	objects caskin.Objects
	roles   caskin.Roles
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
	for i, v := range c.objects {
		if c.snapshotObject[i].RelativeObjectID != 0 {
			v.SetObjectID(c.objects[c.snapshotObject[i].RelativeObjectID].GetID())
		}
		if c.snapshotObject[i].RelativeParentID != 0 {
			v.SetParentID(c.objects[c.snapshotObject[i].RelativeParentID].GetID())
		}
	}
	for i, v := range c.roles {
		if c.snapshotRole[i].RelativeObjectID != 0 {
			v.SetObjectID(c.roles[c.snapshotRole[i].RelativeObjectID].GetID())
		}
		if c.snapshotRole[i].RelativeParentID != 0 {
			v.SetParentID(c.roles[c.snapshotRole[i].RelativeParentID].GetID())
		}
	}
}

func (c *Creator) GetPolicy() []*caskin.Policy {
	var policy []*caskin.Policy
	for _, v := range c.snapshotPolicy {
		r := c.factory.NewRole()
		o := c.factory.NewObject()
		if v.RelativeRoleID != 0 {
			r = c.roles[v.RelativeRoleID]
		} else {
			r.SetID(v.AbsoluteRoleID)
		}
		if v.RelativeObjectID != 0 {
			o = c.objects[v.RelativeObjectID]
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
