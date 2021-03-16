package web_feature

import (
	"sync"
	"time"

	"github.com/awatercolorpen/caskin"
)

var (
	featureRootObject caskin.Object
	once              sync.Once
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
		err = createFeatureRootObject(db, factory, domain)
	})
	return
}

func GetFeatureRootObject() caskin.Object {
	return featureRootObject
}

func createFeatureRootObject(db caskin.MetaDB, factory caskin.ObjectFactory, domain caskin.Domain) (err error) {
	featureRootObject = staticFeatureRootObject(factory)
	featureRootObject.SetDomainID(domain.GetID())
	for i := 0; i < 3; i++ {
		if err = db.Take(featureRootObject); err == nil {
			return
		}
		if err = db.Upsert(featureRootObject); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
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
