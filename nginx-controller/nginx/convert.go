package nginx

import (
	"errors"
	"fmt"
	"strconv"

	"k8s.io/kubernetes/pkg/api/meta"
	"k8s.io/kubernetes/pkg/runtime"
)

var (
	// ErrorKeyNotFound is returned if the key was not found in the map
	ErrorKeyNotFound = errors.New("Key not found in map")
)

// There seems to be no composite interface in the kubernetes api package,
// so we have to declare our own.
type apiObject interface {
	meta.Object
	runtime.Object
}

// GetMapKeyAsBool searches the map for the given key and parses the key as bool
func GetMapKeyAsBool(m map[string]string, key string, context apiObject) (bool, error) {
	if str, exists := m[key]; exists {
		b, err := strconv.ParseBool(str)
		if err != nil {
			return false, fmt.Errorf("%s %v/%v '%s' contains invalid bool: %v, ignoring", context.GetObjectKind().GroupVersionKind().Kind, context.GetNamespace(), context.GetName(), key, err)
		}
		return b, nil
	}
	return false, ErrorKeyNotFound
}

// GetMapKeyAsInt tries to find and parse a key in a map as int64
func GetMapKeyAsInt(m map[string]string, key string, context apiObject) (int64, error) {
	if str, exists := m[key]; exists {
		i, err := strconv.ParseInt(str, 10, 64)
		if err != nil {
			return 0, fmt.Errorf("%s %v/%v '%s' contains invalid integer: %v, ignoring", context.GetObjectKind().GroupVersionKind().Kind, context.GetNamespace(), context.GetName(), key, err)
		}
		return i, nil
	}
	return 0, ErrorKeyNotFound
}
