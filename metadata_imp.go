package caskin

import (
	"gorm.io/gorm"
)

type builtinMetadataDB[U User, R Role, O Object, D Domain] struct {
	DB *gorm.DB
}

func (b *builtinMetadataDB[U, R, O, D]) Create(item any) error {
	return b.DB.Create(item).Error
}

func (b *builtinMetadataDB[U, R, O, D]) Recover(item any) error {
	if err := b.TakeUnscoped(item); err != nil {
		return err
	}
	return b.DB.Unscoped().Model(item).Update("delete_at", nil).Error
}

func (b *builtinMetadataDB[U, R, O, D]) Update(item any) error {
	return b.DB.Updates(item).Error
}

func (b *builtinMetadataDB[U, R, O, D]) UpsertType(item any) UpsertType {
	if err := b.Take(item); err == nil {
		return UpsertTypeUpdate
	}
	if err := b.TakeUnscoped(item); err == nil {
		return UpsertTypeRecover
	}
	return UpsertTypeCreate
}

func (b *builtinMetadataDB[U, R, O, D]) Take(item any) error {
	return b.DB.Where(item).Take(item).Error
}

func (b *builtinMetadataDB[U, R, O, D]) TakeUnscoped(item any) error {
	return b.DB.Unscoped().Where(item).Take(item).Error
}

func (b *builtinMetadataDB[U, R, O, D]) Find(items any, cond ...any) error {
	return b.DB.Find(items, cond...).Error
}

func (b *builtinMetadataDB[U, R, O, D]) DeleteByID(item any, id uint64) error {
	return b.DB.Delete(item, id).Error
}

func (b *builtinMetadataDB[U, R, O, D]) GetUserByID(id []uint64) ([]User, error) {
	users, err := GetByID[U](b, id)
	if err != nil {
		return nil, err
	}
	out := make([]User, len(users))
	for i, v := range users {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetRoleInDomain(domain Domain) ([]Role, error) {
	var roles []R
	if err := b.DB.Find(&roles, "domain_id = ?", domain.GetID()).Error; err != nil {
		return nil, err
	}
	out := make([]Role, len(roles))
	for i, v := range roles {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetRoleByID(id []uint64) ([]Role, error) {
	roles, err := GetByID[R](b, id)
	if err != nil {
		return nil, err
	}
	out := make([]Role, len(roles))
	for i, v := range roles {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetObjectInDomain(domain Domain, objectType ...ObjectType) ([]Object, error) {
	d := b.DB.Where("domain_id = ?", domain.GetID())
	if len(objectType) > 0 {
		d = d.Where("type IN ?", objectType)
	}

	var objects []O
	if err := d.Find(&objects).Error; err != nil {
		return nil, err
	}
	out := make([]Object, len(objects))
	for i, v := range objects {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetObjectByID(id []uint64) ([]Object, error) {
	objects, err := GetByID[O](b, id)
	if err != nil {
		return nil, err
	}
	out := make([]Object, len(objects))
	for i, v := range objects {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetDomainByID(id []uint64) ([]Domain, error) {
	domains, err := GetByID[D](b, id)
	if err != nil {
		return nil, err
	}
	out := make([]Domain, len(domains))
	for i, v := range domains {
		out[i] = v
	}
	return out, nil
}

func (b *builtinMetadataDB[U, R, O, D]) GetAllDomain() ([]Domain, error) {
	var domains []D
	if err := b.DB.Find(&domains).Error; err != nil {
		return nil, err
	}
	out := make([]Domain, len(domains))
	for i, v := range domains {
		out[i] = v
	}
	return out, nil
}

// GetByID retrieves records of type T filtered by a list of IDs.
func GetByID[T any](db MetaDB, id []uint64) ([]T, error) {
	var out []T
	if err := db.Find(&out, "id IN ?", id); err != nil {
		return nil, err
	}
	return out, nil
}
