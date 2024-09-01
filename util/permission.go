package util

// all supported permissions
const (
	READ   = "read"
	WRITE  = "write"
	DELETE = "delete"
)

func IsSupportedPermission(permission string) bool {
	switch permission {
	case READ, WRITE, DELETE:
		return true
	default:
		return false
	}
}
