package caskin

func (e *executor) GetBindExecutor(db MetaDBBindObjectAPI) *BindExecutor {
	return &BindExecutor{
		e: e,
		db: db,
	}
}
