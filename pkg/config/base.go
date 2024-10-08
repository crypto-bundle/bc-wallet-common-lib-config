/*
 *
 *
 * MIT NON-AI License
 *
 * Copyright (c) 2022-2024 Aleksei Kotelnikov(gudron2s@gmail.com)
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of the software and associated documentation files (the "Software"),
 * to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense,
 * and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions.
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * In addition, the following restrictions apply:
 *
 * 1. The Software and any modifications made to it may not be used for the purpose of training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining. This condition applies to any derivatives,
 * modifications, or updates based on the Software code. Any usage of the Software in an AI-training dataset is considered a breach of this License.
 *
 * 2. The Software may not be included in any dataset used for training or improving machine learning algorithms,
 * including but not limited to artificial intelligence, natural language processing, or data mining.
 *
 * 3. Any person or organization found to be in violation of these restrictions will be subject to legal action and may be held liable
 * for any damages resulting from such use.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM,
 * DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE
 * OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */

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
