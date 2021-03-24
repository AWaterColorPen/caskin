package db

import (
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type Factory struct {
	agent   *Agent
	factory caskin.EntryFactory
}

func (f *Factory) NewCreator(domain caskin.Domain) *Creator {
	return &Creator{
		snapshot: f.agent.Snapshot,
		factory: f.factory,
		domain: domain,
	}
}

func NewFactory(db *gorm.DB, factory caskin.EntryFactory) (*Factory, error) {
	if err := db.AutoMigrate(&DomainCreatorObject{}, &DomainCreatorRole{}, &DomainCreatorPolicy{}); err != nil {
		return nil, err
	}
	agent := &Agent{db: db}
	return &Factory{agent: agent}, nil
}

type SnapshotFunc = func() ([]*DomainCreatorObject, []*DomainCreatorRole, []*DomainCreatorPolicy)
