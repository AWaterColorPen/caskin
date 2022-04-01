package caskin

import "fmt"

var (
	// errors about entry
	ErrNil                  = fmt.Errorf("nil data")
	ErrEmptyID              = fmt.Errorf("empty id")
	ErrAlreadyExists        = fmt.Errorf("already exists")
	ErrNotExists            = fmt.Errorf("not exists")
	ErrInValidObject        = fmt.Errorf("invalid object")
	ErrInValidObjectType    = fmt.Errorf("invalid object type")
	ErrCantChangeObjectType = fmt.Errorf("can't change object type")
	ErrInValidParentObject  = fmt.Errorf("invalid parent object")
	ErrParentCanNotBeItself = fmt.Errorf("parent id can't be it self id")
	ErrParentToDescendant   = fmt.Errorf("can't change parent to descendant")
	ErrInValidAction        = fmt.Errorf("invalid action")

	// errors about permission
	ErrNoReadPermission    = fmt.Errorf("no read permission")
	ErrNoWritePermission   = fmt.Errorf("no write permission")
	ErrNoManagePermission  = fmt.Errorf("no manage permission")
	ErrNoBackendPermission = fmt.Errorf("no backend api permission")

	// errors about superadmin
	ErrIsNotSuperAdmin = fmt.Errorf("is not superadmin")

	// errors about caskin initialization
	ErrInitializationNilEnforcer     = fmt.Errorf("enforcer is nil")
	ErrInitializationNilEntryFactory = fmt.Errorf("entry factory is nil")
	ErrInitializationNilMetaDB       = fmt.Errorf("metadata database is nil")

	// errors about current provider
	ErrProviderGet = fmt.Errorf("provider can't get current status")

	// errors about user role pair
	ErrInputPairArrayNotBelongSameUser = fmt.Errorf("input user role pair array are not belong to same user")
	ErrInputPairArrayNotBelongSameRole = fmt.Errorf("input user role pair array are not belong to same role")

	// errors about policy list
	ErrInputPolicyListNotBelongSameRole   = fmt.Errorf("input policy list are not belong to same role")
	ErrInputPolicyListNotBelongSameObject = fmt.Errorf("input policy list are not belong to same object")

	// errors about special object
	ErrObjectTypeObjectIDMustBeItselfID = fmt.Errorf("the object id of object type's object must be itself's id")
	ErrEmptyParentIdOrNotSuperadmin     = fmt.Errorf("the parent id is empty or you are operating root object without superadmin authority")
)
