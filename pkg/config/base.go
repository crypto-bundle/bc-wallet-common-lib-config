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

// BaseConfig is config for application base entity like environment, application run mode and etc
type BaseConfig struct {
	// -------------------
	// Application configs
	// -------------------
	// Application environment.
	// allowed: local, dev, testing, staging, production
	Environment string `envconfig:"APP_ENV" default:"development"`
	// Debug mode
	Debug     bool   `envconfig:"APP_DEBUG" default:"false"`
	StageName string `envconfig:"APP_STAGE" default:"dev"`

	// ----------------------------
	// Calculated config parameters
	hostname        string
	applicationName string
	applicationPID  int

	// Dependencies
	ldFlagManagerSrv ldFlagManagerService
}

// Prepare variables to static configuration
func (c *BaseConfig) Prepare() error {
	host, err := os.Hostname()
	if err != nil {
		return err
	}

	c.hostname = host
	c.applicationPID = os.Getpid()

	return nil
}

func (c *BaseConfig) PrepareWith(cfgSrvList ...interface{}) error {
	for _, cfgSrv := range cfgSrvList {
		switch castedCfg := cfgSrv.(type) {
		case ldFlagManagerService:
			c.ldFlagManagerSrv = castedCfg
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

// GetStageName is for getting log stage name environment
func (c *BaseConfig) GetStageName() string {
	return c.StageName
}

// GetApplicationPID is for getting application process identifier
func (c *BaseConfig) GetApplicationPID() int {
	return c.applicationPID
}

// GetApplicationName is for getting application name
func (c *BaseConfig) GetApplicationName() string {
	return c.applicationName
}

// SetApplicationName is for setting application name
func (c *BaseConfig) SetApplicationName(appName string) {
	c.applicationName = appName
}

func (c *BaseConfig) GetVersion() string {
	return c.ldFlagManagerSrv.GetVersion()
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
		applicationName:  applicationName,
		ldFlagManagerSrv: newDefaultLdFlagManager(),
	}
}
