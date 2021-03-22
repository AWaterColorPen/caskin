package web_feature

import (
	"github.com/awatercolorpen/caskin"
)

type WebFeature struct {
	caskin    *caskin.Caskin
	options   *Options
	modelText string
}

func (w *WebFeature) GetExecutor(provider caskin.CurrentProvider) *Executor {
	e := w.caskin.GetExecutor(provider)
	return &Executor{
		e:                         e,
		objectFactory:             w.caskin.GetOptions().EntryFactory.NewObject,
		operationDomain:           w.operationDomain(),
		enableBackendAPIAuthCache: w.enableBackendAPIAuthCache(),
		modelText:                 w.modelText,
	}
}

func (w *WebFeature) operationDomain() caskin.Domain {
	if w.options == nil || w.options.Domain == nil {
		return w.caskin.GetOptions().GetSuperadminDomain()
	}
	return w.options.Domain
}

func (w *WebFeature) enableBackendAPIAuthCache() bool {
	return !(w.options != nil && w.options.DisableCache)
}

func New(c *caskin.Caskin, options *Options) (w *WebFeature, err error) {
	modelText, err := caskin.CasbinModelText(c.GetOptions())
	if err != nil {
		return
	}
	w = &WebFeature{
		caskin:    c,
		options:   options,
		modelText: modelText,
	}

	err = w.caskin.GetOptions().MetaDB.AutoMigrate(&WebFeatureVersion{})
	if err != nil {
		return
	}

	once.Do(func() {
		factory := w.caskin.GetOptions().EntryFactory.NewObject
		domain := w.operationDomain()
		db := w.caskin.GetOptions().MetaDB
		err = ManualCreateRootObject(db, factory, domain)
	})

	return
}
