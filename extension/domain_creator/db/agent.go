package db

import "gorm.io/gorm"

type Agent struct {
    db *gorm.DB
}

func (a *Agent) Snapshot() ([]*DomainCreatorObject, []*DomainCreatorRole, []*DomainCreatorPolicy) {
    var object []*DomainCreatorObject
    var role []*DomainCreatorRole
    var policy []*DomainCreatorPolicy
    if err := a.db.Transaction(func(tx *gorm.DB) error {
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

func (a *Agent) Create(item interface{}) error {
    if err := checkPolicy(item); err != nil {
        return err
    }
    return a.db.Create(item).Error
}

func (a *Agent) Recover(item interface{}) error {
    if err := checkPolicy(item); err != nil {
        return err
    }
    return a.db.Create(item).Error
}

func (a *Agent) Update(item interface{}) error {
    if err := checkPolicy(item); err != nil {
        return err
    }
    return a.db.Updates(item).Error
}

func (a *Agent) Delete(item interface{}, id uint64) error {
    return a.db.Delete(item, id).Error
}

func checkPolicy(item interface{}) error {
    if v, ok := item.(*DomainCreatorPolicy); ok {
        return v.IsValid()
    }
    return nil
}
