package errors

import "sync"

//nolint:gochecknoglobals // it's ok
var errorsFmtService errorFormatterService = nil

func InitInternalFmt(fmtSvc errorFormatterService) {
	sync.OnceFunc(func() {
		if errorsFmtService == nil {
			errorsFmtService = fmtSvc
		}
	})
}

func ErrorWithCode(err error, code int) error {
	return errorsFmtService.ErrorWithCode(err, code)
}

func ErrWithCode(err error, code int) error {
	return errorsFmtService.ErrWithCode(err, code)
}

func ErrorGetCode(err error) int {
	return errorsFmtService.ErrorGetCode(err)
}

func ErrGetCode(err error) int {
	return errorsFmtService.ErrGetCode(err)
}

func ErrorNoWrap(err error) error {
	return errorsFmtService.ErrorNoWrap(err)
}

func ErrNoWrap(err error) error {
	return errorsFmtService.ErrNoWrap(err)
}

func ErrorOnly(err error, details ...string) error {
	return errorsFmtService.ErrorOnly(err, details...)
}

func Error(err error, details ...string) error {
	return errorsFmtService.Error(err, details...)
}

func Errorf(err error, format string, args ...interface{}) error {
	return errorsFmtService.Errorf(err, format, args...)
}

func NewError(details ...string) error {
	return errorsFmtService.NewError(details...)
}

func NewErrorf(format string, args ...interface{}) error {
	return errorsFmtService.NewErrorf(format, args...)
}
