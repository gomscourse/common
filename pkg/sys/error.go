package sys

import (
	"errors"
	"github.com/gomscourse/common/pkg/sys/codes"
)

type commonError struct {
	msg  string
	code codes.Code
}

func NewCommonError(msg string, code codes.Code) *commonError {
	return &commonError{msg: msg, code: code}
}

func (ce *commonError) Error() string {
	return ce.msg
}

func (ce *commonError) Code() codes.Code {
	return ce.code
}

func IsCommonError(err error) bool {
	var ce *commonError
	return errors.As(err, &ce)
}

func GetCommonError(err error) *commonError {
	var ce *commonError
	if !errors.As(err, &ce) {
		return nil
	}

	return ce
}
