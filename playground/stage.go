package playground

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
)

// Stage example Stage for easy testing
type Stage struct {
	Caskin         *caskin.Caskin  // caskin instance on stage
	Options        *caskin.Options // caskin options on stage
	Domain         *example.Domain // a domain on stage
	SuperadminUser *example.User   // superadmin user on stage
	AdminUser      *example.User   // a domain admin user on stage
	MemberUser     *example.User   // a domain member user on stage
	SubAdminUser   *example.User   // a domain sub admin user on stage
}

func (s *Stage) AddSubAdmin() error {
	provider := caskin.NewCachedProvider(nil, nil)
	provider.Domain = s.Domain
	provider.User = s.AdminUser
	executor := s.Caskin.GetExecutor(provider)

	subAdmin := &example.User{
		PhoneNumber: "123456789031",
		Email:       "subadmin@qq.com",
	}
	if err := executor.UserCreate(subAdmin); err != nil {
		return err
	}

	objects, err := executor.GetObjects()
	if err != nil {
		return err
	}

	object1 := &example.Object{
		Name:     "object_sub_01",
		Type:     caskin.ObjectTypeObject,
		ObjectID: objects[0].GetID(),
		ParentID: objects[0].GetID(),
	}
	if err := executor.CreateObject(object1); err != nil {
		return err
	}
	object1.ObjectID = object1.ID
	if err := executor.UpdateObject(object1); err != nil {
		return err
	}

	object2 := &example.Object{
		Name:     "role_sub_02",
		Type:     caskin.ObjectTypeRole,
		ObjectID: object1.ID,
		ParentID: objects[1].GetID(),
	}
	if err := executor.CreateObject(object2); err != nil {
		return err
	}

	role := &example.Role{
		Name:     "admin_sub_01",
		ObjectID: object2.ID,
		ParentID: 1,
	}
	if err := executor.CreateRole(role); err != nil {
		return err
	}

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		role: {{User: subAdmin, Role: role}},
	} {
		if err := executor.ModifyUserRolePairPerRole(k, v); err != nil {
			return err
		}
	}

	policy := []*caskin.Policy{
		{role, object1, s.Domain, caskin.Read},
		{role, object1, s.Domain, caskin.Write},
		{role, object2, s.Domain, caskin.Read},
		{role, object2, s.Domain, caskin.Write},
	}
	if err := executor.ModifyPolicyListPerRole(role, policy); err != nil {
		return err
	}

	s.SubAdminUser = subAdmin
	return nil
}

func (s *Stage) NoSuperadmin() error {
	provider := caskin.NewCachedProvider(nil, nil)
	executor := s.Caskin.GetExecutor(provider)
	if err := executor.SuperadminDelete(s.SuperadminUser); err != nil {
		return err
	}

	s.SuperadminUser = nil
	return nil
}

func NewStageWithSqlitePath(sqlitePath string, options ...func(*manager.Configuration)) (*Stage, error) {
	m, err := NewManager(sqlitePath, options...)
	if err != nil {
		return nil, err
	}
	return NewStageWithManger(m)
}
