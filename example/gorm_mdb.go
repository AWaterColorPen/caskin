package example

import (
	"errors"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type gormMDB struct {
	db *gorm.DB
}

func (g *gormMDB) CreateUser(user caskin.User) error {
	return g.db.Create(user).Error
}

func (g *gormMDB) RecoverUser(user caskin.User) error {
	return g.recover(user)
}

func (g *gormMDB) UpdateUser(user caskin.User) error {
	return g.db.Updates(user).Error
}

func (g *gormMDB) TakeUser(user caskin.User) error {
	return g.db.Where(user).Take(user).Error
}

func (g *gormMDB) GetUserInDomain(domain caskin.Domain) ([]caskin.User, error) {
	panic("implement me")
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

func (g *gormMDB) UpsertUser(user caskin.User) error {
	return g.upsert(user)
}

func (g *gormMDB) DeleteUserByID(id uint64) error {
	return g.db.Delete(&User{ID: id}).Error
}

func (g *gormMDB) CreateRole(role caskin.Role) error {
	return g.db.Create(role).Error
}

func (g *gormMDB) RecoverRole(role caskin.Role) error {
	return g.recover(role)
}

func (g *gormMDB) UpdateRole(role caskin.Role) error {
	return g.db.Updates(role).Error
}

func (g *gormMDB) TakeRole(role caskin.Role) error {
	return g.db.Where(role).Take(role).Error
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

func (g *gormMDB) UpsertRole(role caskin.Role) error {
	return g.upsert(role)
}

func (g *gormMDB) DeleteRoleByID(id uint64) error {
	return g.db.Delete(&Role{}, id).Error
}

func (g *gormMDB) CreateObject(object caskin.Object) error {
	return g.db.Create(object).Error
}

func (g *gormMDB) RecoverObject(object caskin.Object) error {
	return g.recover(object)
}

func (g *gormMDB) UpdateObject(object caskin.Object) error {
	return g.db.Updates(object).Error
}

func (g *gormMDB) TakeObject(object caskin.Object) error {
	return g.db.Where(object).Take(object).Error
}

func (g *gormMDB) GetObjectInDomain(domain caskin.Domain, objectType ...caskin.ObjectType) ([]caskin.Object, error) {
	o := &Object{DomainID: domain.GetID()}
	if len(objectType) > 0 {
		o.Type = objectType[0]
	}

	var object []*Object
	if err := g.db.Find(&object, "type IN ?", objectType).Error; err != nil {
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

func (g *gormMDB) UpsertObject(object caskin.Object) error {
	return g.upsert(object)
}

func (g *gormMDB) DeleteObjectByID(id uint64) error {
	return g.db.Delete(&Object{}, id).Error
}

func (g *gormMDB) TakeDeletedObject(object caskin.Object) error {
	return g.db.Unscoped().Where(object).Take(object).Error
}

func (g *gormMDB) CreateDomain(domain caskin.Domain) error {
	return g.db.Create(domain).Error
}

func (g *gormMDB) RecoverDomain(domain caskin.Domain) error {
	return g.recover(domain)
}

func (g *gormMDB) UpdateDomain(domain caskin.Domain) error {
	return g.db.Updates(domain).Error
}

func (g *gormMDB) TakeDomain(domain caskin.Domain) error {
	return g.db.Where(domain).Take(domain).Error
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

func (g *gormMDB) upsert(entry entry) error {
	if entry.GetID() == 0 {
		return g.insertOrRecover(entry)
	}
	return g.db.Updates(entry).Error
}

func (g *gormMDB) recover(item interface{}) error {
	if err := g.db.Unscoped().Where(item).Take(item).Error; err != nil {
		return err
	}
	return g.db.Model(item).Update("delete_at", nil).Error
}

func (g *gormMDB) insertOrRecover(item interface{}) error {
	return g.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Unscoped().Where(item).Take(item).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return tx.Create(item).Error
			}
			return err
		}
		return tx.Model(item).Update("delete_at", nil).Error
	})
}

func NewGormMDBByDB(db *gorm.DB) caskin.MetaDB {
	return &gormMDB{db: db}
}

type entry interface {
	GetID() uint64
}
