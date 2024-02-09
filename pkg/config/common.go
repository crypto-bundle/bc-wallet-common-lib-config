package config

import (
	"time"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

type ldFlagManagerService interface {
	GetVersion() string
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
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() uint64
	GetBuildDate() time.Time
}

type configVariablesPoolService interface {
	addSecretVariable(common.Field) error
	addEnvVariable(common.Field) error
}

type secretManagerService interface {
	GetByName(keyName string) (string, bool)
}

type secretAccessorService interface {
	GetCredentialsBytes() (b []byte, err error)
	GetCredentialsBytesByPath(path string) (b []byte, err error)
	GetCredentialsByPathAndKey(path, field string) (string, error)
	GetCredentialsByPathAndKeys(path string, fields ...string) (map[string]string, error)
}
