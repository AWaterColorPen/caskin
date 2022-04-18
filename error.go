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

	ErrNoReadPermission    = fmt.Errorf("no read permission")
	ErrNoWritePermission   = fmt.Errorf("no write permission")
	ErrNoManagePermission  = fmt.Errorf("no manage permission")
	ErrNoBackendPermission = fmt.Errorf("no backend api permission")

	ErrIsNotSuperAdmin = fmt.Errorf("is not superadmin")

	ErrInitializationNilEnforcer = fmt.Errorf("enforcer is nil")

	ErrInputPairArrayNotBelongSameUser    = fmt.Errorf("input user role pair array are not belong to same user")
	ErrInputPairArrayNotBelongSameRole    = fmt.Errorf("input user role pair array are not belong to same role")
	ErrInputPolicyListNotBelongSameRole   = fmt.Errorf("input policy list are not belong to same role")
	ErrInputPolicyListNotBelongSameObject = fmt.Errorf("input policy list are not belong to same object")
)
