package data 

import (
	"strconv"
	"fmt"
	"errors"
	"strings"
)

type Runtime int32

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r Runtime) MarshalJSON() ([]byte, error) {

	jsonValue := fmt.Sprintf("%d min", r)
	quotedJSONValue := strconv.Quote(jsonValue)
	return []byte(quotedJSONValue), nil
}

func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "min"{
		return ErrInvalidRuntimeFormat
	}
	res, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return ErrInvalidRuntimeFormat
	}

	*r = Runtime(res)
	return nil
}