package caskin_test

import (
    "github.com/casbin/casbin/v2"
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

    _, err := caskin.New(options, )
    assert.Error(t, err)
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

    model, err := caskin.CasbinModel(options)
    if err != nil {
        return nil, err
    }

    enforcer, err := casbin.NewSyncedEnforcer(model, adapter)
    if err != nil {
        return nil, err
    }

    return caskin.New(options,
        caskin.DomainCreatorOption(nil),
        caskin.EnforcerOption(enforcer),
        caskin.EntryFactoryOption(&example.EntryFactory{}),
        caskin.MetaDBOption(example.NewGormMDBByDB(db)),
        )
}

func getTestDB(tb testing.TB) (*gorm.DB, error) {
    dsn := filepath.Join(tb.TempDir(), "sqlite")
    return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
