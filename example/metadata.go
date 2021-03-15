package example

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type gormMDB struct {
	caskin.BaseMetadataDB
}

func (g *gormMDB) GetUserByID(id []uint64) ([]caskin.User, error) {
	var user []*User
	if err := g.DB.Find(&user, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.User
	linq.From(user).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetRoleInDomain(domain caskin.Domain) ([]caskin.Role, error) {
	var role []*Role
	if err := g.DB.Where(&Role{DomainID: domain.GetID()}).Find(&role).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Role
	linq.From(role).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetRoleByID(id []uint64) ([]caskin.Role, error) {
	var role []*Role
	if err := g.DB.Find(&role, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Role
	linq.From(role).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetObjectInDomain(domain caskin.Domain, objectType ...caskin.ObjectType) ([]caskin.Object, error) {
	d := g.DB.Where("domain_id = ?", domain.GetID())
	if len(objectType) > 0 {
		d = d.Where("type IN ?", objectType)
	}

	var object []*Object
	if err := d.Find(&object).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Object
	linq.From(object).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetObjectByID(id []uint64) ([]caskin.Object, error) {
	var object []*Object
	if err := g.DB.Find(&object, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Object
	linq.From(object).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetAllDomain() ([]caskin.Domain, error) {
	var domain []*Domain
	if err := g.DB.Find(&domain).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Domain
	linq.From(domain).ToSlice(&ret)
	return ret, nil
}

func NewGormMDBByDB(db *gorm.DB) caskin.MetaDB {
	return &gormMDB{
		BaseMetadataDB: caskin.BaseMetadataDB{DB: db},
	}
}
