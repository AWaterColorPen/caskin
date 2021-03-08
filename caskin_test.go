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
	_, err := newCaskin(t)
	assert.NoError(t, err)
}

func newCaskin(tb testing.TB) (*caskin.Caskin, error) {
	options := &caskin.Options{
		SuperadminOption: &caskin.SuperadminOption{
			Enable:             true,
		},
	}
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

func getTestDB(tb testing.TB) (*gorm.DB, error) {
	dsn := filepath.Join(tb.TempDir(), "sqlite")
	// dsn := filepath.Join("./", "sqlite")
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


func TestCaskin_GetExecutor(t *testing.T) {
	_, err := getStage(t)
	assert.NoError(t, err)
}


func getStage(t *testing.T) (*example.Stage, error) {
	c, err := newCaskin(t)
	if err != nil {
		return nil, err
	}

	provider := &example.Provider{}
	executor := c.GetExecutor(provider)

	domain := &example.Domain{Name: "domain_01"}
	if err := executor.CreateDomain(domain); err != nil {
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

	for _, v := range []*caskin.RolesForUser{
		{User:  admin, Roles: []caskin.Role{roles[0]}},
		{User:  member, Roles: []caskin.Role{roles[1]}},
	} {
		if err := executor.ModifyRolesForUser(v); err != nil {
			return nil, err
		}
	}

	stage := &example.Stage{
		Caskin: c,
		Domain: domain,
		SuperadminUser: superadmin,
		AdminUser: admin,
		MemberUser: member,
	}

	return stage, nil
}