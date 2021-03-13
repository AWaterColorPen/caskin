package web_feature

import (
	"github.com/awatercolorpen/caskin"
	"gorm.io/gorm"
)

type MetaDataDBAPI struct {
	db *gorm.DB
}

func (m *MetaDataDBAPI) Create(item caskin.ObjectData, bind caskin.Object) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(bind).Error; err != nil {
			return err
		}
		item.SetObjectID(bind.GetID())
		return tx.Create(item).Error
	})
}

func (m *MetaDataDBAPI) Recover(item caskin.ObjectData, bind caskin.Object) error {
	// TODO
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(bind).Error; err != nil {
			return err
		}
		item.SetObjectID(bind.GetID())
		return tx.Create(item).Error
	})
}

func (m *MetaDataDBAPI) Update(item caskin.ObjectData, bind caskin.Object) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(bind).Error; err != nil {
			return err
		}
		return tx.Updates(item).Error
	})
}

func (m *MetaDataDBAPI) DeleteByID(item caskin.ObjectData, bind caskin.Object) error {
	return m.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(bind, bind.GetID()).Error; err != nil {
			return err
		}
		return tx.Delete(item, item.GetID()).Error
	})
}
