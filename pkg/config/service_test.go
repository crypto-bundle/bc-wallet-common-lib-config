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
