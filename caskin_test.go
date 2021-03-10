package caskin_test

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"path/filepath"
	"testing"
)

func TestNewCaskin(t *testing.T) {
	_, err := newCaskin(t, &caskin.Options{})
	assert.NoError(t, err)
}

func TestNewStage(t *testing.T) {
	stage, err := newStage(t)
	assert.NoError(t, err)
	assert.NoError(t, stageAddSubAdmin(stage))
}

func getTestDB(tb testing.TB) (*gorm.DB, error) {
	dsn := filepath.Join(tb.TempDir(), "sqlite")
	//dsn := filepath.Join("./", "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

var casbinModelMap = map[bool]model.Model{}

func getCasbinModel(options *caskin.Options) (model.Model, error) {
	k := options.IsEnableSuperAdmin()
	if _, ok := casbinModelMap[k]; !ok {
		m, err := caskin.CasbinModel(options)
		if err != nil {
			return nil, err
		}
		casbinModelMap[k] = m
	}

	return casbinModelMap[k], nil
}

func newCaskin(tb testing.TB, options *caskin.Options) (*caskin.Caskin, error) {
	db, err := getTestDB(tb)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(
		&example.User{},
		&example.Domain{},
		&example.Role{},
		&example.Object{})
	if err != nil {
		return nil, err
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := getCasbinModel(options)
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return caskin.New(options,
		caskin.DomainCreatorOption(example.NewDomainCreator),
		caskin.EnforcerOption(enforcer),
		caskin.EntryFactoryOption(&example.EntryFactory{}),
		caskin.MetaDBOption(example.NewGormMDBByDB(db)),
	)
}

func newStage(t *testing.T) (*example.Stage, error) {
	options := &caskin.Options{
		SuperadminOption: &caskin.SuperadminOption{
			Enable: true,
		},
	}
	c, err := newCaskin(t, options)
	if err != nil {
		return nil, err
	}

	provider := &example.Provider{}
	executor := c.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_01"}
	if err := executor.CreateDomain(domain); err != nil {
		return nil, err
	}
	if err := executor.ReInitializeDomain(domain); err != nil {
		return nil, err
	}

	superadmin := &example.User{
		PhoneNumber: "12345678901",
		Email:       "superadmin@qq.com",
	}
	admin := &example.User{
		PhoneNumber: "12345678902",
		Email:       "admin@qq.com",
	}
	member := &example.User{
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

	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: admin, Role: roles[0]}},
		roles[1]: {{User: member, Role: roles[1]}},
	} {
		if err := executor.ModifyUserRolePairPerRole(k, v); err != nil {
			return nil, err
		}
	}

	stage := &example.Stage{
		Caskin:         c,
		Options:        options,
		Domain:         domain,
		SuperadminUser: superadmin,
		AdminUser:      admin,
		MemberUser:     member,
	}

	return stage, nil
}

func stageAddSubAdmin(stage *example.Stage) error {
	provider := &example.Provider{
		Domain: stage.Domain,
		User:   stage.AdminUser,
	}
	executor := stage.Caskin.GetExecutor(provider)

	subAdmin := &example.User{
		PhoneNumber: "123456789031",
		Email:       "subadmin@qq.com",
	}
	if err := executor.CreateUser(subAdmin); err != nil {
		return err
	}

	object1 := &example.Object{
		Name: "object_sub_01",
		Type: caskin.ObjectTypeObject,
		ObjectID: 1,
		ParentID: 1,
	}
	if err := executor.CreateObject(object1); err != nil {
		return err
	}
	object1.ObjectID = object1.ID
	if err := executor.UpdateObject(object1); err != nil {
		return err
	}

	object2 := &example.Object{
		Name: "role_sub_02",
		Type: caskin.ObjectTypeRole,
		ObjectID: object1.ID,
		ParentID: 2,
	}
	if err := executor.CreateObject(object2); err != nil {
		return err
	}

	role := &example.Role{
		Name: "admin_sub_01",
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
		{role, object1, stage.Domain, caskin.Read},
		{role, object1, stage.Domain, caskin.Write},
		{role, object2, stage.Domain, caskin.Read},
		{role, object2, stage.Domain, caskin.Write},
	}
	if err := executor.ModifyPolicyListPerRole(role, policy); err != nil {
		return err
	}

	stage.SubAdminUser = subAdmin

	return nil
}

