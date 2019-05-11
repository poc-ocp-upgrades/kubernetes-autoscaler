package model

import (
	"fmt"
)

type KeyError struct{ key interface{} }

func NewKeyError(key interface{}) KeyError {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return KeyError{key}
}
func (e KeyError) Error() string {
	_logClusterCodePath()
	defer _logClusterCodePath()
	return fmt.Sprintf("KeyError: %s", e.key)
}
