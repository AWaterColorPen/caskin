package web_feature

import (
	"time"

	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

type Root struct {
	Super    caskin.Object `json:"super"`
	Feature  caskin.Object `json:"feature"`
	Frontend caskin.Object `json:"frontend"`
	Backend  caskin.Object `json:"backend"`
}

func (r *Root) GetFeatureRootObject() caskin.Object {
	return r.Feature
}

func (r *Root) GetFrontendRootObject() caskin.Object {
	return r.Frontend
}

func (r *Root) GetBackendRootObject() caskin.Object {
	return r.Backend
}

func (r *Root) SetFeatureRoot(object caskin.Object) {
	r.setRootID(object, r.Feature)
}

func (r *Root) SetFrontendRoot(object caskin.Object) {
	r.setRootID(object, r.Frontend)
}

func (r *Root) SetBackendRoot(object caskin.Object) {
	r.setRootID(object, r.Backend)
}

func (r *Root) setRootID(object caskin.Object, root caskin.Object) {
	if object.GetParentID() == 0 {
		object.SetParentID(root.GetID())
	}
	object.SetObjectID(r.Super.GetID())
}

func InitRootObject(db caskin.MetaDB, factory caskin.ObjectFactory, domain caskin.Domain) (*Root, error) {
	r := &Root{}
	r.Super = staticSuperRootObject(factory)
	r.Super.SetDomainID(domain.GetID())
	if err := doOnceSuperRoot(db, r.Super); err != nil {
		return nil, err
	}

	r.Feature = staticRootObject(staticFeatureRoot(), factory)
	r.Frontend = staticRootObject(staticFrontendRoot(), factory)
	r.Backend = staticRootObject(staticBackendRoot(), factory)
	for _, v := range []caskin.Object{r.Feature, r.Frontend, r.Backend} {
		v.SetDomainID(domain.GetID())
		v.SetObjectID(r.Super.GetID())
	}

	if err := doOnceCustomizedData(db, r.Feature, staticFeatureRoot(), factory); err != nil {
		return nil, err
	}
	if err := doOnceCustomizedData(db, r.Frontend, staticFrontendRoot(), factory); err != nil {
		return nil, err
	}
	if err := doOnceCustomizedData(db, r.Backend, staticBackendRoot(), factory); err != nil {
		return nil, err
	}
	return r, nil
}

func staticFeatureRoot() caskin.CustomizedData {
	return &Feature{
		Name:        DefaultFeatureRootName,
		Description: DefaultFeatureRootDescription,
		Group:       DefaultFeatureRootGroup,
	}
}

func staticFrontendRoot() caskin.CustomizedData {
	return &Frontend{
		Key:         DefaultFrontendRootKey,
		Type:        DefaultFrontendRootType,
		Description: DefaultFrontendRootDescription,
		Group:       DefaultFrontendRootGroup,
	}
}

func staticBackendRoot() caskin.CustomizedData {
	return &Backend{
		Path:        DefaultBackendRootPath,
		Method:      DefaultBackendRootMethod,
		Description: DefaultBackendRootDescription,
		Group:       DefaultBackendRootGroup,
	}
}

func staticSuperRootObject(factory caskin.ObjectFactory) caskin.Object {
	o := factory()
	o.SetName(DefaultSuperRootName)
	o.SetObjectType(caskin.ObjectTypeObject)
	return o
}

func staticRootObject(customized caskin.CustomizedData, factory caskin.ObjectFactory) caskin.Object {
	o := factory()
	o.SetName(customized.GetName())
	o.SetObjectType(customized.GetObjectType())
	return o
}

func doOnceSuperRoot(db caskin.MetaDB, item interface{}) (err error) {
	for j := 0; j < 3; j++ {
		if err = db.Take(item); err == nil {
			return
		}
		if err = db.Create(item); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}

func doOnceCustomizedData(db caskin.MetaDB, item caskin.Object, customized caskin.CustomizedData, factory caskin.ObjectFactory) (err error) {
	var cond []interface{}
	linq.From(customized.JSONQuery()).ToSlice(&cond)
	cond = append(cond, item)
	for j := 0; j < 3; j++ {
		item.SetCustomizedData(nil)
		o := factory()
		if err = db.First(o, cond...); err == nil {
			item.SetID(o.GetID())
			item.SetParentID(o.GetParentID())
			item.SetCustomizedData(o.GetCustomizedData())
			return
		}

		caskin.CustomizedData2Object(customized, item)
		if err = db.Create(item); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}
