package example

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/manager"
)

// Stage example Stage for easy testing
type Stage struct {
	Manager        *manager.Manager // manager extension
	Caskin         *caskin.Caskin   // caskin instance on stage
	Options        *caskin.Options  // caskin options on stage
	Domain         *Domain          // a domain on stage
	SuperadminUser *User            // superadmin user on stage
	AdminUser      *User            // a domain admin user on stage
	MemberUser     *User            // a domain member user on stage
	SubAdminUser   *User            // a domain sub admin user on stage
}

func (s *Stage) AddSubAdmin() error {
	provider := caskin.NewCachedProvider(nil, nil)
	provider.Domain = s.Domain
	provider.User = s.AdminUser
	executor := s.Caskin.GetExecutor(provider)

	subAdmin := &User{
		PhoneNumber: "123456789031",
		Email:       "subadmin@qq.com",
	}
	if err := executor.CreateUser(subAdmin); err != nil {
		return err
	}

	objects, err := executor.GetObjects()
	if err != nil {
		return err
	}

	object1 := &Object{
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

	object2 := &Object{
		Name:     "role_sub_02",
		Type:     caskin.ObjectTypeRole,
		ObjectID: object1.ID,
		ParentID: objects[1].GetID(),
	}
	if err := executor.CreateObject(object2); err != nil {
		return err
	}

	role := &Role{
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
	if err := executor.DeleteSuperadminUser(s.SuperadminUser); err != nil {
		return err
	}

	s.SuperadminUser = nil
	return nil
}

func NewStageWithManger(m *manager.Manager) (*Stage, error) {
	c, err := m.GetCaskin()
	if err != nil {
		return nil, err
	}
	provider := caskin.NewCachedProvider(nil, nil)
	executor := c.GetExecutor(provider)

	domain := &Domain{Name: "domain_01"}
	if err := executor.CreateDomain(domain); err != nil {
		return nil, err
	}
	if err := executor.ReInitializeDomain(domain); err != nil {
		return nil, err
	}

	superadmin := &User{
		PhoneNumber: "12345678901",
		Email:       "superadmin@qq.com",
	}
	admin := &User{
		PhoneNumber: "12345678902",
		Email:       "admin@qq.com",
	}
	member := &User{
		PhoneNumber: "12345678903",
		Email:       "member@qq.com",
	}
	for _, v := range []caskin.User{superadmin, admin, member} {
		if err := executor.CreateUser(v); err != nil {
			return nil, err
		}
	}

	if err := executor.AddSuperadminUser(superadmin); err != nil {
		return nil, err
	}

	provider.Domain = domain
	provider.User = superadmin
	roles, err := executor.GetRoles()
	if err != nil {
		return nil, err
	}

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: admin, Role: roles[0]}},
		roles[1]: {{User: member, Role: roles[1]}},
	} {
		if err := executor.ModifyUserRolePairPerRole(k, v); err != nil {
			return nil, err
		}
	}

	stage := &Stage{
		Manager:        m,
		Caskin:         c,
		Options:        c.GetOptions(),
		Domain:         domain,
		SuperadminUser: superadmin,
		AdminUser:      admin,
		MemberUser:     member,
	}

	return stage, nil
}

func NewStageWithSqlitePath(sqlitePath string, options ...func(*manager.Configuration)) (*Stage, error) {
	m, err := NewManager(sqlitePath, options...)
	if err != nil {
		return nil, err
	}
	return NewStageWithManger(m)
}
