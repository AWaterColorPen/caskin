package caskin

// GetDirectory
// get choose directory
// 1. current user has read of object to get directory
// 2. get object by type
func (s *server) GetDirectory(user User, domain Domain, ty ObjectType) ([]*Directory, error) {
	// object, err := s.GetObject(user, domain, Read, ty)
	// if err != nil {
	// 	return nil, err
	// }
	return nil, nil
}

// func ObjectToDirectoryWithItemCount(db *gorm.DB, ty string, object []Object) ([]*Directory, error) {
// 	// models, _, err := itemModelsByItemType(ty)
// 	// if err != nil {
// 	// 	return nil, err
// 	// }
// 	id := ID(object)
// 	var itemCounts []map[uint64]uint64
// 	for _, model := range models {
// 		count, err := countDirectoryItem(db, model, id)
// 		if err != nil {
// 			return nil, err
// 		}
// 		itemCounts = append(itemCounts, count)
// 	}
// 	return object2Directory(object, itemCounts...), nil
// }

func object2Directory(object []Object, itemCounts ...map[uint64]uint64) []*Directory {
	var out []*Directory
	for _, v := range object {
		u := &Directory{Object: v}
		for _, count := range itemCounts {
			u.TopItemCount += count[u.GetID()]
		}
		out = append(out, u)
	}
	return out
}

type CountDirectoryItem = func(objectID []uint64) (map[uint64]uint64, error)
type DeleteDirectoryItem = func(objectID []uint64) error

// func countDirectoryItem(db *gorm.DB, model any, objectID []uint64) (map[uint64]uint64, error) {
// 	rows, err := db.
// 		Model(model).
// 		Select("object_id, COUNT(*) as count").
// 		Where("object_id in (?)", objectID).
// 		Group("object_id").Rows()
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	out := map[uint64]uint64{}
// 	for rows.Next() {
// 		var id, count uint64
// 		if err = rows.Scan(&id, &count); err != nil {
// 			return nil, err
// 		}
// 		out[id] = count
// 	}
//
// 	return out, nil
// }
//
// func deleteDirectoryItem(db *gorm.DB, model any, objectID []uint64) error {
// 	return db.Delete(model, "object_id in (?)", objectID).Error
// }
