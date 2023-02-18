package config

import "time"

type ldFlagManagerService interface {
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() uint64
	GetBuildDate() time.Time
}

type baseConfigService interface {
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
	GetApplicationPID() int
	GetVersion() string
	GetReleaseTag() string
	GetCommitID() string
	GetShortCommitID() string
	GetBuildNumber() uint64
	GetBuildDateTS() uint64
	GetBuildDate() time.Time
}

type targetConfigService interface {
	Prepare() error
	PrepareWith(baseConfigSrv baseConfigService) error
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
