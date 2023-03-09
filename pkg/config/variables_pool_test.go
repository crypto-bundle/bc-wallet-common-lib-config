package config

import (
	"fmt"
	"os"
	"strconv"
	"testing"
)

func TestVarPoolBaseEnvVariables(t *testing.T) {
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

	isDebug, _ := strconv.ParseBool(InitialEnvVariables["APP_DEBUG"])

	expectedResult := &BaseConfig{
		Environment: InitialEnvVariables["APP_ENV"],
		Debug:       isDebug,
		StageName:   InitialEnvVariables["APP_STAGE"],
	}

	baseCfg := &BaseConfig{}
	cfgVarPool := newConfigVarsPool(nil, baseCfg)
	err := cfgVarPool.Process()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if baseCfg.Debug != expectedResult.Debug {
		t.Errorf("not equal Debug")
	}

	if baseCfg.StageName != expectedResult.StageName {
		t.Errorf("not equal StageName")
	}

	if baseCfg.Environment != expectedResult.Environment {
		t.Errorf("not equal Environment")
	}

	t.Log("success")
}

func TestVarPoolSecretVariables(t *testing.T) {
	const initialDbPort uint16 = 12345
	var InitialEnvVariables = map[string]string{
		"DATABASE_DRIVER":                    "postgresql",
		"DATABASE_PORT":                      fmt.Sprintf("%d", initialDbPort),
		"TEST_FIELD_FOR_OVERWRITE_BY_SECRET": "initial_ENV_value",
	}

	var InitialSecretVariables = map[string]string{
		"DATABASE_USER":                      "secret_user",
		"DATABASE_PASSWORD":                  "secret_password",
		"TEST_FIELD_FOR_OVERWRITE_BY_SECRET": "initial_SECRET_value",
	}

	var MockSecretService = &mockSecretManager{
		ValuesPool: InitialSecretVariables,
	}
	for key, value := range InitialEnvVariables {
		err := os.Setenv(key, value)
		if err != nil {
			return
		}
	}

	type DbConfig struct {
		DatabaseDriver              string `envconfig:"DATABASE_DRIVER" required:"true"`
		DatabasePort                uint16 `envconfig:"DATABASE_PORT" default:"54321"`
		DatabaseUser                string `envconfig:"DATABASE_USER" secret:"true"`
		DatabasePassword            string `envconfig:"DATABASE_PASSWORD" secret:"true"`
		TestFieldForSecretOverwrite string `envconfig:"TEST_FIELD_FOR_OVERWRITE_BY_SECRET" secret:"true"`
	}

	testTypeStructSecrets := DbConfig{}
	expectedResult := &DbConfig{
		DatabaseUser:                "secret_user",
		DatabasePassword:            "secret_password",
		DatabaseDriver:              "postgresql",
		DatabasePort:                initialDbPort,
		TestFieldForSecretOverwrite: InitialSecretVariables["TEST_FIELD_FOR_OVERWRITE_BY_SECRET"],
	}

	cfgVarPool := newConfigVarsPool(MockSecretService, testTypeStructSecrets)
	err := cfgVarPool.Process()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if testTypeStructSecrets.DatabaseUser != expectedResult.DatabaseUser {
		t.Errorf("not equal DatabseUser")
	}

	if testTypeStructSecrets.DatabasePassword != expectedResult.DatabasePassword {
		t.Errorf("not equal DatabasePassword")
	}

	if testTypeStructSecrets.DatabaseDriver != expectedResult.DatabaseDriver {
		t.Errorf("not equal DatabaseDriver")
	}

	if testTypeStructSecrets.DatabasePort != expectedResult.DatabasePort {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStructSecrets.TestFieldForSecretOverwrite != expectedResult.TestFieldForSecretOverwrite {
		t.Errorf("not equal DatabasePort")
	}

	t.Log("success")
}
