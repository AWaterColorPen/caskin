package web_feature

func (e *Executor) Dump() (*Dump, error) {
    feature, frontend, backend, err := e.get3pair()
    if err != nil {
        return nil, err
    }
    all := e.allWebFeatureRelation(e.operationDomain)
	return NewDump(feature, frontend, backend, all), nil
}
