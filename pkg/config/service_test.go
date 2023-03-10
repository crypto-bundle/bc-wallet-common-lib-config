package config

import (
	"context"
	"os"
	"testing"
)

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

	cfgManagerSrv := NewConfigManager()

	err := cfgManagerSrv.PrepareTo(baseCfg).Do(context.Background())
	if err != nil {
		t.Errorf("%s", err)
		return
	}

}
