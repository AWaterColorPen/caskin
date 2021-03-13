package web_feature

// package handler
//
// import (
//     "fmt"
//     "github.com/casbin/casbin/v2"
//     "strings"
//     "time"
//
//     "github.com/tencentad/martech/api/types"
//     "github.com/tencentad/martech/pkg/orm"
//     "github.com/gin-gonic/gin"
//     "gorm.io/gorm"
// )
//
//
// // BackendPostHandler 新增页面实体信息
// func BackendPostHandler(c *gin.Context) {
//     item := &types.Backend{}
//     dbEditHandler(c, item, func(db *gorm.DB) error {
//         name := fmt.Sprint(item.Path, "_", item.Method)
//         // old
//         if item.ID != 0 {
//             old := &types.Backend{ID: item.ID}
//             if err := orm.TakeBackend(db, old); err != nil {
//                 return err
//             }
//
//             if item.ObjectID != old.ObjectID {
//                 return fmt.Errorf("can't update object_id")
//             }
//
//             // update object and backend
//             return db.Transaction(func(tx *gorm.DB) error {
//                 item.Object = &types.Object{
//                     ID:   item.ObjectID,
//                     Name: name,
//                 }
//
//                 if err := orm.UpsertBackend(tx, item); err != nil {
//                     return err
//                 }
//                 return orm.UpsertObject(tx, item.Object)
//             })
//         }
//
//         // create new one
//         return db.Transaction(func(tx *gorm.DB) error {
//             o := &types.Object{
//                 Name: name,
//                 Type: types.ObjectTypeBackend,
//                 TenantID: SuperAdminDomainID,
//             }
//             item.Object = o
//
//             return orm.UpsertBackend(tx, item)
//         })
//     })
// }
//
// // BackendDeleteHandler 删除某个页面实体
// func BackendDeleteHandler(c *gin.Context) {
//     dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
//         item := &types.Backend{ID: id}
//         if err := orm.TakeBackend(db, item); err != nil {
//             return err
//         }
//
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
//         e := GetEnforcer(c)
//
//         // delete object as item in current tenant
//         os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
//         for _, v := range os {
//             if _, err := e.RemoveNamedGroupingPolicy("g2", object, v, domain); err != nil {
//                 return err
//             }
//         }
//
//         // delete object as group in current tenant
//         if _, err := e.RemoveFilteredNamedGroupingPolicy("g2", 1, object, domain); err != nil {
//             return err
//         }
//
//         return db.Transaction(func(tx *gorm.DB) error {
//             if err := orm.DeleteObjectById(db, item.ObjectID); err != nil {
//                 return err
//             }
//             return orm.DeleteBackendById(db, id)
//         })
//     })
// }
//
// // FeatureGetHandler 拉取功能列表
// func FeatureGetHandler(c *gin.Context) {
//     dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
//         fRoot := getFeatureRoot()
//         feature, err := orm.GetAllFeature(db)
//         if err != nil {
//             return nil, err
//         }
//         fm := map[uint64]*types.Feature{}
//         for _, v := range feature {
//             fm[v.ObjectID] = v
//         }
//
//         e := GetEnforcer(c)
//         tree := map[uint64][]uint64{}
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         rules := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
//         for _, rule := range rules {
//             if !casbinObjectIs(rule[1]) {
//                 continue
//             }
//             x := casbinObjectDecode(rule[0]).ID
//             y := casbinObjectDecode(rule[1]).ID
//             tree[x] = append(tree[x], y)
//         }
//
//         frontend, err := orm.GetAllFrontend(db)
//         if err != nil {
//             return nil, err
//         }
//         for _, v := range frontend {
//             for _, u := range tree[v.ObjectID] {
//                 if m, ok := fm[u]; ok {
//                     m.Frontend = append(m.Frontend, v)
//                 }
//             }
//         }
//
//         backend, err := orm.GetAllBackend(db)
//         if err != nil {
//             return nil, err
//         }
//         for _, v := range backend {
//             for _, u := range tree[v.ObjectID] {
//                 if m, ok := fm[u]; ok {
//                     m.Backend = append(m.Backend, v)
//                 }
//             }
//         }
//
//         for _, v := range feature {
//             if v.ID != fRoot.ID {
//                 v.ParentID = fRoot.ID
//             }
//         }
//
//         return feature, nil
//     })
// }
//
// // FeaturePostHandler 新增或者编辑功能信息
// func FeaturePostHandler(c *gin.Context) {
//     item := &types.Feature{}
//     dbGetHandler(c, item, func(db *gorm.DB) (interface{}, error) {
//         froot := getFeatureRootObject()
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         e := GetEnforcer(c)
//         if item.ObjectID == froot.ID {
//             return nil, fmt.Errorf("can't modify root data")
//         }
//         // old
//         if item.ID != 0 {
//             old := &types.Feature{ID: item.ID}
//             if err := orm.TakeFeature(db, old); err != nil {
//                 return nil, err
//             }
//
//             if item.ObjectID != old.ObjectID {
//                 return nil, fmt.Errorf("can't update object_id")
//             }
//
//             // update object and feature
//             if err := db.Transaction(func(tx *gorm.DB) error {
//                 item.Object = &types.Object{
//                     ID:   item.ObjectID,
//                     Name: item.Name,
//                 }
//
//                 if err := orm.UpsertFeature(tx, item); err != nil {
//                     return err
//                 }
//                 return orm.UpsertObject(tx, item.Object)
//             }); err != nil {
//                 return nil, err
//             }
//
//             if _, err := e.AddNamedGroupingPolicy("g2", casbinObjectEncode(item.Object), casbinObjectEncode(froot), domain); err != nil {
//                 return nil, err
//             }
//
//             item.CreatedAt = time.Time{}
//             item.UpdatedAt = time.Time{}
//             return item, orm.TakeFeature(db, item)
//         }
//
//         // create new one
//         if err := db.Transaction(func(tx *gorm.DB) error {
//             o := &types.Object{
//                 Name: item.Name,
//                 Type: types.ObjectTypeFeature,
//                 TenantID: SuperAdminDomainID,
//             }
//             item.Object = o
//
//             if err := orm.UpsertFeature(tx, item); err != nil {
//                 return err
//             }
//
//             o.Object = casbinObjectEncode(o)
//             return orm.UpsertObject(tx, o)
//         }); err != nil {
//             return nil, err
//         }
//
//         if _, err := e.AddNamedGroupingPolicy("g2", casbinObjectEncode(item.Object), casbinObjectEncode(froot), domain); err != nil {
//             return nil, err
//         }
//
//         item.CreatedAt = time.Time{}
//         item.UpdatedAt = time.Time{}
//         return item, orm.TakeFeature(db, item)
//     })
// }
//
// // FeatureRelateHandler 更新一个功能的关联数据
// func FeatureRelateHandler(c *gin.Context) {
//     item := &types.Feature{}
//     dbEditHandler(c, item, func(db *gorm.DB) error {
//         // verify input
//         if item.ID == 0 {
//             return ErrEmptyID
//         }
//         if err := orm.TakeFeature(db, item); err != nil {
//             return err
//         }
//
//         var fid, bid []uint64
//         for _, v := range item.Frontend {
//             fid = append(fid, v.ID)
//         }
//         for _, v := range item.Backend {
//             bid = append(bid, v.ID)
//         }
//
//         frontend, err := orm.ListFrontendById(db, fid)
//         if err != nil {
//             return err
//         }
//
//         backend, err := orm.ListBackendById(db, bid)
//         if err != nil {
//             return err
//         }
//
//         var oid []uint64
//         for _, v := range frontend {
//             oid = append(oid, v.ObjectID)
//         }
//         for _, v := range backend {
//             oid = append(oid, v.ObjectID)
//         }
//
//         // make source and target role list
//         var source, target []interface{}
//
//         e := GetEnforcer(c)
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
//         os, _ := e.GetModel()["g"]["g2"].RM.GetUsers(object, domain)
//         for _, v := range os {
//             source = append(source, v)
//         }
//
//         for _, v := range oid {
//             target = append(target, casbinObjectEncode(&types.Object{ID: v}))
//         }
//
//         // get diff to add and remove
//         add, remove := Diff(source, target)
//         for _, v := range add {
//             if _, err := e.AddNamedGroupingPolicy("g2", v, object, domain); err != nil {
//                 return err
//             }
//         }
//
//         for _, v := range remove {
//             if _, err := e.RemoveNamedGroupingPolicy("g2", v, object, domain); err != nil {
//                 return err
//             }
//         }
//
//         return nil
//     })
// }
//
// // FeatureDeleteHandler 删除某个功能实体
// func FeatureDeleteHandler(c *gin.Context) {
//     dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
//         item := &types.Feature{ID: id}
//         froot := getFeatureRootObject()
//         if item.ObjectID == froot.ID {
//             return fmt.Errorf("can't modify root data")
//         }
//
//         if err := orm.TakeFeature(db, item); err != nil {
//             return err
//         }
//
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
//         e := GetEnforcer(c)
//
//         // delete object as item in current tenant
//         os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
//         for _, v := range os {
//             if _, err := e.RemoveNamedGroupingPolicy("g2", object, v, domain); err != nil {
//                 return err
//             }
//         }
//
//         // delete object as group in current tenant
//         if _, err := e.RemoveFilteredNamedGroupingPolicy("g2", 1, object, domain); err != nil {
//             return err
//         }
//
//         return db.Transaction(func(tx *gorm.DB) error {
//             if err := orm.DeleteObjectById(db, item.ObjectID); err != nil {
//                 return err
//             }
//             return orm.DeleteFeatureById(db, id)
//         })
//     })
// }
//
// // FeatureSyncHandler 同步功能到其他租户
// func FeatureSyncHandler(c *gin.Context) {
//     dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
//         e := GetEnforcer(c)
//         // get object's tree in current tenant
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         rules := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
//         froot := getFeatureRootObject()
//         tenant, err := orm.GetAllTenant(db)
//         if err != nil {
//             return nil, err
//         }
//
//         for _, t := range tenant {
//             d := casbinTenantEncode(t)
//             set := getAllSubSon(e, casbinObjectEncode(froot), d)
//             nrs := e.GetFilteredNamedGroupingPolicy("g2", 2, d)
//
//             var olds [][]string
//             for _, rule := range nrs {
//                 if _, ok := set[rule[0]]; !ok {
//                     continue
//                 }
//                 if _, ok := set[rule[1]]; !ok {
//                     continue
//                 }
//                 olds = append(olds, rule)
//             }
//
//             var news [][]string
//             for _, rule := range rules {
//                 news = append(news, []string{rule[0], rule[1], d})
//             }
//
//             // make source and target feature g2 list
//             var source, target []interface{}
//             for _, v := range olds {
//                 source = append(source, strings.Join(v, ","))
//             }
//             for _, v := range news {
//                 target = append(target, strings.Join(v, ","))
//             }
//
//             add, remove := Diff(source, target)
//             for _, v := range add {
//                 s := strings.Split(v.(string), ",")
//                 if _, err := e.AddNamedGroupingPolicy("g2", s[0], s[1], d); err != nil {
//                     return nil, err
//                 }
//             }
//             for _, v := range remove {
//                 s := strings.Split(v.(string), ",")
//                 if _, err := e.RemoveNamedGroupingPolicy("g2", s[0], s[1], d); err != nil {
//                     return nil, err
//                 }
//             }
//         }
//
//         return nil, nil
//     })
// }
//
// func getAllSubSon(e casbin.IEnforcer, name string, domain string) map[string]bool {
//     list := []string{name}
//     visit := map[string]bool{
//         name: true,
//     }
//
//     for i := 0; i < len(list); i++ {
//         ll, _ := e.GetModel()["g"]["g2"].RM.GetUsers(list[i], domain)
//         for _, v := range ll {
//             if _, ok := visit[v]; !ok {
//                 visit[v] = true
//                 list = append(list, v)
//             }
//         }
//     }
//
//     return visit
// }
//
// func getFeatureRootObject() *types.Object {
//     o := &types.Object{Name: "功能", Type: types.ObjectTypeFeature, TenantID: SuperAdminDomainID}
//     db := orm.GetDB()
//     if err := orm.TakeObject(db, o); err == nil {
//         return o
//     }
//     _ = orm.UpsertObject(db, o)
//     return o
// }
//
// func getFeatureRoot() *types.Feature {
//     o := getFeatureRootObject()
//     f := &types.Feature{Name: "功能", Description: "功能根节点", ObjectID: o.ID}
//     db := orm.GetDB()
//     if err := orm.TakeFeature(db, f); err == nil {
//         return f
//     }
//     _ = orm.UpsertFeature(db, f)
//     return f
// }
//
// package handler
//
// import (
// "fmt"
// "strings"
// "time"
//
// "github.com/tencentad/martech/api/types"
// "github.com/tencentad/martech/pkg/orm"
// "github.com/casbin/casbin/v2"
// "github.com/gin-gonic/gin"
// "gorm.io/gorm"
// )
//
// // FeatureGetHandler 拉取功能列表
// func FeatureGetHandler(c *gin.Context) {
//     dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
//         fRoot := getFeatureRoot()
//         feature, err := orm.GetAllFeature(db)
//         if err != nil {
//             return nil, err
//         }
//         fm := map[uint64]*types.Feature{}
//         for _, v := range feature {
//             fm[v.ObjectID] = v
//         }
//
//         e := GetEnforcer(c)
//         tree := map[uint64][]uint64{}
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         rules := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
//         for _, rule := range rules {
//             if !casbinObjectIs(rule[1]) {
//                 continue
//             }
//             x := casbinObjectDecode(rule[0]).ID
//             y := casbinObjectDecode(rule[1]).ID
//             tree[x] = append(tree[x], y)
//         }
//
//         frontend, err := orm.GetAllFrontend(db)
//         if err != nil {
//             return nil, err
//         }
//         for _, v := range frontend {
//             for _, u := range tree[v.ObjectID] {
//                 if m, ok := fm[u]; ok {
//                     m.Frontend = append(m.Frontend, v)
//                 }
//             }
//         }
//
//         backend, err := orm.GetAllBackend(db)
//         if err != nil {
//             return nil, err
//         }
//         for _, v := range backend {
//             for _, u := range tree[v.ObjectID] {
//                 if m, ok := fm[u]; ok {
//                     m.Backend = append(m.Backend, v)
//                 }
//             }
//         }
//
//         for _, v := range feature {
//             if v.ID != fRoot.ID {
//                 v.ParentID = fRoot.ID
//             }
//         }
//
//         return feature, nil
//     })
// }
//
// // FeaturePostHandler 新增或者编辑功能信息
// func FeaturePostHandler(c *gin.Context) {
//     item := &types.Feature{}
//     dbGetHandler(c, item, func(db *gorm.DB) (interface{}, error) {
//         froot := getFeatureRootObject()
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         e := GetEnforcer(c)
//         if item.ObjectID == froot.ID {
//             return nil, fmt.Errorf("can't modify root data")
//         }
//         // old
//         if item.ID != 0 {
//             old := &types.Feature{ID: item.ID}
//             if err := orm.TakeFeature(db, old); err != nil {
//                 return nil, err
//             }
//
//             if item.ObjectID != old.ObjectID {
//                 return nil, fmt.Errorf("can't update object_id")
//             }
//
//             // update object and feature
//             if err := db.Transaction(func(tx *gorm.DB) error {
//                 item.Object = &types.Object{
//                     ID:   item.ObjectID,
//                     Name: item.Name,
//                 }
//
//                 if err := orm.UpsertFeature(tx, item); err != nil {
//                     return err
//                 }
//                 return orm.UpsertObject(tx, item.Object)
//             }); err != nil {
//                 return nil, err
//             }
//
//             if _, err := e.AddNamedGroupingPolicy("g2", casbinObjectEncode(item.Object), casbinObjectEncode(froot), domain); err != nil {
//                 return nil, err
//             }
//
//             item.CreatedAt = time.Time{}
//             item.UpdatedAt = time.Time{}
//             return item, orm.TakeFeature(db, item)
//         }
//
//         // create new one
//         if err := db.Transaction(func(tx *gorm.DB) error {
//             o := &types.Object{
//                 Name: item.Name,
//                 Type: types.ObjectTypeFeature,
//                 TenantID: SuperAdminDomainID,
//             }
//             item.Object = o
//
//             if err := orm.UpsertFeature(tx, item); err != nil {
//                 return err
//             }
//
//             o.Object = casbinObjectEncode(o)
//             return orm.UpsertObject(tx, o)
//         }); err != nil {
//             return nil, err
//         }
//
//         if _, err := e.AddNamedGroupingPolicy("g2", casbinObjectEncode(item.Object), casbinObjectEncode(froot), domain); err != nil {
//             return nil, err
//         }
//
//         item.CreatedAt = time.Time{}
//         item.UpdatedAt = time.Time{}
//         return item, orm.TakeFeature(db, item)
//     })
// }
//
// // FeatureRelateHandler 更新一个功能的关联数据
// func FeatureRelateHandler(c *gin.Context) {
//     item := &types.Feature{}
//     dbEditHandler(c, item, func(db *gorm.DB) error {
//         // verify input
//         if item.ID == 0 {
//             return ErrEmptyID
//         }
//         if err := orm.TakeFeature(db, item); err != nil {
//             return err
//         }
//
//         var fid, bid []uint64
//         for _, v := range item.Frontend {
//             fid = append(fid, v.ID)
//         }
//         for _, v := range item.Backend {
//             bid = append(bid, v.ID)
//         }
//
//         frontend, err := orm.ListFrontendById(db, fid)
//         if err != nil {
//             return err
//         }
//
//         backend, err := orm.ListBackendById(db, bid)
//         if err != nil {
//             return err
//         }
//
//         var oid []uint64
//         for _, v := range frontend {
//             oid = append(oid, v.ObjectID)
//         }
//         for _, v := range backend {
//             oid = append(oid, v.ObjectID)
//         }
//
//         // make source and target role list
//         var source, target []interface{}
//
//         e := GetEnforcer(c)
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
//         os, _ := e.GetModel()["g"]["g2"].RM.GetUsers(object, domain)
//         for _, v := range os {
//             source = append(source, v)
//         }
//
//         for _, v := range oid {
//             target = append(target, casbinObjectEncode(&types.Object{ID: v}))
//         }
//
//         // get diff to add and remove
//         add, remove := Diff(source, target)
//         for _, v := range add {
//             if _, err := e.AddNamedGroupingPolicy("g2", v, object, domain); err != nil {
//                 return err
//             }
//         }
//
//         for _, v := range remove {
//             if _, err := e.RemoveNamedGroupingPolicy("g2", v, object, domain); err != nil {
//                 return err
//             }
//         }
//
//         return nil
//     })
// }
//
// // FeatureDeleteHandler 删除某个功能实体
// func FeatureDeleteHandler(c *gin.Context) {
//     dbDeleteHandler(c, func(db *gorm.DB, id uint64) error {
//         item := &types.Feature{ID: id}
//         froot := getFeatureRootObject()
//         if item.ObjectID == froot.ID {
//             return fmt.Errorf("can't modify root data")
//         }
//
//         if err := orm.TakeFeature(db, item); err != nil {
//             return err
//         }
//
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         object := casbinObjectEncode(&types.Object{ID: item.ObjectID})
//         e := GetEnforcer(c)
//
//         // delete object as item in current tenant
//         os, _ := e.GetModel()["g"]["g2"].RM.GetRoles(object, domain)
//         for _, v := range os {
//             if _, err := e.RemoveNamedGroupingPolicy("g2", object, v, domain); err != nil {
//                 return err
//             }
//         }
//
//         // delete object as group in current tenant
//         if _, err := e.RemoveFilteredNamedGroupingPolicy("g2", 1, object, domain); err != nil {
//             return err
//         }
//
//         return db.Transaction(func(tx *gorm.DB) error {
//             if err := orm.DeleteObjectById(db, item.ObjectID); err != nil {
//                 return err
//             }
//             return orm.DeleteFeatureById(db, id)
//         })
//     })
// }
//
// // FeatureSyncHandler 同步功能到其他租户
// func FeatureSyncHandler(c *gin.Context) {
//     dbListHandler(c, func(db *gorm.DB) (interface{}, error) {
//         e := GetEnforcer(c)
//         // get object's tree in current tenant
//         domain := casbinTenantEncode(&types.Tenant{ID: SuperAdminDomainID})
//         rules := e.GetFilteredNamedGroupingPolicy("g2", 2, domain)
//         froot := getFeatureRootObject()
//         tenant, err := orm.GetAllTenant(db)
//         if err != nil {
//             return nil, err
//         }
//
//         for _, t := range tenant {
//             d := casbinTenantEncode(t)
//             set := getAllSubSon(e, casbinObjectEncode(froot), d)
//             nrs := e.GetFilteredNamedGroupingPolicy("g2", 2, d)
//
//             var olds [][]string
//             for _, rule := range nrs {
//                 if _, ok := set[rule[0]]; !ok {
//                     continue
//                 }
//                 if _, ok := set[rule[1]]; !ok {
//                     continue
//                 }
//                 olds = append(olds, rule)
//             }
//
//             var news [][]string
//             for _, rule := range rules {
//                 news = append(news, []string{rule[0], rule[1], d})
//             }
//
//             // make source and target feature g2 list
//             var source, target []interface{}
//             for _, v := range olds {
//                 source = append(source, strings.Join(v, ","))
//             }
//             for _, v := range news {
//                 target = append(target, strings.Join(v, ","))
//             }
//
//             add, remove := Diff(source, target)
//             for _, v := range add {
//                 s := strings.Split(v.(string), ",")
//                 if _, err := e.AddNamedGroupingPolicy("g2", s[0], s[1], d); err != nil {
//                     return nil, err
//                 }
//             }
//             for _, v := range remove {
//                 s := strings.Split(v.(string), ",")
//                 if _, err := e.RemoveNamedGroupingPolicy("g2", s[0], s[1], d); err != nil {
//                     return nil, err
//                 }
//             }
//         }
//
//         return nil, nil
//     })
// }
//
// func getAllSubSon(e casbin.IEnforcer, name string, domain string) map[string]bool {
//     list := []string{name}
//     visit := map[string]bool{
//         name: true,
//     }
//
//     for i := 0; i < len(list); i++ {
//         ll, _ := e.GetModel()["g"]["g2"].RM.GetUsers(list[i], domain)
//         for _, v := range ll {
//             if _, ok := visit[v]; !ok {
//                 visit[v] = true
//                 list = append(list, v)
//             }
//         }
//     }
//
//     return visit
// }
//
// func getFeatureRootObject() *types.Object {
//     o := &types.Object{Name: "功能", Type: types.ObjectTypeFeature, TenantID: SuperAdminDomainID}
//     db := orm.GetDB()
//     if err := orm.TakeObject(db, o); err == nil {
//         return o
//     }
//     _ = orm.UpsertObject(db, o)
//     return o
// }
//
// func getFeatureRoot() *types.Feature {
//     o := getFeatureRootObject()
//     f := &types.Feature{Name: "功能", Description: "功能根节点", ObjectID: o.ID}
//     db := orm.GetDB()
//     if err := orm.TakeFeature(db, f); err == nil {
//         return f
//     }
//     _ = orm.UpsertFeature(db, f)
//     return f
// }
