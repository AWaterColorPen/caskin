package caskin_test

import (
	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"path/filepath"
	"testing"
	"time"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestNewCaskin(t *testing.T) {
	_, err := newCaskin(t)
	assert.NoError(t, err)
}

func newCaskin(tb testing.TB) (*caskin.Caskin, error) {
	options := &caskin.Options{
		SuperAdminOption: &caskin.SuperAdminOption{
			Enable:             true,
			RealSuperadminInDB: true,
			Role:               nil,
			Domain:             nil,
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
	dsn := filepath.Join("./", "sqlite")
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

func getInitializeTestCaskin(t *testing.T) (*caskin.Caskin, error) {
	c, err := newCaskin(t)
	if err != nil {
		return nil, err
	}

	provider := &example.Provider{}
	executor := c.GetExecutor(provider)

	domain := &example.Domain{
		CreatedAt: time.Time{},
		UpdatedAt: time.Time{},
		DeletedAt: gorm.DeletedAt{},
		Name:      "test_domain_01",
	}
	// 创建domain
	if err := executor.CreateDomain(domain); err != nil {
		return nil, err
	}

	superAdmin := &example.User{
		CreatedAt:   time.Time{},
		UpdatedAt:   time.Time{},
		DeletedAt:   gorm.DeletedAt{},
		PhoneNumber: "12345678901",
		Email:       "superadmin@qq.com",
	}
	if err := executor.CreateUser(superAdmin); err != nil {
		return nil, err
	}

	if err := executor.AddSuperadminUser(superAdmin); err != nil {
		return nil, err
	}

	provider.Domain = domain
	provider.User = superAdmin

	executor = c.GetExecutor(provider)

	roles, err := executor.GetRoles()

	rolesForUser := &caskin.RolesForUser{
		User:  superAdmin,
		Roles: roles,
	}

	if err := executor.ModifyRolesForUser(rolesForUser); err != nil {
		return nil, err
	}

	// 创建superAdmin
	return c, nil
}

func TestCaskin_GetExecutor(t *testing.T) {
	_, err := getInitializeTestCaskin(t)
	assert.NoError(t, err)
}
