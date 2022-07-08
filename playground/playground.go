package playground

import (
	"path/filepath"

	"github.com/awatercolorpen/caskin"
	"github.com/awatercolorpen/caskin/example"
	"gorm.io/gorm"
)

var DictionaryDsn = "configs/caskin.toml"

// Playground example playground for easy testing
type Playground struct {
	Service    caskin.IService
	DB         *gorm.DB
	Superadmin *example.User   // superadmin user on playground
	Admin      *example.User   // a domain admin user on stage
	Member     *example.User   // a domain member user on stage
	Domain     *example.Domain // domain on playground
}

func (p *Playground) Setup() error {
	service := p.Service

	p.Domain = &example.Domain{Name: "school-1"}
	p.Superadmin = &example.User{Email: "superadmin@qq.com"}
	p.Admin = &example.User{Email: "teacher@qq.com"}
	p.Member = &example.User{Email: "student@qq.com"}

	if err := service.CreateDomain(p.Domain); err != nil {
		return err
	}
	if err := service.ResetDomain(p.Domain); err != nil {
		return err
	}
	if err := service.ResetFeature(p.Domain); err != nil {
		return err
	}
	for _, v := range []caskin.User{p.Superadmin, p.Admin, p.Member} {
		if err := service.CreateUser(v); err != nil {
			return err
		}
	}

	if err := service.AddSuperadmin(p.Superadmin); err != nil {
		return err
	}

	roles, err := service.GetRole(p.Superadmin, p.Domain)
	if err != nil {
		return err
	}
	for k, v := range map[caskin.Role][]*caskin.UserRolePair{
		roles[0]: {{User: p.Admin, Role: roles[0]}},
		roles[1]: {{User: p.Member, Role: roles[1]}},
	} {
		if err = service.ModifyUserRolePerRole(p.Superadmin, p.Domain, k, v); err != nil {
			return err
		}
	}

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
	err = db.AutoMigrate(&example.User{}, &example.Role{}, &example.Object{}, &example.Domain{}, &example.OneObjectData{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	option := &caskin.Options{
		Dictionary: &caskin.DictionaryOption{Dsn: DictionaryDsn},
		DB:         dbOption,
	}

	caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

	service, err := caskin.New(option)
	if err != nil {
		return nil, err
	}
	playground := &Playground{Service: service, DB: db}
	return playground, playground.Setup()
}

func NewPlaygroundWithSqlitePathAndWatcher(sqlitePath string, watcher *caskin.WatcherOption) (*Playground, error) {
	dbOption := &caskin.DBOption{
		DSN:  filepath.Join(sqlitePath, "sqlite"),
		Type: "sqlite",
	}
	db, err := dbOption.NewDB()
	if err != nil {
		return nil, err
	}
	err = db.AutoMigrate(&example.User{}, &example.Role{}, &example.Object{}, &example.Domain{}, &example.OneObjectData{})
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	option := &caskin.Options{
		Dictionary: &caskin.DictionaryOption{Dsn: DictionaryDsn},
		DB:         dbOption,
		Watcher:    watcher,
	}

	caskin.Register[*example.User, *example.Role, *example.Object, *example.Domain]()

	service, err := caskin.New(option)
	if err != nil {
		return nil, err
	}
	playground := &Playground{Service: service, DB: db}
	return playground, playground.Setup()
}
