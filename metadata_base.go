package caskin

import "gorm.io/gorm"

type BaseMetadataDB struct {
	DB *gorm.DB
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

func (b *BaseMetadataDB) Upsert(item interface{}) error {
	if v, ok := item.(idInterface); ok && v.GetID() != 0 {
		return b.Update(item)
	}
	if err := b.Recover(item); err == nil {
		return nil
	}
	return b.Create(item)
}

func (b *BaseMetadataDB) Take(item interface{}) error {
	return b.DB.Where(item).Take(item).Error
}

func (b *BaseMetadataDB) TakeUnscoped(item interface{}) error {
	return b.DB.Unscoped().Where(item).Take(item).Error
}

func (b *BaseMetadataDB) DeleteByID(item interface{}, id uint64) error {
	return b.DB.Delete(item, id).Error
}
