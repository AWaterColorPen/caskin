package caskin

type IDictionary interface {
	IFeatureDictionary
	ICreatorDictionary
}

type IFeatureDictionary interface {
	GetFeature() ([]*Feature, error)
	GetBackend() ([]*Backend, error)
	GetFrontend() ([]*Frontend, error)
	GetFeatureByKey(key string) (*Feature, error)
	GetBackendByKey(key string) (*Backend, error)
	GetFrontendByKey(key string) (*Frontend, error)
	GetPackage() ([]*Package, error)
}

type ICreatorDictionary interface {
	GetCreatorObject() ([]*CreatorObject, error)
	GetCreatorRole() ([]*CreatorRole, error)
	GetCreatorPolicy() ([]*CreatorPolicy, error)
}
