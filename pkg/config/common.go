package config

import "time"

type ldFlagManagerService interface {
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() int64
	GetBuildDate() time.Time
}

type configService interface {
	Prepare() error
	PrepareWith(cfgSrv ...interface{}) error
}

type baseConfigService interface {
	configService

	GetHostName() string
	GetEnvironmentName() string
	IsProd() bool
	IsStage() bool
	IsTest() bool
	IsDev() bool
	IsDebug() bool
	IsLocal() bool
	GetMinimalLogLevel() string
	GetStageName() string
	GetApplicationName() string
	SetApplicationName() string
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
	addSecretVariable(field) error
	addEnvVariable(field) error
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
