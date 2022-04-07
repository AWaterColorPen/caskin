package caskin

func (s *server) FeatureReload(domain Domain) error {
	sourceG2 := s.getSourceFeatureG2(domain)
	targetG2 := s.getTargetFeatureG2()
	var source, target []string
	for k, relation := range sourceG2 {
		for _, v := range relation {
			source = append(source, (&InheritanceEdge[string]{}).Encode(k, v))
		}
	}
	for k, relation := range targetG2 {
		for _, v := range relation {
			target = append(target, (&InheritanceEdge[string]{}).Encode(k, v))
		}
	}

	var toDelete []Object
	for k := range sourceG2 {
		if _, ok := targetG2[k]; !ok {
			toDelete = append(toDelete, &NamedObject{Name: k})
		}
	}
	for _, v := range toDelete {
		if err := s.Enforcer.RemoveObjectInDomain(v, domain); err != nil {
			return err
		}
	}

	graph := MergeInheritanceGraph(sourceG2, targetG2)
	index := TopSort(graph)
	sorter := NewEdgeSorter(index)

	add, remove := Diff2(source, target)
	if err := inheritanceAction(add, domain, sorter.RootFirstSort, s.Enforcer.AddParentForObjectInDomain); err != nil {
		return err
	}

	if err := inheritanceAction(remove, domain, sorter.LeafFirstSort, s.Enforcer.RemoveParentForObjectInDomain); err != nil {
		return err
	}

	return nil
}

func (s *server) getSourceFeatureG2(domain Domain) map[string][]string {
	var queue []string
	inQueue := map[string]bool{}
	for _, v := range queue {
		inQueue[v] = true
	}

	m := map[string][]string{}
	for i := 0; i < len(queue); i++ {
		m[queue[i]] = []string{}
		ll := s.Enforcer.GetChildrenForObjectInDomain(&NamedObject{Name: queue[i]}, domain)
		for _, v := range ll {
			if _, ok := inQueue[v.Encode()]; !ok {
				queue = append(queue, v.Encode())
				inQueue[v.Encode()] = true
			}
			m[queue[i]] = append(m[queue[i]], v.Encode())
		}
	}

	return m
}

func (s *server) getTargetFeatureG2() map[string][]string {
	return nil
}

func inheritanceAction(in []string, domain Domain, sortFn func([]*InheritanceEdge[string]), action func(Object, Object, Domain) error) error {
	var edges []*InheritanceEdge[string]
	for _, v := range in {
		edge := &InheritanceEdge[string]{}
		_ = edge.Decode(v)
		edges = append(edges, edge)
	}

	sortFn(edges)
	for _, v := range edges {
		if err := action(&NamedObject{Name: v.V}, &NamedObject{Name: v.U}, domain); err != nil {
			return err
		}
	}
	return nil
}
