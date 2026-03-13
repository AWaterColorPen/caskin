package caskin

import (
	"cmp"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/gorm-adapter/v3"
)

type server struct {
	Enforcer   IEnforcer
	DB         MetaDB
	Dictionary IDictionary

	CurrentUser   User   // for ICurrentService
	CurrentDomain Domain // for ICurrentService
}

// New creates a new caskin service instance with the given options.
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

// Filter filters source by checking permission for (u, d, action) on each element.
func Filter[T any](e IEnforcer, u User, d Domain, action Action, source []T) []T {
	var result []T
	for _, v := range source {
		if Check(e, u, d, v, action) {
			result = append(result, v)
		}
	}
	return result
}

// Check returns true if u has action permission on one in domain d.
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

// Diff returns the elements in target but not in source (add) and
// elements in source but not in target (remove).
func Diff[T cmp.Ordered](source, target []T) (add, remove []T) {
	sourceSet := make(map[T]struct{}, len(source))
	targetSet := make(map[T]struct{}, len(target))
	for _, v := range source {
		sourceSet[v] = struct{}{}
	}
	for _, v := range target {
		targetSet[v] = struct{}{}
	}
	for _, v := range target {
		if _, ok := sourceSet[v]; !ok {
			add = append(add, v)
		}
	}
	for _, v := range source {
		if _, ok := targetSet[v]; !ok {
			remove = append(remove, v)
		}
	}
	return
}

// DiffPolicy returns policies in target but not in source (add) and
// policies in source but not in target (remove).
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
