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
	"context"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

type mockSecretManager struct {
	ValuesPool map[string]string
}

func (m *mockSecretManager) GetByName(keyName string) (string, bool) {
	result, isExists := m.ValuesPool[keyName]

	return result, isExists
}

func (m *mockSecretManager) GetByNameAndPath(keyName string) (string, bool) {
	result, isExists := m.ValuesPool[keyName]

	return result, isExists
}

func newMockLdFlagManager(releaseTag string,
	commitID string,
	shortCommitID string,
	buildNumber string,
) *ldFlagManager {
	buildTime := time.Now()

	buildNumberRaw, err := strconv.ParseUint(buildNumber, 10, 0)
	if err != nil {
		buildNumberRaw = 0
	}

	return &ldFlagManager{
		buildDateAt:   buildTime,
		buildDateTS:   uint64(buildTime.Unix()),
		releaseTag:    releaseTag,
		commitID:      commitID,
		shortCommitID: shortCommitID,
		buildNumber:   buildNumberRaw,
	}
}

func TestBaseEnvVariables(t *testing.T) {
	var InitialEnvVariables = map[string]string{
		"APP_ENV":          "development",
		"APP_DEBUG":        "false",
		"APP_LOGGER_LEVEL": "debug",
		"APP_STAGE":        "dev",
	}
	for key, value := range InitialEnvVariables {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}

	baseCfg := &BaseConfig{}

	cfgManagerSrv := NewConfigManager(common.NewMockErrFormatter())

	err := cfgManagerSrv.PrepareTo(baseCfg).Do(context.Background())
	if err != nil {
		t.Errorf("%s", err)
		return
	}
}

func TestBaseEnvVariablesPrepareWith(t *testing.T) {
	const (
		ldFlagMockVersion     = "v0.0.0"
		ldFlagMockReleaseTag  = "v0.0.0~mock-release"
		ldFlagMockCommit      = "0000000000000000mock00000000000000000000"
		ldFlagMockShortCommit = "00mock00"
		ldFlagMockBuildNumber = "0"
	)

	var InitialEnvVariables = map[string]string{
		"APP_ENV":   "development",
		"APP_DEBUG": "false",
		"APP_STAGE": "dev",
	}
	for key, value := range InitialEnvVariables {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}

	isDebug, _ := strconv.ParseBool(InitialEnvVariables["APP_DEBUG"])

	expectedResult := &BaseConfig{
		Environment: InitialEnvVariables["APP_ENV"],
		Debug:       isDebug,
		StageName:   InitialEnvVariables["APP_STAGE"],

		hostname:         "",
		applicationName:  "",
		applicationPID:   0,
		ldFlagManagerSrv: nil,
		e:                nil,
	}

	baseCfg := &BaseConfig{}
	mockLdFlagManager := newMockLdFlagManager(ldFlagMockReleaseTag,
		ldFlagMockCommit,
		ldFlagMockShortCommit,
		ldFlagMockBuildNumber)

	cfgManagerSrv := NewConfigManager(common.NewMockErrFormatter())
	err := cfgManagerSrv.PrepareTo(baseCfg).With(mockLdFlagManager).
		Do(context.Background())
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if baseCfg.GetEnvironmentName() != expectedResult.GetEnvironmentName() {
		t.Errorf("not equal EnvironmentName")
	}

	if (baseCfg.IsDebug() != expectedResult.IsDebug()) && baseCfg.IsDebug() {
		t.Errorf("not equal IsDebug")
	}

	if baseCfg.GetStageName() != expectedResult.GetStageName() {
		t.Errorf("not equal StageName")
	}
}
