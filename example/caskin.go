package example

import (
	"path/filepath"

	"github.com/awatercolorpen/caskin"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}

func NewCaskin(options *caskin.Options, sqlitePath string) (*caskin.Caskin, error) {
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

	m, err := caskin.CasbinModel(options.IsDisableSuperadmin())
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}

	return caskin.New(options,
		caskin.DomainCreatorOption(NewDomainCreator),
		caskin.EnforcerOption(enforcer),
		caskin.EntryFactoryOption(&EntryFactory{}),
		caskin.MetaDBOption(NewGormMDBByDB(db)),
	)
}
