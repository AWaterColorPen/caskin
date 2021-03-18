package web_feature

import (
	"fmt"

	"github.com/awatercolorpen/caskin"
)

func (e *Executor) BuildVersion() error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	relations := e.e.Enforcer.GetObjectInheritanceRelationInDomain(e.operationDomain)
	metadata := Relations(relations)
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

func (e *Executor) SyncVersionToOneDomain(version *WebFeatureVersion, domain caskin.Domain) error {
	if err := e.versionPermissionCheck(); err != nil {
		return err
	}
	if err := e.isValidVersion(version); err != nil {
		return err
	}
	set := e.getRootAndDescendant(GetFeatureRootObject(), domain)
	relations := e.e.Enforcer.GetObjectInheritanceRelationInDomain(domain)

	encode := func(k, v interface{}) string {
		return fmt.Sprintf("%v%v%v", k, caskin.DefaultSeparator, v)
	}
	decode := func(in interface{}) (x, y uint64, err error) {
		format := fmt.Sprintf("%%d%v%%d", caskin.DefaultSeparator)
		_, err = fmt.Sscanf(in.(string), format, &x, &y)
		return
	}

	var source, target []interface{}
	for k, relation := range relations {
		if _, ok := set[k]; !ok {
			continue
		}
		for _, v := range relation {
			if _, ok := set[v]; !ok {
				continue
			}
			source = append(source, encode(k, v))
		}
	}
	for k, relation := range version.MetaData {
		for _, v := range relation {
			target = append(target, encode(k, v))
		}
	}
	add, remove := caskin.Diff(source, target)
	for _, v := range add {
		x, y, err := decode(v)
		if err != nil {
			return err
		}
		ox, oy := e.objectFactory(), e.objectFactory()
		ox.SetID(x)
		oy.SetID(y)
		if err := e.e.Enforcer.AddParentForObjectInDomain(oy, ox, domain); err != nil {
			return err
		}
	}
	for _, v := range remove {
		x, y, err := decode(v)
		if err != nil {
			return err
		}
		ox, oy := e.objectFactory(), e.objectFactory()
		ox.SetID(x)
		oy.SetID(y)
		if err := e.e.Enforcer.RemoveParentForObjectInDomain(oy, ox, domain); err != nil {
			return err
		}
	}

	return nil
}

func (e *Executor) getRootAndDescendant(object caskin.Object, domain caskin.Domain) map[interface{}]bool {
	list := []caskin.Object{object}
	visit := map[interface{}]bool{
		object.GetID(): true,
	}

	for i := 0; i < len(list); i++ {
		ll := e.e.Enforcer.GetChildrenForObjectInDomain(list[i], domain)
		for _, v := range ll {
			if _, ok := visit[v.GetID()]; !ok {
				visit[v] = true
				list = append(list, v)
			}
		}
	}

	return visit
}

func (e *Executor) isValidVersion(version *WebFeatureVersion) error {
	return e.e.DB.Take(version)
}

func (e *Executor) versionPermissionCheck() error {
	if err := e.operationPermissionCheck(); err != nil {
		return err
	}
	return e.e.IsSuperadminCheck()
}
