package web_feature

import "github.com/awatercolorpen/caskin"

func (e *Executor) GetRelation() (Relations, error) {
	return nil, nil
}

func (e *Executor) GetRelationByFeature(feature *Feature, object caskin.Object) (Relation, error) {
	return nil, nil
}

func (e *Executor) ModifyRelationPerFeature(feature *Feature, object caskin.Object, relation Relation) error {
	return nil
}

func (e *Executor) SyncRelationsToOneDomain(relations Relations, domain caskin.Domain) error {
	return nil
}

func (e *Executor) SyncRelationsToAllDomain() error {
	return nil
}