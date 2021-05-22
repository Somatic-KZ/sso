package database

import (
	"errors"
	"fmt"
	"strings"
)

type ErrInvalidDatastoreName []string

func (ds ErrInvalidDatastoreName) Error() error {
	return fmt.Errorf("datastore: invalid datastore name. Must be one of: %s", strings.Join(ds, ", "))
}

var (
	ErrEmptyConfigStruct       = errors.New("empty config structure")
	ErrDatastoreNotImplemented = errors.New("datastore not implemented")
)
