package db

import (
	"github.com/awatercolorpen/caskin"
)

type Creator struct {
	db MetaDB
	domain  caskin.Domain
	objects caskin.Objects
	roles   caskin.Roles
}

func (c *Creator) BuildCreator() ([]caskin.Role, []caskin.Object) {

	return c.roles, c.objects
}

func (c *Creator) SetRelation() {
}

func (c *Creator) GetPolicy() []*caskin.Policy {
	return []*caskin.Policy{
	}
}

func (c *Creator) GetRoles() []caskin.Role {
	return c.roles
}

func (c *Creator) GetObjects() []caskin.Object {
	return c.objects
}
