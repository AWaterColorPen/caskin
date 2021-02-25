package caskin

import "fmt"

var (
	ErrEmptyID       = fmt.Errorf("empty id")
	ErrAlreadyExists = fmt.Errorf("already exists")
	ErrNotExists     = fmt.Errorf("not exists")

	ErrNoReadPermission  = fmt.Errorf("no read permission")
	ErrNoWritePermission = fmt.Errorf("no write permission")

	ErrIsNotSuperAdmin = fmt.Errorf("is no superadmin")
)
