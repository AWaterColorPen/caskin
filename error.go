package caskin

import "fmt"

var (
	ErrNil                   = fmt.Errorf("nil data")
	ErrEmptyID               = fmt.Errorf("empty id")
	ErrAlreadyExists         = fmt.Errorf("already exists")
	ErrNotExists             = fmt.Errorf("not exists")
	ErrInValidObject         = fmt.Errorf("invalid object")
	ErrInValidObjectType     = fmt.Errorf("invalid object type")
	ErrCantChangeObjectType  = fmt.Errorf("can't change object type")
	ErrCantOperateRootObject = fmt.Errorf("can't operate root object")
	ErrParentCanNotBeItself  = fmt.Errorf("parent id can't be it self id")
	ErrParentToDescendant    = fmt.Errorf("can't change parent to descendant")
	ErrInValidRequest        = fmt.Errorf("invalid request")

	ErrNoReadPermission    = fmt.Errorf("no read permission")
	ErrNoWritePermission   = fmt.Errorf("no write permission")
	ErrNoManagePermission  = fmt.Errorf("no manage permission")
	ErrNoBackendPermission = fmt.Errorf("no backend api permission")

	ErrIsNotSuperadmin           = fmt.Errorf("is not superadmin")
	ErrInitializationNilEnforcer = fmt.Errorf("enforcer is nil")
	ErrInValidCurrent            = fmt.Errorf("invalid current api")
)
