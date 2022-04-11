package playground

import (
	"path/filepath"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
)

var DictionaryDsn = "configs/caskin.toml"

// Playground example playground for easy testing
type Playground struct {
	Service        caskin.IService
	Domain         *example.Domain // a domain on stage
	SuperadminUser *example.User   // superadmin user on stage
	AdminUser      *example.User   // a domain admin user on stage
	MemberUser     *example.User   // a domain member user on stage
	SubAdminUser   *example.User   // a domain sub admin user on stage
}

func (s *Playground) AddSubAdmin() error {
	// 
	// provider.Domain = s.Domain
	// provider.User = s.AdminUser
	// executor := s.Caskin.GetExecutor(provider)
	//
	// subAdmin := &example.User{
	// 	PhoneNumber: "123456789031",
	// 	Email:       "subadmin@qq.com",
	// }
	// if err := service.UserCreate(subAdmin); err != nil {
	// 	return err
	// }
	//
	// objects, err := service.ObjectGet()
	// if err != nil {
	// 	return err
	// }
	//
	// object1 := &example.Object{
	// 	Name:     "object_sub_01",
	// 	Type:     caskin.ObjectTypeObject,
	// 	ObjectID: objects[0].GetID(),
	// 	ParentID: objects[0].GetID(),
	// }
	// if err := service.CreateObject(object1); err != nil {
	// 	return err
	// }
	// object1.ObjectID = object1.ID
	// if err := service.UpdateObject(object1); err != nil {
	// 	return err
	// }
	//
	// object2 := &example.Object{
	// 	Name:     "role_sub_02",
	// 	Type:     caskin.ObjectTypeRole,
	// 	ObjectID: object1.ID,
	// 	ParentID: objects[1].GetID(),
	// }
	// if err := service.CreateObject(object2); err != nil {
	// 	return err
	// }
	//
	// role := &example.Role{
	// 	Name:     "admin_sub_01",
	// 	ObjectID: object2.ID,
	// 	ParentID: 1,
	// }
	// if err := service.RoleCreate(role); err != nil {
	// 	return err
	// }
	//
	// for k, v := range map[caskin.Role][]*caskin.UserRolePair{
	// 	role: {{User: subAdmin, Role: role}},
	// } {
	// 	if err := service.UserRolePerRoleModify(k, v); err != nil {
	// 		return err
	// 	}
	// }
	//
	// policy := []*caskin.Policy{
	// 	{role, object1, s.Domain, caskin.Read},
	// 	{role, object1, s.Domain, caskin.Write},
	// 	{role, object2, s.Domain, caskin.Read},
	// 	{role, object2, s.Domain, caskin.Write},
	// }
	// if err := service.PolicyPerRoleModify(role, policy); err != nil {
	// 	return err
	// }

	// s.SubAdminUser = subAdmin
	return nil
}

func NewPlaygroundWithSqlitePath(sqlitePath string) (*Playground, error) {
	dbOption := &caskin.DBOption{
		DSN:  filepath.Join(sqlitePath, "sqlite"),
		Type: "sqlite",
	}

	db, err := dbOption.NewDB()
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&example.User{}, &example.Domain{}, &example.Object{}, &example.OneObjectData{})
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}

	m, err := caskin.CasbinModel()
	if err != nil {
		return nil, err
	}

	enforcer, err := casbin.NewSyncedEnforcer(m, adapter)
	if err != nil {
		return nil, err
	}
	option := &caskin.Options{
		Dictionary: &caskin.DictionaryOption{Dsn: DictionaryDsn},
		DB:         dbOption,
		Enforcer:   enforcer,
		MetaDB:     example.NewGormMDBByDB(db),
	}
	service, err := caskin.New(option)
	if err != nil {
		return nil, err
	}
	return &Playground{Service: service}, nil
}
