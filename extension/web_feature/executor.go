package web_feature

import "github.com/awatercolorpen/caskin"

type Executor struct {
	e             *caskin.Executor
	objectFactory caskin.ObjectFactory
	options       Options
}

func (e *Executor) operationPermissionCheck() error {

	return nil
}

func (e *Executor) staticFeatureRootObject() caskin.Object {
	root := &Feature{
		Name:        DefaultFeatureRootName,
		Description: DefaultFeatureRootDescriptionName,
		Group:       DefaultFeatureRootGroupName,
	}
	o := e.objectFactory()
	o.SetDomainID(e.options.Domain.GetID())
	caskin.CustomizedData2Object(root, o)
	return o
}
