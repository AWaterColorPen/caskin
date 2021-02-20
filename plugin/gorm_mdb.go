package plugin

import (
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type gormMDB struct {
	db *gorm.DB
}

func (g *gormMDB) GetUserByID(uint64s []uint64) ([]caskin.User, error) {
	panic("implement me")
}

func (g *gormMDB) TakeRole(role caskin.Role) error {
	panic("implement me")
}

func (g *gormMDB) GetRoleInDomain(domain caskin.Domain) ([]caskin.Role, error) {
	panic("implement me")
}

func (g *gormMDB) GetRoleByID(uint64s []uint64) ([]caskin.Role, error) {
	panic("implement me")
}

func (g *gormMDB) UpsertRole(role caskin.Role) error {
	panic("implement me")
}

func (g *gormMDB) DeleteRoleByID(u uint64) error {
	panic("implement me")
}

func (g *gormMDB) TakeObject(role caskin.Role) error {
	panic("implement me")
}

func (g *gormMDB) GetObjectInDomain(domain caskin.Domain, objectType ...caskin.ObjectType) ([]caskin.Object, error) {
	panic("implement me")
}

func (g *gormMDB) GetObjectByID(uint64s []uint64) ([]caskin.Object, error) {
	panic("implement me")
}

func (g *gormMDB) UpsertObject(object caskin.Object) error {
	panic("implement me")
}

func (g *gormMDB) DeleteObjectByID(u uint64) error {
	panic("implement me")
}

func (g *gormMDB) GetAllDomain() ([]caskin.Domain, error) {
	panic("implement me")
}

func (g *gormMDB) DeleteDomainByID(u uint64) error {
	panic("implement me")
}

func NewGormMDBByDB(db *gorm.DB) caskin.MetaDB {
	return &gormMDB{
		db: db,
	}
}
