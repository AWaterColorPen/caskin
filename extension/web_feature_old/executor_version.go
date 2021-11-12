package web_feature_old

import (
	"github.com/awatercolorpen/caskin"
)

func (e *Executor) BuildVersion() error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	dump, err := e.Dump()
	if err != nil {
		return err
	}
	metadata := dump.ToRelation()
	version := &WebFeatureVersion{
		SHA256:   metadata.Version(),
		MetaData: metadata,
	}
	return e.e.DB.Create(version)
}

func (e *Executor) GetVersion() ([]*WebFeatureVersion, error) {
	if err := e.versionPermissionCheck(); err != nil {
		return nil, err
	}
	var versions []*WebFeatureVersion
	return versions, e.e.DB.Find(&versions)
}

func (e *Executor) GetLatestVersion() (*WebFeatureVersion, error) {
	versions, err := e.GetVersion()
	if err != nil {
		return nil, err
	}
	if len(versions) == 0 {
		return nil, caskin.ErrNotExists
	}
	return versions[len(versions)-1], nil
}

func (e *Executor) SyncLatestVersionToAllDomain() error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	version, err := e.GetLatestVersion()
	if err != nil {
		return err
	}
	domains, err := e.e.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err := e.SyncVersionToOneDomain(version, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) SyncVersionToAllDomain(version *WebFeatureVersion) error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	if err := e.isValidVersion(version); err != nil {
		return err
	}
	domains, err := e.e.GetAllDomain()
	if err != nil {
		return err
	}
	for _, v := range domains {
		if err := e.syncVersionToOneDomain(version, v); err != nil {
			return err
		}
	}
	return nil
}

func (e *Executor) SyncVersionToOneDomain(version *WebFeatureVersion, domain caskin.Domain) error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	if err := e.isValidVersion(version); err != nil {
		return err
	}
	return e.syncVersionToOneDomain(version, domain)
}

func (e *Executor) syncVersionToOneDomain(version *WebFeatureVersion, domain caskin.Domain) error {
	if v, ok := domain.(VersionedDomain); ok {
		if v.GetVersion() == version.SHA256 {
			return nil
		}
	}
	if err := e.syncVersionToOneDomainInternal(version, domain); err != nil {
		return err
	}
	if v, ok := domain.(VersionedDomain); ok {
		v.SetVersion(version.SHA256)
		return e.e.DB.Update(v)
	}
	return nil
}

func (e *Executor) syncVersionToOneDomainInternal(version *WebFeatureVersion, domain caskin.Domain) error {
	relations := e.allWebFeatureRelation(domain)
	var source, target []interface{}
	for k, relation := range relations {
		for _, v := range relation {
			source = append(source, relationEncode(k, v))
		}
	}
	targetRelations := version.MetaData.MergedRelation()
	for k, relation := range targetRelations {
		for _, v := range relation {
			target = append(target, relationEncode(k, v))
		}
	}

	var toDelete []uint64
	for k := range relations {
		if _, ok := targetRelations[k]; !ok {
			toDelete = append(toDelete, k)
		}
	}
	for _, v := range toDelete {
		o := e.objectFactory()
		o.SetID(v)
		if err := e.e.Enforcer.RemoveObjectInDomain(o, domain); err != nil {
			return err
		}
	}

	graph := caskin.MergedInheritanceRelations(relations, targetRelations)
	index := caskin.TopSort(graph)
	sorter := caskin.NewEdgeSorter(index)

	add, remove := caskin.Diff(source, target)
	if err := relationsAction(add, sorter.RootFirstSort, domain, e.objectFactory, e.e.Enforcer.AddParentForObjectInDomain); err != nil {
		return err
	}

	if err := relationsAction(remove, sorter.LeafFirstSort, domain, e.objectFactory, e.e.Enforcer.RemoveParentForObjectInDomain); err != nil {
		return err
	}

	return nil
}

func (e *Executor) isValidVersion(version *WebFeatureVersion) error {
	if err := e.e.DB.Take(version); err != nil {
		return err
	}
	dump, err := e.Dump()
	if err != nil {
		return err
	}
	return version.IsCompatible(dump)
}

func (e *Executor) versionPermissionCheck() error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.IsSuperadminCheck()
}
