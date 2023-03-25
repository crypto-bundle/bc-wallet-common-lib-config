package jsonconfig

import "time"

type configService interface {
	Prepare() error
	PrepareWith(cfgSrv ...interface{}) error
}

type secretManagerService interface {
	GetByName(keyName string) (string, bool)
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
