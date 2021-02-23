package caskin

type CurrentUserProvider interface {
	Get() (User, Domain, error)
}

type management struct {
	mdb MetaDB
	e   ienforcer
}

func (m *management) GetExecutor(provider CurrentUserProvider) *executor {
	return &executor{
		mdb:      m.mdb,
		e:        m.e,
		provider: provider,
	}
}
