package example

import (
	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/extension/manager"
	"github.com/casbin/casbin/v2"
	gormadapter "github.com/casbin/gorm-adapter/v3"
)

func NewManager(configuration *manager.Configuration, sqlitePath string) (*manager.Manager, error) {
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
		Enforcer:      enforcer,
		EntryFactory:  &EntryFactory{},
		MetaDB:        NewGormMDBByDB(db),
		Extension: &manager.Extension{
			WebFeature:    0,
		},
	}

	return manager.NewManager(config)
}
