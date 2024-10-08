package config

import (
	"time"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

type ldFlagManagerService interface {
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}

type dependentConfigService interface {
	Prepare() error
	PrepareWith(cfgSrv ...interface{}) error
}

type configInitService interface {
	InitWith(cfgSrv ...interface{}) error
}

type configService interface {
	Prepare() error
}

type baseConfigService interface {
	dependentConfigService

	GetHostName() string
	GetEnvironmentName() string
	IsProd() bool
	IsStage() bool
	IsTest() bool
	IsDev() bool
	IsDebug() bool
	IsLocal() bool
	GetStageName() string
	GetApplicationName() string
	SetApplicationName(appName string)
	GetApplicationPID() int
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}

type configVariablesPoolService interface {
	addSecretVariable(field common.Field) error
	addEnvVariable(field common.Field) error
}

type secretManagerService interface {
	GetByName(keyName string) (string, bool)
}

//nolint:interfacebloat // it's ok here, we need it we must use it as one big interface
type errorFormatterService interface {
	ErrorWithCode(err error, code int) error
	ErrWithCode(err error, code int) error
	ErrorGetCode(err error) int
	ErrGetCode(err error) int
	// ErrorNoWrap function for pseudo-wrap error, must be used in case of linter warnings...
	ErrorNoWrap(err error) error
	// ErrNoWrap same with ErrorNoWrap function, just alias for ErrorNoWrap, just short function name...
	ErrNoWrap(err error) error
	ErrorOnly(err error, details ...string) error
	Error(err error, details ...string) error
	Errorf(err error, format string, args ...interface{}) error
	NewError(details ...string) error
	NewErrorf(format string, args ...interface{}) error
}
