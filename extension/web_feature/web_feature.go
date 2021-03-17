package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type WebFeature struct {
	caskin  *caskin.Caskin
	options *Options
}

func (w *WebFeature) GetExecutor(provider caskin.CurrentProvider) *Executor {
	e := w.caskin.GetExecutor(provider)
	return &Executor{
		e:               e,
		objectFactory:   w.caskin.GetOptions().EntryFactory.NewObject,
		operationDomain: w.operationDomain(),
	}
}

func (w *WebFeature) operationDomain() caskin.Domain {
	if w.options == nil || w.options.Domain == nil {
		return w.caskin.GetOptions().GetSuperadminDomain()
	}
	return w.options.Domain
}

func New(caskin *caskin.Caskin, options *Options) (w *WebFeature, err error) {
	w = &WebFeature{
		caskin:  caskin,
		options: options,
	}
	once.Do(func() {
		factory := w.caskin.GetOptions().EntryFactory.NewObject
		domain := w.operationDomain()
		db := w.caskin.GetOptions().MetaDB
		err = ManualCreateRootObject(db, factory, domain)
	})
	return
}