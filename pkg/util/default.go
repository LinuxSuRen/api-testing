// Package util provides a set of common functions
package util

// MakeSureNotNil makes sure the parameter is not nil
func MakeSureNotNil[T any](inter T) T {
	switch val := any(inter).(type) {
	case func():
		if val == nil {
			val = func() {
				// only making sure this is not nil
			}
			return any(val).(T)
		}
	case map[string]string:
		if val == nil {
			val = map[string]string{}
			return any(val).(T)
		}
	}
	return inter
}

// ContentType is the HTTP header key
const (
	ContentType       = "Content-Type"
	MultiPartFormData = "multipart/form-data"
	Form              = "application/x-www-form-urlencoded"
)
