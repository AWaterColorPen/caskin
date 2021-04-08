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
	root    *Root
}

func (w *WebFeature) GetRoot() *Root {
	return w.root
}

func (w *WebFeature) GetExecutor(provider caskin.CurrentProvider) *Executor {
	e := w.options.Caskin.GetExecutor(provider)
	return &Executor{
		e:               e,
		root:            w.root,
		objectFactory:   w.options.ObjectFactory,
		operationDomain: w.options.DomainFactory(),
		modelText:       w.options.ModelText,
	}
}

func New(options *Options) (*WebFeature, error) {
	if err := options.MetaDB.AutoMigrate(&WebFeatureVersion{}); err != nil {
		return nil, err
	}

	root, err := InitRootObject(options.MetaDB, options.ObjectFactory, options.DomainFactory())
	if err != nil {
		return nil, err
	}

	return &WebFeature{options: options, root: root}, nil
}
