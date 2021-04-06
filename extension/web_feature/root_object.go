package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"sync"
	"time"

	"github.com/awatercolorpen/caskin"
)

var (
	superRootObject    caskin.Object
	featureRootObject  caskin.Object
	frontendRootObject caskin.Object
	backendRootObject  caskin.Object
	once               sync.Once
)

func GetFeatureRootObject() caskin.Object {
	return featureRootObject
}

func GetFrontendRootObject() caskin.Object {
	return frontendRootObject
}

func GetBackendRootObject() caskin.Object {
	return backendRootObject
}

func ManualCreateRootObject(db caskin.MetaDB, factory caskin.ObjectFactory, domain caskin.Domain) (err error) {
	superRootObject = staticSuperRootObject(factory)
	superRootObject.SetDomainID(domain.GetID())
	if err = doOnceSuperRoot(db, superRootObject); err != nil {
		return err
	}

	featureRootObject = staticRootObject(staticFeatureRoot(), factory)
	frontendRootObject = staticRootObject(staticFrontendRoot(), factory)
	backendRootObject = staticRootObject(staticBackendRoot(), factory)
	for _, v := range []caskin.Object{featureRootObject, frontendRootObject, backendRootObject} {
		v.SetDomainID(domain.GetID())
		v.SetObjectID(superRootObject.GetID())
	}

	if err = doOnceCustomizedData(db, featureRootObject, staticFeatureRoot(), factory); err != nil {
		return err
	}
	if err = doOnceCustomizedData(db, frontendRootObject, staticFrontendRoot(), factory); err != nil {
		return err
	}
	if err = doOnceCustomizedData(db, backendRootObject, staticBackendRoot(), factory); err != nil {
		return err
	}
	return
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

func setRootID(object caskin.Object, root caskin.Object) {
	if object.GetParentID() == 0 {
		object.SetParentID(root.GetID())
	}
	object.SetObjectID(superRootObject.GetID())
}

func setFeatureRoot(object caskin.Object) {
	setRootID(object, GetFeatureRootObject())
}

func setFrontendRoot(object caskin.Object) {
	setRootID(object, GetFrontendRootObject())
}

func setBackendRoot(object caskin.Object) {
	setRootID(object, GetBackendRootObject())
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
