package caskin

// Default names used when no overrides are provided via [Options].
var (
	// DefaultSuperadminRoleName is the default name for the superadmin role
	// stored in the casbin policy. Override via [Options.DefaultSuperadminRoleName].
	DefaultSuperadminRoleName = "superadmin_role"
	// DefaultSuperadminDomainName is the default name for the superadmin domain.
	// Override via [Options.DefaultSuperadminDomainName].
	DefaultSuperadminDomainName = "superadmin_domain"
)

// Option is a functional option that mutates an [Options] struct.
// Pass one or more Options to [New] to customise the service.
type Option func(*Options)

// Options holds the configuration for creating a caskin [IService] via [New].
type Options struct {
	// DefaultSuperadminDomainName overrides the built-in superadmin domain name
	// (default: "superadmin_domain"). Must match the name used in the database.
	DefaultSuperadminDomainName string `json:"default_superadmin_domain_name"`
	// DefaultSuperadminRoleName overrides the built-in superadmin role name
	// (default: "superadmin_role"). Must match the name stored in casbin.
	DefaultSuperadminRoleName string `json:"default_superadmin_role_name"`
	// Dictionary configures the feature/creator dictionary backend. If nil,
	// an empty in-memory dictionary is used.
	Dictionary *DictionaryOption `json:"dictionary"`
	// DB configures the metadata database connection.
	DB *DBOption `json:"db"`
	// Watcher configures the optional casbin policy watcher (e.g. Redis).
	// When nil, no watcher is set up.
	Watcher *WatcherOption `json:"watcher"`
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}
