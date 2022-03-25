package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type Options struct {
	Caskin        *caskin.Caskin
	DomainFactory caskin.DomainFactory
	ObjectFactory caskin.ObjectFactory
	MetaDB        caskin.MetaDB
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
	}
}

func New(options *Options) (*WebFeature, error) {
	// root, err := InitRootObject(options.MetaDB, options.DomainFactory())
	// if err != nil {
	// 	return nil, err
	// }
	//
	// return &WebFeature{options: options, root: root}, nil
	return nil, nil
}
