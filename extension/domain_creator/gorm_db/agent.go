package gorm_db

import "gorm.io/gorm"

type Agent struct {
	db *gorm.DB
}

func (a *Agent) Snapshot() ([]*DomainCreatorObject, []*DomainCreatorRole, []*DomainCreatorPolicy) {
	return dbSnapshot(a.db)
}

func (a *Agent) GetDomainCreatorObject() ([]*DomainCreatorObject, error) {
	var out []*DomainCreatorObject
	return out, a.db.Find(&out).Error
}

func (a *Agent) GetDomainCreatorRole() ([]*DomainCreatorRole, error) {
	var out []*DomainCreatorRole
	return out, a.db.Find(&out).Error
}

func (a *Agent) GetDomainCreatorPolicy() ([]*DomainCreatorPolicy, error) {
	var out []*DomainCreatorPolicy
	return out, a.db.Find(&out).Error
}

func (a *Agent) Create(item relativeIDAndAbsoluteID) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(item).Error; err != nil {
			return err
		}
		return isValidSnapshot(tx)
	})
}

// func (a *Agent) Recover(item relativeIDAndAbsoluteID) error {
//     if err := item.IsValid(); err != nil {
//         return err
//     }
//     return a.gorm_db.Transaction(func(tx *gorm.DB) error {
//         if err := tx.Create(item).Error; err != nil {
//             return err
//         }
//         return isValidSnapshot(tx)
//     })
// }

func (a *Agent) Update(item relativeIDAndAbsoluteID) error {
	if err := item.IsValid(); err != nil {
		return err
	}
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Updates(item).Error; err != nil {
			return err
		}
		return isValidSnapshot(tx)
	})
}

func (a *Agent) Delete(item interface{}, id uint64) error {
	return a.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(item, id).Error; err != nil {
			return err
		}
		return isValidSnapshot(tx)
	})
}

func dbSnapshot(db *gorm.DB) ([]*DomainCreatorObject, []*DomainCreatorRole, []*DomainCreatorPolicy) {
	var object []*DomainCreatorObject
	var role []*DomainCreatorRole
	var policy []*DomainCreatorPolicy
	if err := db.Transaction(func(tx *gorm.DB) error {
		if err1 := tx.Find(&object).Error; err1 != nil {
			return err1
		}
		if err1 := tx.Find(&role).Error; err1 != nil {
			return err1
		}
		if err1 := tx.Find(&policy).Error; err1 != nil {
			return err1
		}
		return nil
	}); err != nil {
		return nil, nil, nil
	}
	return object, role, policy
}

func isValidSnapshot(db *gorm.DB) error {
	object, role, policy := dbSnapshot(db)
	om := map[uint64]bool{}
	for _, v := range object {
		om[v.ID] = true
	}
	rm := map[uint64]bool{}
	for _, v := range role {
		rm[v.ID] = true
	}

	for _, v := range object {
		if err := checkRelativeIDInMap(v.RelativeObjectID, om); err != nil {
			return err
		}
		if err := checkRelativeIDInMap(v.RelativeParentID, om); err != nil {
			return err
		}
	}
	for _, v := range role {
		if err := checkRelativeIDInMap(v.RelativeObjectID, om); err != nil {
			return err
		}
		if err := checkRelativeIDInMap(v.RelativeParentID, rm); err != nil {
			return err
		}
	}
	for _, v := range policy {
		if err := checkRelativeIDInMap(v.RelativeObjectID, om); err != nil {
			return err
		}
		if err := checkRelativeIDInMap(v.RelativeRoleID, rm); err != nil {
			return err
		}
	}
	return nil
}

func checkRelativeIDInMap(id uint64, m map[uint64]bool) error {
	if id == 0 {
		return nil
	}
	if _, ok := m[id]; !ok {
		return ErrRelativeIDOutOfIndex
	}
	return nil
}
