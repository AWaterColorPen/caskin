package caskin

import "github.com/ahmetb/go-linq/v3"

func (e *baseService) DefaultObjectUpdater() *treeNodeUpdater {
	return NewTreeNodeUpdater(e.DefaultObjectParentGetFunc(), e.DefaultObjectParentAddFunc(), e.DefaultObjectParentDelFunc())
}

func (e *baseService) DefaultObjectDeleteFunc() TreeNodeDeleteFunc {
	return func(p treeNode, d Domain) error {
		if err := e.Enforcer.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *baseService) DefaultObjectChildrenGetFunc() TreeNodeChildrenGetFunc {
	return e.childrenOrParentGetFn(func(p treeNode, domain Domain) any {
		os := e.Enforcer.GetChildrenForObjectInDomain(p.(Object), domain)
		om := IDMap(os)
		os2, _ := e.DB.GetObjectInDomain(domain, p.(Object).GetObjectType())
		for _, v := range os2 {
			if v.GetParentID() != p.GetID() {
				continue
			}
			if _, ok := om[v.GetID()]; !ok {
				om[v.GetID()] = v
				os = append(os, v)
			}
		}
		return os
	})
}

func (e *baseService) DefaultObjectParentGetFunc() TreeNodeParentGetFunc {
	return e.childrenOrParentGetFn(func(p treeNode, domain Domain) any {
		return e.Enforcer.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *baseService) DefaultObjectParentAddFunc() TreeNodeParentAddFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *baseService) DefaultObjectParentDelFunc() TreeNodeParentDelFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *baseService) DefaultRoleUpdater() *treeNodeUpdater {
	return NewTreeNodeUpdater(e.DefaultRoleParentGetFunc(), e.DefaultRoleParentAddFunc(), e.DefaultRoleParentDelFunc())
}

func (e *baseService) DefaultRoleDeleteFunc() TreeNodeDeleteFunc {
	return func(p treeNode, d Domain) error {
		if err := e.Enforcer.RemoveRoleInDomain(p.(Role), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *baseService) DefaultRoleChildrenGetFunc() TreeNodeChildrenGetFunc {
	return e.childrenOrParentGetFn(func(p treeNode, domain Domain) any {
		rs := e.Enforcer.GetChildrenForRoleInDomain(p.(Role), domain)
		rm := IDMap(rs)
		rs2, _ := e.DB.GetRoleInDomain(domain)
		for _, v := range rs2 {
			if v.GetParentID() != p.GetID() {
				continue
			}
			if _, ok := rm[v.GetID()]; !ok {
				rm[v.GetID()] = v
				rs = append(rs, v)
			}
		}
		return rs
	})
}

func (e *baseService) DefaultRoleParentGetFunc() TreeNodeParentGetFunc {
	return e.childrenOrParentGetFn(func(p treeNode, domain Domain) any {
		return e.Enforcer.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *baseService) DefaultRoleParentAddFunc() TreeNodeParentAddFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func (e *baseService) DefaultRoleParentDelFunc() TreeNodeParentDelFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.RemoveParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func (e *baseService) childrenOrParentGetFn(fn func(treeNode, Domain) any) TreeNodeChildrenGetFunc {
	return func(p treeNode, domain Domain) []treeNode {
		var out []treeNode
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}
