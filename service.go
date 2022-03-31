package caskin

type IService interface {
}

type ICurrentService interface {
}

type IBaseService interface {
	CreateUser(User) error
	RecoverUser(User) error
	DeleteUser(User) error
	UpdateUser(User) error
	CreateDomain(Domain) error
	RecoverDomain(Domain) error
	DeleteDomain(Domain) error
}

type IFeatureService interface {
}

type ICreatorService interface {
}
