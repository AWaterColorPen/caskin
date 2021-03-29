package domain_creator_test

import (
	"github.com/awatercolorpen/caskin/extension/domain_creator"
	"path/filepath"
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func TestFactory(t *testing.T) {
	_, err := newFactory(t.TempDir())
	assert.NoError(t, err)
}

func newFactory(sqlitePath string) (*domain_creator.Factory, error) {
	gormDB, err := getDB(sqlitePath)
	if err != nil {
		return nil, err
	}
	factory, err := domain_creator.NewFactory(gormDB, &example.EntryFactory{})
	if err != nil {
		return nil, err
	}
	agent := factory.GetAgent()
	object := []*domain_creator.DomainCreatorObject{
		{Name: string(caskin.ObjectTypeObject), Type: caskin.ObjectTypeObject, RelativeObjectID: 1},
		{Name: string(caskin.ObjectTypeRole), Type: caskin.ObjectTypeRole, RelativeObjectID: 2},
		{Name: string(caskin.ObjectTypeDefault), Type: caskin.ObjectTypeDefault, RelativeObjectID: 3},
	}
	for _, v := range object {
		if err := agent.Create(v); err != nil {
			return nil, err
		}
	}
	role := []*domain_creator.DomainCreatorRole{
		{Name: "admin", RelativeObjectID: 2},
		{Name: "member", RelativeObjectID: 2},
	}
	for _, v := range role {
		if err := agent.Create(v); err != nil {
			return nil, err
		}
	}
	policy := []*domain_creator.DomainCreatorPolicy{
		{RelativeRoleID: 1, RelativeObjectID: 1, Action: caskin.Read},
		{RelativeRoleID: 1, RelativeObjectID: 1, Action: caskin.Write},
		{RelativeRoleID: 1, RelativeObjectID: 2, Action: caskin.Read},
		{RelativeRoleID: 1, RelativeObjectID: 2, Action: caskin.Write},
		{RelativeRoleID: 1, RelativeObjectID: 3, Action: caskin.Read},
		{RelativeRoleID: 1, RelativeObjectID: 3, Action: caskin.Write},
		{RelativeRoleID: 2, RelativeObjectID: 3, Action: caskin.Read},
		{RelativeRoleID: 2, RelativeObjectID: 3, Action: caskin.Write},
	}
	for _, v := range policy {
		if err := agent.Create(v); err != nil {
			return nil, err
		}
	}
	return factory, nil
}

func getDB(path string) (*gorm.DB, error) {
	dsn := filepath.Join(path, "sqlite")
	return gorm.Open(sqlite.Open(dsn), &gorm.Config{})
}
