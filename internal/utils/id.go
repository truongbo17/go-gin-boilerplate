package utils

import (
	"errors"
	"strconv"
)

var ErrInvalidID = errors.New("invalid positive 32-bit ID")

// ParseID parses an unsigned database ID independently of an HTTP transport.
func ParseID(value string) (uint, error) {
	id, err := strconv.ParseUint(value, 10, 32)
	if err != nil || id == 0 {
		return 0, ErrInvalidID
	}
	return uint(id), nil
}
