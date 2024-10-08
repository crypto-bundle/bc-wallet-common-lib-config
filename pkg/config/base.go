package config

import (
	"os"
	"time"
)

const (
	EnvLocal      = "local"
	EnvDev        = "development"
	EnvStaging    = "staging"
	EnvTesting    = "testing"
	EnvProduction = "production"
)

var (
	_ baseConfigService = (*BaseConfig)(nil)
)

// BaseConfig is config for application base entity like environment, application run mode and etc...
type BaseConfig struct {
	ldFlagManagerSrv ldFlagManagerService
	e                errorFormatterService

	Environment      string `envconfig:"APP_ENV" default:"development"`
	StageName        string `envconfig:"APP_STAGE" default:"dev"`
	LocalEnvFilePath string `envconfig:"APP_LOCAL_ENV_FILE_PATH" default:"./env"`
	hostname         string
	applicationName  string
	applicationPID   int
	Debug            bool `envconfig:"APP_DEBUG" default:"false"`
}

// Prepare variables to static configuration...
func (c *BaseConfig) Prepare() error {
	host, err := os.Hostname()
	if err != nil {
		return c.e.ErrorOnly(err)
	}

	c.hostname = host
	c.applicationPID = os.Getpid()

	return nil
}

func (c *BaseConfig) PrepareWith(cfgSrvList ...interface{}) error {
	for _, cfgSrv := range cfgSrvList {
		switch castedCfgDep := cfgSrv.(type) {
		case ldFlagManagerService:
			c.ldFlagManagerSrv = castedCfgDep
		case errorFormatterService:
			c.e = castedCfgDep
		default:
			continue
		}
	}

	return nil
}

// GetHostName ...
func (c *BaseConfig) GetHostName() string {
	return c.hostname
}

// GetEnvironmentName ...
func (c *BaseConfig) GetEnvironmentName() string {
	return c.Environment
}

// IsProd ...
func (c *BaseConfig) IsProd() bool {
	return c.Environment == EnvProduction
}

// IsStage ...
func (c *BaseConfig) IsStage() bool {
	return c.Environment == EnvStaging
}

// IsTest ...
func (c *BaseConfig) IsTest() bool {
	return c.Environment == EnvStaging || c.Environment == EnvTesting
}

// IsDev ...
func (c *BaseConfig) IsDev() bool {
	return c.Environment == EnvLocal || c.Environment == EnvDev
}

// IsDebug ...
func (c *BaseConfig) IsDebug() bool {
	return c.Debug
}

// IsLocal ...
func (c *BaseConfig) IsLocal() bool {
	return c.Environment == EnvLocal
}

func (c *BaseConfig) GetLocalEnvFilePath() string {
	return c.LocalEnvFilePath
}

// GetStageName is for getting log stage name environment...
func (c *BaseConfig) GetStageName() string {
	return c.StageName
}

// GetApplicationPID is for getting application process identifier...
func (c *BaseConfig) GetApplicationPID() int {
	return c.applicationPID
}

// GetApplicationName is for getting application name...
func (c *BaseConfig) GetApplicationName() string {
	return c.applicationName
}

// SetApplicationName is for setting application name...
func (c *BaseConfig) SetApplicationName(appName string) {
	c.applicationName = appName
}

func (c *BaseConfig) GetReleaseTag() string {
	return c.ldFlagManagerSrv.GetReleaseTag()
}

func (c *BaseConfig) GetCommitID() string {
	return c.ldFlagManagerSrv.GetCommitID()
}

func (c *BaseConfig) GetShortCommitID() string {
	return c.ldFlagManagerSrv.GetShortCommitID()
}

func (c *BaseConfig) GetBuildNumber() uint64 {
	return c.ldFlagManagerSrv.GetBuildNumber()
}

func (c *BaseConfig) GetBuildDateTS() int64 {
	return c.ldFlagManagerSrv.GetBuildDateTS()
}

func (c *BaseConfig) GetBuildDate() time.Time {
	return c.ldFlagManagerSrv.GetBuildDate()
}

func NewBaseConfig(applicationName string) *BaseConfig {
	return &BaseConfig{
		Environment:      "",    // will be filled in config filling stage by config service
		Debug:            false, // will be filled in config filling stage by config service
		StageName:        "",    // will be filled in config filling stage by config service
		LocalEnvFilePath: "",    // will be filled in config filling stage by config service
		hostname:         "",    // will be filled in config filling stage by config service on call Prepare function
		applicationName:  applicationName,
		applicationPID:   0, // will be filled in config filling stage by config service on call Prepare function
		ldFlagManagerSrv: newDefaultLdFlagManager(),
		e:                nil, // will be filled in config filling stage by config service on call PrepareWith function
	}
}
