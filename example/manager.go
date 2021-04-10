package example

import (
	"path/filepath"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/manager"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func NewManager(sqlitePath string, options ...func(*manager.Configuration)) (*manager.Manager, error) {
	db, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&OneObjectData{})
	if err != nil {
		return nil, err
	}

	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := caskin.CasbinModel()
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	config := &manager.Configuration{
		DomainCreator: NewDomainCreator,
		Enforcer: enforcer,
		EntryFactory: &EntryFactory{},
		MetaDB: NewGormMDBByDB(db),
	}

	for _, v := range options {
		v(config)
	}

	return manager.NewManager(config)
}
