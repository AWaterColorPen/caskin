package caskin

var (
	DefaultSuperadminRoleName   = "superadmin_role"
	DefaultSuperadminDomainName = "superadmin_domain"
)

type Option func(*Options)

// Options configuration for caskin
type Options struct {
	// default caskin option
	DefaultSuperadminDomainName string            `json:"default_superadmin_domain_name"`
	DefaultSuperadminRoleName   string            `json:"default_superadmin_role_name"`
	Dictionary                  *DictionaryOption `json:"dictionary"`
	DB                          *DBOption         `json:"db"`
	Watcher                     *WatcherOption    `json:"watcher"`
}

func (o *Options) newOptions(opts ...Option) *Options {
	for _, v := range opts {
		v(o)
	}
	return o
}
