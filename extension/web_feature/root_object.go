package web_feature

import (
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

func StaticFeatureRoot() caskin.CustomizedData {
	return &Feature{
		Name:        DefaultFeatureRootName,
		Description: DefaultFeatureRootDescription,
		Group:       DefaultFeatureRootGroup,
	}
}

func StaticFrontendRoot() caskin.CustomizedData {
	return &Frontend{
		Key:         DefaultFrontendRootKey,
		Type:        DefaultFrontendRootType,
		Description: DefaultFrontendRootDescription,
		Group:       DefaultFrontendRootGroup,
	}
}

func StaticBackendRoot() caskin.CustomizedData {
	return &Backend{
		Path:        DefaultBackendRootPath,
		Method:      DefaultBackendRootMethod,
		Description: DefaultBackendRootDescription,
		Group:       DefaultBackendRootGroup,
	}
}

func StaticSuperRootObject(factory caskin.ObjectFactory) caskin.Object {
	o := factory()
	o.SetName(DefaultSuperRootName)
	o.SetObjectType(caskin.ObjectTypeObject)
	return o
}

func StaticRootObject(customized caskin.CustomizedData, factory caskin.ObjectFactory) caskin.Object {
	o := factory()
	caskin.CustomizedData2Object(customized, o)
	return o
}

func GetFeatureRootObject() caskin.Object {
	return featureRootObject
}

func GetFrontendRootObject() caskin.Object {
	return frontendRootObject
}

func GetBackendRootObject() caskin.Object {
	return backendRootObject
}

func setRootID(object caskin.Object, root caskin.Object) {
	if object.GetParentID() == 0 {
		object.SetParentID(root.GetID())
	}
	if object.GetObject().GetID() == 0 {
		object.SetObjectID(superRootObject.GetID())
	}
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

func ManualCreateRootObject(db caskin.MetaDB, factory caskin.ObjectFactory, domain caskin.Domain) (err error) {
	superRootObject = StaticSuperRootObject(factory)
	superRootObject.SetDomainID(domain.GetID())
	if err = doOnce(db, superRootObject); err != nil {
		return err
	}

	featureRootObject = StaticRootObject(StaticFeatureRoot(), factory)
	frontendRootObject = StaticRootObject(StaticFrontendRoot(), factory)
	backendRootObject = StaticRootObject(StaticBackendRoot(), factory)
	for _, v := range []caskin.Object{featureRootObject, frontendRootObject, backendRootObject} {
		v.SetDomainID(domain.GetID())
		v.SetObjectID(superRootObject.GetID())
		if err = doOnce(db, v); err != nil {
			return err
		}
	}
	return
}

func doOnce(db caskin.MetaDB, item interface{}) (err error) {
	for j := 0; j < 3; j++ {
		if err = db.Take(item); err == nil {
			return
		}
		if err = db.Upsert(item); err == nil {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	return
}
