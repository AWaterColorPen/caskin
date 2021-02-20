package caskin

type CurrentUserProvider interface {
	Get() (User, Domain, error)
}

type Kin struct {
	mdb MetaDB
	e   ienforcer
}
