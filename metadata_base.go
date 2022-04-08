package caskin

import (
	"gorm.io/gorm"
)

type BaseMetadataDB struct {
	DB *gorm.DB
}

func (b *BaseMetadataDB) AutoMigrate(dst ...any) error {
	return b.DB.AutoMigrate(dst...)
}

func (b *BaseMetadataDB) Create(item any) error {
	return b.DB.Create(item).Error
}

func (b *BaseMetadataDB) Recover(item any) error {
	if err := b.TakeUnscoped(item); err != nil {
		return err
	}
	return b.DB.Unscoped().Model(item).Update("delete_at", nil).Error
}

func (b *BaseMetadataDB) Update(item any) error {
	return b.DB.Updates(item).Error
}

func (b *BaseMetadataDB) UpsertType(item any) UpsertType {
	if err := b.Take(item); err == nil {
		return UpsertTypeUpdate
	}
	if err := b.TakeUnscoped(item); err == nil {
		return UpsertTypeRecover
	}
	return UpsertTypeCreate
}

func (b *BaseMetadataDB) Take(item any) error {
	return b.DB.Where(item).Take(item).Error
}

func (b *BaseMetadataDB) TakeUnscoped(item any) error {
	return b.DB.Unscoped().Where(item).Take(item).Error
}

func (b *BaseMetadataDB) Find(items any, cond ...any) error {
	return b.DB.Find(items, cond...).Error
}

func (b *BaseMetadataDB) DeleteByID(item any, id uint64) error {
	return b.DB.Delete(item, id).Error
}
