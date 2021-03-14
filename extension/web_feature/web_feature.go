package web_feature

import (
	"sync"

	"github.com/awatercolorpen/caskin"
)

var (
	once              sync.Once
	featureRootObject caskin.Object
)

type WebFeature struct {
	caskin  *caskin.Caskin
	options *Options
}

func (w *WebFeature) GetExecutor(provider caskin.CurrentProvider) *Executor {
	e := w.caskin.GetExecutor(provider)
	return &Executor{
		e:             e,
		objectFactory: w.caskin.GetOptions().EntryFactory.NewObject,
	}
}

func (w *WebFeature) operationDomain() caskin.Domain {
	if w.options.Domain == nil {
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
		featureRootObject = staticFeatureRootObject(factory)
		domain := w.operationDomain()
		featureRootObject.SetDomainID(domain.GetID())
		db := w.caskin.GetOptions().MetaDB
		if err = db.Take(featureRootObject); err == nil {
			return
		}
		err = db.Upsert(featureRootObject)
	})
	return
}

func GetFeatureRootObject() caskin.Object {
	return featureRootObject
}

func staticFeatureRootObject(factory caskin.ObjectFactory) caskin.Object {
	root := &Feature{
		Name:        DefaultFeatureRootName,
		Description: DefaultFeatureRootDescriptionName,
		Group:       DefaultFeatureRootGroupName,
	}
	o := factory()
	caskin.CustomizedData2Object(root, o)
	return o
}

func setFeatureParentID(object caskin.Object) {
	if object.GetParentID() == 0 {
		root := GetFeatureRootObject()
		object.SetParentID(root.GetID())
	}
}
