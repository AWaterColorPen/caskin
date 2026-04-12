package caskin

// IDictionary combines the feature dictionary and the creator dictionary
// into a single interface. Implementations provide the static configuration
// for features (backends, frontends) and the seed data for creator objects,
// roles, and policies.
type IDictionary interface {
	IFeatureDictionary
	ICreatorDictionary
}

// IFeatureDictionary provides read access to the feature registry that
// defines available backends (API endpoints) and frontend (UI element) guards.
type IFeatureDictionary interface {
	// GetFeature returns all registered features.
	GetFeature() ([]*Feature, error)
	// GetBackend returns all registered backend permission guards.
	GetBackend() ([]*Backend, error)
	// GetFrontend returns all registered frontend permission guards.
	GetFrontend() ([]*Frontend, error)
	// GetFeatureByKey looks up a feature by its unique key.
	GetFeatureByKey(key string) (*Feature, error)
	// GetBackendByKey looks up a backend guard by its unique key.
	GetBackendByKey(key string) (*Backend, error)
	// GetFrontendByKey looks up a frontend guard by its unique key.
	GetFrontendByKey(key string) (*Frontend, error)
	// GetPackage returns all feature packages (bundles of features).
	GetPackage() ([]*Package, error)
}

// ICreatorDictionary provides the seed data used to initialise a new domain.
// When [IBaseService.CreateDomain] is called, these objects, roles, and
// policies are created automatically.
type ICreatorDictionary interface {
	// GetCreatorObject returns the object templates to seed into a new domain.
	GetCreatorObject() ([]*CreatorObject, error)
	// GetCreatorRole returns the role templates to seed into a new domain.
	GetCreatorRole() ([]*CreatorRole, error)
	// GetCreatorPolicy returns the policy templates to seed into a new domain.
	GetCreatorPolicy() ([]*CreatorPolicy, error)
}
