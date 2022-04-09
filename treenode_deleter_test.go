package caskin_test

import (
	"testing"
)

func TestTreeNodeEntryDeleter_RetryWithRelation(t *testing.T) {
	//stage, _ := playground.NewPlaygroundWithSqlitePath(t.TempDir())
	//
	//assert.NoError(t, stage.AddSubAdmin())
	//
	//service := stage.Service
	//
	//object1 := &example.Object{
	//	Name:     "object_sub_02",
	//	Type:     caskin.ObjectTypeObject,
	//	ObjectID: 4,
	//	ParentID: 4,
	//}
	//
	//provider.Domain = stage.Domain
	//provider.User = stage.SubAdminUser
	//assert.NoError(t, service.CreateObject(object1))
	//
	//provider.User = stage.SuperadminUser
	//list1, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, list1, 6)
	//
	//assert.NoError(t, service.Enforcer.RemoveObjectInDomain(list1[5], stage.Domain))
	//list2, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, list2, 6)
	//
	//provider.User = stage.AdminUser
	//assert.NoError(t, service.DeleteObject(list2[3]))
	//
	//provider.User = stage.SuperadminUser
	//list3, err := service.ObjectGet()
	//assert.NoError(t, err)
	//assert.Len(t, list3, 4)
}
