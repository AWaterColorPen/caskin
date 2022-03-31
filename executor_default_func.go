package caskin

import "github.com/ahmetb/go-linq/v3"

func (e *Executor) DefaultObjectUpdater() *treeNodeEntryUpdater {
	return NewTreeNodeEntryUpdater(e.DefaultObjectParentGetFunc(), e.DefaultObjectParentAddFunc(), e.DefaultObjectParentDelFunc())
}

func (e *Executor) DefaultObjectDeleteFunc() TreeNodeEntryDeleteFunc {
	return func(p TreeNodeEntry, d Domain) error {
		if err := e.Enforcer.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *Executor) DefaultObjectChildrenGetFunc() TreeNodeEntryChildrenGetFunc {
	return e.childrenOrParentGetFn(func(p TreeNodeEntry, domain Domain) any {
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

func (e *Executor) DefaultObjectParentGetFunc() TreeNodeEntryParentGetFunc {
	return e.childrenOrParentGetFn(func(p TreeNodeEntry, domain Domain) any {
		return e.Enforcer.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (e *Executor) DefaultObjectParentAddFunc() TreeNodeEntryParentAddFunc {
	return func(p1 TreeNodeEntry, p2 TreeNodeEntry, domain Domain) error {
		return e.Enforcer.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *Executor) DefaultObjectParentDelFunc() TreeNodeEntryParentDelFunc {
	return func(p1 TreeNodeEntry, p2 TreeNodeEntry, domain Domain) error {
		return e.Enforcer.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (e *Executor) DefaultRoleUpdater() *treeNodeEntryUpdater {
	return NewTreeNodeEntryUpdater(e.DefaultRoleParentGetFunc(), e.DefaultRoleParentAddFunc(), e.DefaultRoleParentDelFunc())
}

func (e *Executor) DefaultRoleDeleteFunc() TreeNodeEntryDeleteFunc {
	return func(p TreeNodeEntry, d Domain) error {
		if err := e.Enforcer.RemoveRoleInDomain(p.(Role), d); err != nil {
			return err
		}
		return e.DB.DeleteByID(p, p.GetID())
	}
}

func (e *Executor) DefaultRoleChildrenGetFunc() TreeNodeEntryChildrenGetFunc {
	return e.childrenOrParentGetFn(func(p TreeNodeEntry, domain Domain) any {
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

func (e *Executor) DefaultRoleParentGetFunc() TreeNodeEntryParentGetFunc {
	return e.childrenOrParentGetFn(func(p TreeNodeEntry, domain Domain) any {
		return e.Enforcer.GetParentsForRoleInDomain(p.(Role), domain)
	})
}

func (e *Executor) DefaultRoleParentAddFunc() TreeNodeEntryParentAddFunc {
	return func(p1 TreeNodeEntry, p2 TreeNodeEntry, domain Domain) error {
		return e.Enforcer.AddParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func (e *Executor) DefaultRoleParentDelFunc() TreeNodeEntryParentDelFunc {
	return func(p1 TreeNodeEntry, p2 TreeNodeEntry, domain Domain) error {
		return e.Enforcer.RemoveParentForRoleInDomain(p1.(Role), p2.(Role), domain)
	}
}

func (e *Executor) childrenOrParentGetFn(fn func(TreeNodeEntry, Domain) any) TreeNodeEntryChildrenGetFunc {
	return func(p TreeNodeEntry, domain Domain) []TreeNodeEntry {
		var out []TreeNodeEntry
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}
