package caskin

import (
	"gorm.io/gorm"
)

type BaseMetadataDB struct {
	DB *gorm.DB
}

func (b *BaseMetadataDB) AutoMigrate(dst ...interface{}) error {
	return b.DB.AutoMigrate(dst...)
}

func (b *BaseMetadataDB) Create(item interface{}) error {
	return b.DB.Create(item).Error
}

func (b *BaseMetadataDB) Recover(item interface{}) error {
	if err := b.TakeUnscoped(item); err != nil {
		return err
	}
	return b.DB.Model(item).Update("delete_at", nil).Error
}

func (b *BaseMetadataDB) Update(item interface{}) error {
	return b.DB.Updates(item).Error
}

func (b *BaseMetadataDB) UpsertType(item interface{}) UpsertType {
	if err := b.Take(item); err == nil {
		return UpsertTypeUpdate
	}
	if err := b.TakeUnscoped(item); err == nil {
		return UpsertTypeRecover
	}
	return UpsertTypeCreate
}

func (b *BaseMetadataDB) Take(item interface{}) error {
	return b.DB.Where(item).Take(item).Error
}

func (b *BaseMetadataDB) TakeUnscoped(item interface{}) error {
	return b.DB.Unscoped().Where(item).Take(item).Error
}

func (b *BaseMetadataDB) First(item interface{}, cond ...interface{}) error {
	return b.DB.First(item, cond...).Error
}

func (b *BaseMetadataDB) Find(items interface{}, cond ...interface{}) error {
	return b.DB.Find(items, cond...).Error
}

func (b *BaseMetadataDB) DeleteByID(item interface{}, id uint64) error {
	return b.DB.Delete(item, id).Error
}
