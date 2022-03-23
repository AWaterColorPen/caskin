package web_feature_old

import (
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

func InitRootObject(db caskin.MetaDB, domain caskin.Domain) (*Root, error) {
	r := &Root{}

	for _, v := range []caskin.Object{r.Feature, r.Frontend, r.Backend} {
		v.SetDomainID(domain.GetID())
		v.SetObjectID(r.Super.GetID())
	}
	return r, nil
}
