package caskin

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
	"golang.org/x/exp/constraints"
)

type server struct {
	Enforcer   IEnforcer
	DB         MetaDB
	Dictionary IDictionary

	CurrentUser   User   // for ICurrentService
	CurrentDomain Domain // for ICurrentService
}

func New(options *Options, opts ...Option) (IService, error) {
	options = options.newOptions(opts...)
	// set default caskin option
	if options.DefaultSuperadminDomainName != "" {
		DefaultSuperadminDomainName = options.DefaultSuperadminDomainName
	}
	if options.DefaultSuperadminRoleName != "" {
		DefaultSuperadminRoleName = options.DefaultSuperadminRoleName
	}

	dictionary, err := NewDictionary(options.Dictionary)
	if err != nil {
		return nil, err
	}
	db, err := options.DB.NewDB()
	if err != nil {
		return nil, err
	}
	model, err := CasbinModel()
	if err != nil {
		return nil, err
	}
	adapter, err := gormadapter.NewAdapterByDB(db)
	if err != nil {
		return nil, err
	}
	syncedEnforcer, err := casbin.NewSyncedEnforcer(model, adapter)
	if err != nil {
		return nil, err
	}
	if err = SetWatcher(syncedEnforcer, options.Watcher); err != nil {
		return nil, err
	}

	return &server{
		Enforcer:   NewEnforcer(syncedEnforcer, DefaultFactory()),
		DB:         DefaultFactory().MetadataDB(db),
		Dictionary: dictionary,
	}, nil
}

// Filter do filter source permission by u, d, action
func Filter[T any](e IEnforcer, u User, d Domain, action Action, source []T) []T {
	var result []T
	for _, v := range source {
		if Check(e, u, d, v, action) {
			result = append(result, v)
		}
	}
	return result
}

// Check object/object_data permission by u, d, action
func Check[T any](e IEnforcer, u User, d Domain, one T, action Action) bool {
	if data, ok := any(one).(ObjectData); ok {
		o := DefaultFactory().NewObject()
		o.SetID(data.GetObjectID())
		ok, _ = e.Enforce(u, o, d, action)
		return ok
	}
	if o, ok := any(one).(Object); ok {
		ok, _ = e.Enforce(u, o, d, action)
		return ok
	}
	return false
}

// Diff do diff source, target list to get add, remove list
func Diff[T constraints.Ordered](source, target []T) (add, remove []T) {
	linq.From(source).Except(linq.From(target)).ToSlice(&remove)
	linq.From(target).Except(linq.From(source)).ToSlice(&add)
	return
}

// DiffPolicy diff policy source, target list to get add, remove list
func DiffPolicy(source, target []*Policy) (add, remove []*Policy) {
	sourceMap := make(map[any]*Policy)
	targetMap := make(map[any]*Policy)
	for _, v := range source {
		sourceMap[v.Key()] = v
	}
	for _, v := range target {
		targetMap[v.Key()] = v
	}

	for _, v := range source {
		if _, ok := targetMap[v.Key()]; !ok {
			remove = append(remove, v)
		}
	}
	for _, v := range target {
		if _, ok := sourceMap[v.Key()]; !ok {
			add = append(add, v)
		}
	}
	return
}
