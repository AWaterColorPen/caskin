package example

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type gormMDB struct {
	db *gorm.DB
}

func (g *gormMDB) Create(item interface{}) error {
	return g.db.Create(item).Error
}

func (g *gormMDB) Recover(item interface{}) error {
	if err := g.TakeUnscoped(item); err != nil {
		return err
	}
	return g.db.Model(item).Update("delete_at", nil).Error
}

func (g *gormMDB) Update(item interface{}) error {
	return g.db.Updates(item).Error
}

func (g *gormMDB) Upsert(item interface{}) error {
	v, ok := item.(entry)
	if !ok || v.GetID() == 0 {
		if err := g.Recover(item); err == nil {
			return nil
		}
		return g.Create(item)
	}
	return g.Update(item)
}

func (g *gormMDB) Take(item interface{}) error {
	return g.db.Where(item).Take(item).Error
}

func (g *gormMDB) TakeUnscoped(item interface{}) error {
	return g.db.Unscoped().Where(item).Take(item).Error
}

func (g *gormMDB) GetUserByID(id []uint64) ([]caskin.User, error) {
	var user []*User
	if err := g.db.Find(&user, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.User
	linq.From(user).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) DeleteUserByID(id uint64) error {
	return g.db.Delete(&User{ID: id}).Error
}

func (g *gormMDB) GetRoleInDomain(domain caskin.Domain) ([]caskin.Role, error) {
	var role []*Role
	if err := g.db.Where(&Role{DomainID: domain.GetID()}).Find(&role).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Role
	linq.From(role).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) GetRoleByID(id []uint64) ([]caskin.Role, error) {
	var role []*Role
	if err := g.db.Find(&role, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Role
	linq.From(role).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) DeleteRoleByID(id uint64) error {
	return g.db.Delete(&Role{}, id).Error
}

func (g *gormMDB) GetObjectInDomain(domain caskin.Domain, objectType ...caskin.ObjectType) ([]caskin.Object, error) {
	d := g.db.Where("domain_id = ?", domain.GetID())
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
	if err := g.db.Find(&object, "id IN ?", id).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Object
	linq.From(object).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) DeleteObjectByID(id uint64) error {
	return g.db.Delete(&Object{}, id).Error
}

func (g *gormMDB) GetAllDomain() ([]caskin.Domain, error) {
	var domain []*Domain
	if err := g.db.Find(&domain).Error; err != nil {
		return nil, err
	}

	var ret []caskin.Domain
	linq.From(domain).ToSlice(&ret)
	return ret, nil
}

func (g *gormMDB) DeleteDomainByID(id uint64) error {
	return g.db.Delete(&Domain{}, id).Error
}

func NewGormMDBByDB(db *gorm.DB) caskin.MetaDB {
	return &gormMDB{db: db}
}

type entry interface {
	GetID() uint64
}
