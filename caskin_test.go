package caskin_test

import (
    "github.com/casbin/casbin/v2"
    "github.com/casbin/casbin/v2/model"
    "path/filepath"
    "testing"

    "github.com/awatercolorpen/caskin"
    "github.com/awatercolorpen/caskin/example"
    "github.com/casbin/gorm-adapter/v3"
    "github.com/stretchr/testify/assert"
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func TestNewCaskin(t *testing.T) {
    options := &caskin.Options{}
    _, err := newCaskin(t, options)
    assert.NoError(t, err)
}

func newCaskin(tb testing.TB, options *caskin.Options) (*caskin.Caskin, error) {
    db, err := getTestDB(tb)
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