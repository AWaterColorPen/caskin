package db

import (
	"github.com/awatercolorpen/caskin"
)

type Factory struct {
	db MetaDB
}

func (f *Factory) NewCreator(domain caskin.Domain) *Creator {
	return &Creator{db: f.db, domain: domain}
}

func NewFactory(db MetaDB) (*Factory, error) {
	if err := db.AutoMigrate(&DomainCreatorObject{}, &DomainCreatorRole{}, &DomainCreatorPolicy{}); err != nil {
		return nil, err
	}
	return &Factory{db: db}, nil
}

// func GetAllDomainCreatorObject()