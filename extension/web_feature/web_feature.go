package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type Options struct {
	Caskin        *caskin.Caskin
	DomainFactory caskin.DomainFactory
	ObjectFactory caskin.ObjectFactory
	MetaDB        caskin.MetaDB
	ModelText     string
}

type WebFeature struct {
	options *Options
}

func (w *WebFeature) GetExecutor(provider caskin.CurrentProvider) *Executor {
	e := w.options.Caskin.GetExecutor(provider)
	return &Executor{
		e:               e,
		objectFactory:   w.options.ObjectFactory,
		operationDomain: w.options.DomainFactory(),
		modelText:       w.options.ModelText,
	}
}

func New(options *Options) (w *WebFeature, err error) {
	err = options.MetaDB.AutoMigrate(&WebFeatureVersion{})
	if err != nil {
		return
	}

	w = &WebFeature{options: options}
	once.Do(func() {
		factory := options.ObjectFactory
		domain := options.DomainFactory()
		db := options.MetaDB
		err = ManualCreateRootObject(db, factory, domain)
	})

	return
}
