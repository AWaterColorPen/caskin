package caskin

import "github.com/ahmetb/go-linq/v3"

func (s *server) DefaultObjectUpdater() *treeNodeUpdater {
	return NewTreeNodeUpdater(s.DefaultObjectParentGetFunc(), s.DefaultObjectParentAddFunc(), s.DefaultObjectParentDelFunc())
}

func (s *server) DefaultObjectDeleteFunc() TreeNodeDeleteFunc {
	return func(p treeNode, d Domain) error {
		if err := s.Enforcer.RemoveObjectInDomain(p.(Object), d); err != nil {
			return err
		}
		return s.DB.DeleteByID(p, p.GetID())
	}
}

func (s *server) DefaultObjectChildrenGetFunc() TreeNodeChildrenGetFunc {
	return childrenOrParentGetFn[Object](func(p treeNode, domain Domain) []Object {
		os := s.Enforcer.GetChildrenForObjectInDomain(p.(Object), domain)
		om := IDMap(os)
		os2, _ := s.DB.GetObjectInDomain(domain, p.(Object).GetObjectType())
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

func (s *server) DefaultObjectParentGetFunc() TreeNodeParentGetFunc {
	return childrenOrParentGetFn[Object](func(p treeNode, domain Domain) []Object {
		return s.Enforcer.GetParentsForObjectInDomain(p.(Object), domain)
	})
}

func (s *server) DefaultObjectParentAddFunc() TreeNodeParentAddFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return s.Enforcer.AddParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func (s *server) DefaultObjectParentDelFunc() TreeNodeParentDelFunc {
	return func(p1 treeNode, p2 treeNode, domain Domain) error {
		return s.Enforcer.RemoveParentForObjectInDomain(p1.(Object), p2.(Object), domain)
	}
}

func childrenOrParentGetFn[T treeNode](fn func(treeNode, Domain) []T) TreeNodeChildrenGetFunc {
	return func(p treeNode, domain Domain) []treeNode {
		var out []treeNode
		linq.From(fn(p, domain)).ToSlice(&out)
		return out
	}
}
