package types

import (
	"errors"
)

var InputError = errors.New("input error")

func NewInputError(err error) error {
	return errors.Join(InputError, err)
}
