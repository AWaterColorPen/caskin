package caskin

import "fmt"

var (
	ErrEmptyID           = fmt.Errorf("empty id")
	ErrNoReadPermission  = fmt.Errorf("no read permission")
	ErrNoWritePermission = fmt.Errorf("no write permission")
)
