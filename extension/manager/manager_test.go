package manager_test

import (
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/extension/manager"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func enforcer(db *gorm.DB) (casbin.IEnforcer, error) {
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := caskin.CasbinModel()
	if err != nil {
		return nil, err
	}

	return casbin.NewSyncedEnforcer(m, adapter)
}

func TestManager(t *testing.T) {
	config := &manager.Configuration{}
	_, err := manager.NewManager(config)
	assert.Equal(t, caskin.ErrInitializationNilDomainCreator, err)

	config.DomainCreator = example.NewDomainCreator
	_, err = manager.NewManager(config)
	assert.Equal(t, caskin.ErrInitializationNilEnforcer, err)

	db, err := getDB(t.TempDir())
	assert.NoError(t, err)

	config.Enforcer, err = enforcer(db)
	_, err = manager.NewManager(config)
	assert.Equal(t, caskin.ErrInitializationNilEntryFactory, err)

	config.EntryFactory = &example.EntryFactory{}
	_, err = manager.NewManager(config)
	assert.Equal(t, caskin.ErrInitializationNilMetaDB, err)

	config.MetaDB = example.NewGormMDBByDB(db)
	m1, err := manager.NewManager(config)
	assert.NoError(t, err)

	_, err = m1.GetCaskin()
	assert.NoError(t, err)

	_, err = m1.GetWebFeature()
	assert.Equal(t, manager.ErrNoInitialization, err)

	_, err = m1.GetDomainCreatorFactory()
	assert.Equal(t, manager.ErrNoInitialization, err)
}

func TestExtension(t *testing.T) {
	db, err := getDB(t.TempDir())
	assert.NoError(t, err)

	e, err := enforcer(db)
	assert.NoError(t, err)

	config := &manager.Configuration{
		DomainCreator: example.NewDomainCreator,
		Enforcer:      e,
		EntryFactory:  &example.EntryFactory{},
		MetaDB:        example.NewGormMDBByDB(db),
		Extension: &manager.Extension{
			DomainCreator: 0,
			WebFeature:    0,
		},
	}

	_, err = manager.NewManager(config)
	assert.Equal(t, manager.ErrInitializationNilDB, err)

	config.DB = db
	_, err = manager.NewManager(config)
	assert.Equal(t, manager.ErrExtensionConfigurationConflict, err)

	config.DomainCreator = nil
	m1, err := manager.NewManager(config)
	assert.NoError(t, err)

	_, err = m1.GetCaskin()
	assert.NoError(t, err)

	_, err = m1.GetWebFeature()
	assert.NoError(t, err)

	_, err = m1.GetDomainCreatorFactory()
	assert.NoError(t, err)
}
