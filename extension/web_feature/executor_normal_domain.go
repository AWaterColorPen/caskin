package web_feature

import (
	"github.com/ahmetb/go-linq/v3"
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) NormalDomainGetFeature() ([]caskin.Object, error) {
	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	os := e.filterWithNoError(objects)
	linq.From(os).ToSlice(&objects)
	return objects, nil
}

func (e *Executor) NormalDomainGetPolicy() ([]caskin.Object, error) {
	objects, err := e.e.DB.GetObjectInDomain(e.operationDomain, ObjectTypeFeature)
	if err != nil {
		return nil, err
	}
	os := e.filterWithNoError(objects)
	linq.From(os).ToSlice(&objects)
	return objects, nil
}