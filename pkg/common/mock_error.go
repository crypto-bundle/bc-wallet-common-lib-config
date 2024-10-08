package common

import "errors"

var ErrMockFormatter = errors.New("mock_err_formatter")

type errFmt struct {
}

func (f *errFmt) ErrorWithCode(_ error, _ int) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrWithCode(_ error, _ int) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrorGetCode(_ error) int {
	return -1
}

func (f *errFmt) ErrGetCode(_ error) int {
	return -1
}

func (f *errFmt) ErrorNoWrap(_ error) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrNoWrap(_ error) error {
	return ErrMockFormatter
}

func (f *errFmt) ErrorOnly(_ error, _ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) Error(_ error, _ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) Errorf(_ error, _ string, _ ...interface{}) error {
	return ErrMockFormatter
}

func (f *errFmt) NewError(_ ...string) error {
	return ErrMockFormatter
}

func (f *errFmt) NewErrorf(_ string, _ ...interface{}) error {
	return ErrMockFormatter
}

func NewMockErrFormatter() *errFmt {
	return &errFmt{}
}
