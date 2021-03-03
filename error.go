package caskin

import "fmt"

var (
	// errors about entry
	ErrNil           = fmt.Errorf("nil data")
	ErrEmptyID       = fmt.Errorf("empty id")
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrNotExists     = fmt.Errorf("not exists")
	// errors about permission
	ErrNoReadPermission  = fmt.Errorf("no read permission")
	ErrNoWritePermission = fmt.Errorf("no write permission")
	// errors about superadmin
	ErrIsNotSuperAdmin       = fmt.Errorf("is not superadmin")
	ErrSuperAdminIsNoEnabled = fmt.Errorf("superadmin is not enabled")
	// errors about caskin initialization
	ErrInitializationNilDomainCreator = fmt.Errorf("domain creator is nil")
	ErrInitializationNilEnforcer      = fmt.Errorf("enforcer is nil")
	ErrInitializationNilEntryFactory  = fmt.Errorf("entry factory is nil")
	ErrInitializationNilMetaDB        = fmt.Errorf("metadata database is nil")
	// errors about db
	ErrCannotRecover = fmt.Errorf("record already exists, can not be recovered")
)
