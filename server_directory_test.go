package caskin_test

import (
	"testing"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/awatercolorpen/caskin/playground"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func TestServer_CreateDirectory(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateDirectory(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_01", Type: "test"}
	assert.Equal(t, caskin.ErrAlreadyExists, service.CreateDirectory(stage.Admin, stage.Domain, object2))
	request1 := &caskin.DirectoryRequest{Type: "test", To: object2.GetID()}
	assert.Equal(t, caskin.ErrInValidRequest, service.DeleteDirectory(stage.Admin, stage.Domain, request1))
	request1.ActionDirectory = func([]uint64) error { return nil }
	assert.NoError(t, service.DeleteDirectory(stage.Admin, stage.Domain, request1))

	object3 := &example.Object{Name: "test_sub_01", Type: "test"}
	assert.NoError(t, service.CreateDirectory(stage.Admin, stage.Domain, object3))

}

func TestServer_UpdateDirectory(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)
	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{ID: object1.GetID(), Name: "test_sub_01_new_name", ParentID: object1.GetParentID(), Type: "test"}
	assert.NoError(t, service.UpdateDirectory(stage.Admin, stage.Domain, object2))
}

func TestServer_DeleteDirectory(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateDirectory(stage.Admin, stage.Domain, object1))

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: object1.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data1, "test"))
	request1 := &caskin.DirectoryRequest{Type: "test", To: object1.GetID()}
	request1.ActionDirectory = func(id []uint64) error {
		return stage.DB.Delete(&example.OneObjectData{}, "object_id in (?)", id).Error
	}
	assert.NoError(t, service.DeleteDirectory(stage.Admin, stage.Domain, request1))

	var list []*example.OneObjectData
	err = stage.DB.Find(&list).Error
	assert.NoError(t, err)
	assert.Len(t, list, 0)
}

func TestServer_GetDirectory(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test", ParentID: object1.GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: object1.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data1, "test"))
	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: object2.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data2, "test"))

	request1 := &caskin.DirectoryRequest{Type: "test", To: objects[1].GetID(), SearchType: "all"}
	request1.CountDirectory = countDirectory(stage.DB, &example.OneObjectData{})
	directory1, err := service.GetDirectory(stage.Member, stage.Domain, request1)
	assert.NoError(t, err)
	assert.Len(t, directory1, 2)
	assert.Equal(t, uint64(1), directory1[0].AllDirectoryCount)
	assert.Equal(t, uint64(2), directory1[0].AllItemCount)
	assert.Equal(t, uint64(1), directory1[0].TopDirectoryCount)
	assert.Equal(t, uint64(1), directory1[0].TopItemCount)
}

func TestServer_MoveDirectory(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test", ParentID: object1.GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: object1.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data1, "test"))
	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: object2.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data2, "test"))

	request1 := &caskin.DirectoryRequest{Type: "test", To: object2.GetID(), ID: []uint64{objects[1].GetID(), object1.GetID()}}
	request1.Policy = "continue"
	response1, err := service.MoveDirectory(stage.Admin, stage.Domain, request1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(2), response1.ToDoDirectoryCount)
	assert.Equal(t, uint64(0), response1.DoneDirectoryCount)
	request1.To = objects[1].GetID()
	request1.ID = []uint64{object1.GetID(), object2.GetID()}
	response2, err := service.MoveDirectory(stage.Admin, stage.Domain, request1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), response2.ToDoDirectoryCount)
	assert.Equal(t, uint64(2), response2.DoneDirectoryCount)

	request2 := &caskin.DirectoryRequest{Type: "test", To: objects[1].GetID()}
	request2.CountDirectory = countDirectory(stage.DB, &example.OneObjectData{})
	directory1, err := service.GetDirectory(stage.Admin, stage.Domain, request2)
	assert.NoError(t, err)
	assert.Len(t, directory1, 2)
	assert.Equal(t, uint64(0), directory1[0].AllDirectoryCount)
	assert.Equal(t, uint64(1), directory1[0].AllItemCount)
	assert.Equal(t, uint64(0), directory1[1].TopDirectoryCount)
	assert.Equal(t, uint64(1), directory1[1].TopItemCount)
}

func TestServer_MoveItem(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test", ParentID: object1.GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: object1.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data1, "test"))
	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: object2.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data2, "test"))

	request1 := &caskin.DirectoryRequest{Type: "test", To: object1.GetID(), ID: []uint64{data1.GetID(), data2.GetID()}}
	request1.Policy = "continue"
	response1, err := service.MoveItem(stage.Admin, stage.Domain, &example.OneObjectData{}, request1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), response1.ToDoItemCount)
	assert.Equal(t, uint64(2), response1.DoneItemCount)

	request2 := &caskin.DirectoryRequest{Type: "test", To: objects[1].GetID()}
	request2.CountDirectory = countDirectory(stage.DB, &example.OneObjectData{})
	directory1, err := service.GetDirectory(stage.Admin, stage.Domain, request2)
	assert.NoError(t, err)
	assert.Len(t, directory1, 1)
	assert.Equal(t, uint64(1), directory1[0].AllDirectoryCount)
	assert.Equal(t, uint64(2), directory1[0].AllItemCount)
}

func TestServer_CopyItem(t *testing.T) {
	stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	service := stage.Service

	objects, err := service.GetObject(stage.Admin, stage.Domain, caskin.Manage)
	assert.NoError(t, err)
	assert.Len(t, objects, 2)

	object1 := &example.Object{Name: "test_sub_01", Type: "test", ParentID: objects[1].GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object1))

	object2 := &example.Object{Name: "test_sub_02", Type: "test", ParentID: object1.GetID()}
	assert.NoError(t, service.CreateObject(stage.Admin, stage.Domain, object2))

	data1 := &example.OneObjectData{Name: "object_data_1", ObjectID: object1.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data1, "test"))
	data2 := &example.OneObjectData{Name: "object_data_2", ObjectID: object2.GetID()}
	assert.NoError(t, service.CreateObjectData(stage.Admin, stage.Domain, data2, "test"))

	request1 := &caskin.DirectoryRequest{Type: "test", To: object1.GetID(), ID: []uint64{data2.GetID()}}
	request1.Policy = "continue"
	response1, err := service.CopyItem(stage.Member, stage.Domain, &example.OneObjectData{}, request1)
	assert.NoError(t, err)
	assert.Equal(t, uint64(0), response1.ToDoItemCount)
	assert.Equal(t, uint64(1), response1.DoneItemCount)

	request2 := &caskin.DirectoryRequest{Type: "test", To: objects[1].GetID()}
	request2.CountDirectory = countDirectory(stage.DB, &example.OneObjectData{})
	directory1, err := service.GetDirectory(stage.Admin, stage.Domain, request2)
	assert.NoError(t, err)
	assert.Len(t, directory1, 1)
	assert.Equal(t, uint64(1), directory1[0].AllDirectoryCount)
	assert.Equal(t, uint64(3), directory1[0].AllItemCount)
	assert.Equal(t, uint64(2), directory1[0].TopItemCount)
}

func countDirectory(db *gorm.DB, model any) func(id []uint64) (map[uint64]uint64, error) {
	return func(id []uint64) (map[uint64]uint64, error) {
		rows, err := db.
			Model(model).
			Select("object_id, COUNT(*) as count").
			Where("object_id in (?)", id).
			Group("object_id").Rows()
		if err != nil {
			return nil, err
		}
		out := map[uint64]uint64{}
		for rows.Next() {
			var v, count uint64
			if err = rows.Scan(&v, &count); err != nil {
				return nil, err
			}
			out[v] = count
		}
		return out, nil
	}
}
