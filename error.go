package caskin

import "fmt"

// Sentinel errors returned by caskin operations.
//
// Callers should use errors.Is to check for specific errors:
//
//	if errors.Is(err, caskin.ErrNoReadPermission) { ... }
var (
	// ErrNil is returned when a required argument is nil.
	ErrNil = fmt.Errorf("nil data")
	// ErrEmptyID is returned when an entity has an unset (zero) ID.
	ErrEmptyID = fmt.Errorf("empty id")
	// ErrAlreadyExists is returned when trying to create an entity that
	// already exists (non-deleted) in the database.
	ErrAlreadyExists = fmt.Errorf("already exists")
	// ErrNotExists is returned when the target entity cannot be found.
	ErrNotExists = fmt.Errorf("not exists")
	// ErrInValidObject is returned when the object reference is invalid or
	// does not belong to the current domain.
	ErrInValidObject = fmt.Errorf("invalid object")
	// ErrInValidObjectType is returned when the object type string is not
	// registered in the current factory.
	ErrInValidObjectType = fmt.Errorf("invalid object type")
	// ErrCantChangeObjectType is returned when an update attempts to change
	// an existing object's type, which is not allowed.
	ErrCantChangeObjectType = fmt.Errorf("can't change object type")
	// ErrCantOperateRootObject is returned when an operation would affect the
	// root object of a domain, which is protected.
	ErrCantOperateRootObject = fmt.Errorf("can't operate root object")
	// ErrParentCanNotBeItself is returned when an object's parent ID is set
	// to its own ID, which would create a self-loop.
	ErrParentCanNotBeItself = fmt.Errorf("parent id can't be it self id")
	// ErrParentToDescendant is returned when moving an object would make its
	// new parent one of its own descendants, creating a cycle.
	ErrParentToDescendant = fmt.Errorf("can't change parent to descendant")
	// ErrInValidRequest is returned when the [DirectoryRequest] parameters
	// are inconsistent or missing required fields.
	ErrInValidRequest = fmt.Errorf("invalid request")

	// ErrNoReadPermission is returned when the caller lacks read access to
	// the target object/domain.
	ErrNoReadPermission = fmt.Errorf("no read permission")
	// ErrNoWritePermission is returned when the caller lacks write access to
	// the target object/domain.
	ErrNoWritePermission = fmt.Errorf("no write permission")
	// ErrNoManagePermission is returned when the caller lacks manage access
	// to the target object/domain.
	ErrNoManagePermission = fmt.Errorf("no manage permission")
	// ErrNoBackendPermission is returned when [IFeatureService.AuthBackend]
	// determines the caller is not authorised for the given backend endpoint.
	ErrNoBackendPermission = fmt.Errorf("no backend api permission")

	// ErrIsNotSuperadmin is returned when a superadmin-only operation is
	// attempted by a non-superadmin user, or when decoding a superadmin
	// token fails.
	ErrIsNotSuperadmin = fmt.Errorf("is not superadmin")
	// ErrInValidCurrent is returned when [ICurrentService] methods are called
	// before [ICurrentService.SetCurrent] has been called.
	ErrInValidCurrent = fmt.Errorf("invalid current api")
)
