package caskin

import "github.com/ahmetb/go-linq/v3"

func (e *server) DefaultObjectUpdater() *treeNodeUpdater {
	return NewTreeNodeUpdater(e.DefaultObjectParentGetFunc(), e.DefaultObjectParentAddFunc(), e.DefaultObjectParentDelFunc())
}

func (e *server) DefaultObjectDeleteFunc() TreeNodeDeleteFunc {
	return func(p treeNode, d Domain) error {
		if err := e.Enforcer.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *server) DefaultObjectChildrenGetFunc() TreeNodeChildrenGetFunc {
	return childrenOrParentGetFn[Object](func(p treeNode, domain Domain) []Object {
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

func (e *server) DefaultObjectParentGetFunc() TreeNodeParentGetFunc {
	return childrenOrParentGetFn[Object](func(p treeNode, domain Domain) []Object {
		return e.Enforcer.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *server) DefaultObjectParentAddFunc() TreeNodeParentAddFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *server) DefaultObjectParentDelFunc() TreeNodeParentDelFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *server) DefaultRoleUpdater() *treeNodeUpdater {
	return NewTreeNodeUpdater(e.DefaultRoleParentGetFunc(), e.DefaultRoleParentAddFunc(), e.DefaultRoleParentDelFunc())
}

func (e *server) DefaultRoleDeleteFunc() TreeNodeDeleteFunc {
	return func(p treeNode, d Domain) error {
		if err := e.Enforcer.RemoveRoleInDomain(p.(Role), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *server) DefaultRoleChildrenGetFunc() TreeNodeChildrenGetFunc {
	return childrenOrParentGetFn[Role](func(p treeNode, domain Domain) []Role {
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

func (e *server) DefaultRoleParentGetFunc() TreeNodeParentGetFunc {
	return childrenOrParentGetFn[Role](func(p treeNode, domain Domain) []Role {
		return e.Enforcer.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *server) DefaultRoleParentAddFunc() TreeNodeParentAddFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func (e *server) DefaultRoleParentDelFunc() TreeNodeParentDelFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return e.Enforcer.RemoveParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func childrenOrParentGetFn[T treeNode](fn func(treeNode, Domain) []T) TreeNodeChildrenGetFunc {
	return func(p treeNode, domain Domain) []treeNode {
		var out []treeNode
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}
