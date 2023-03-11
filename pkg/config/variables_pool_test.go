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

	testTypeStructSecrets := &DbConfig{}
	expectedResult := &DbConfig{
		DatabaseUser:                InitialSecretVariables["DATABASE_USER"],
		DatabasePassword:            InitialSecretVariables["DATABASE_PASSWORD"],
		DatabaseDriver:              InitialEnvVariables["DATABASE_DRIVER"],
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
		t.Errorf("not equal TestFieldForSecretOverwrite")
	}
}

type TestDbConfigForPrepare struct {
	DatabaseDriver   string `envconfig:"DATABASE_DRIVER" required:"true"`
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"postgresql.local"`
	DatabasePort     uint16 `envconfig:"DATABASE_PORT" default:"54321"`
	DatabaseUser     string `envconfig:"DATABASE_USER" secret:"true"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" secret:"true"`
	DatabaseName     string `envconfig:"DATABASE_NAME" secret:"true"`

	// calculated fields
	dbDSN string
}

func (c *TestDbConfigForPrepare) Prepare() error {
	c.dbDSN = fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%t",
		c.DatabaseUser, c.DatabasePassword, c.DatabaseHost, c.DatabaseName, false)

	return nil
}

func (c *TestDbConfigForPrepare) PrepareWith(cfgSrvList ...interface{}) error {
	return nil
}

func TestVarPoolVariablesWithSecretAndPrepare(t *testing.T) {
	const initialDbPort uint16 = 12345
	var InitialEnvVariables = map[string]string{
		"DATABASE_DRIVER": "postgresql",
		"DATABASE_PORT":   fmt.Sprintf("%d", initialDbPort),
		"DATABASE_HOST":   "127.0.0.1",
	}

	var InitialSecretVariables = map[string]string{
		"DATABASE_USER":     "secret_user",
		"DATABASE_PASSWORD": "secret_password",
		"DATABASE_NAME":     "test_database",
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

	testTypeStruct := &TestDbConfigForPrepare{}
	expectedResult := &TestDbConfigForPrepare{
		DatabaseUser:     InitialSecretVariables["DATABASE_USER"],
		DatabasePassword: InitialSecretVariables["DATABASE_PASSWORD"],
		DatabaseHost:     InitialEnvVariables["DATABASE_HOST"],
		DatabaseDriver:   InitialEnvVariables["DATABASE_DRIVER"],
		DatabaseName:     InitialSecretVariables["DATABASE_NAME"],
		DatabasePort:     initialDbPort,
		dbDSN: fmt.Sprintf("postgresql://%s:%s@%s/%s?sslmode=%t",
			InitialSecretVariables["DATABASE_USER"], InitialSecretVariables["DATABASE_PASSWORD"],
			InitialEnvVariables["DATABASE_HOST"], InitialSecretVariables["DATABASE_NAME"],
			false),
	}

	cfgVarPool := newConfigVarsPool(MockSecretService, testTypeStruct)
	err := cfgVarPool.Process()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if testTypeStruct.DatabaseUser != expectedResult.DatabaseUser {
		t.Errorf("not equal DatabseUser")
	}

	if testTypeStruct.DatabasePassword != expectedResult.DatabasePassword {
		t.Errorf("not equal DatabasePassword")
	}

	if testTypeStruct.DatabaseDriver != expectedResult.DatabaseDriver {
		t.Errorf("not equal DatabaseDriver")
	}

	if testTypeStruct.DatabasePort != expectedResult.DatabasePort {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.DatabaseHost != expectedResult.DatabaseHost {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.DatabaseName != expectedResult.DatabaseName {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.dbDSN != expectedResult.dbDSN {
		t.Errorf("not equal dbDSN")
	}
}

type TestDbEmbeddedConfig struct {
	EmbeddedFieldOne string `envconfig:"EMBEDDED_FIELD_ONE" required:"true"`
	EmbeddedFieldTwo string `envconfig:"EMBEDDED_FIELD_TWO" required:"true"`
}

type TestDbEmbeddedConfigForPrepare struct {
	DatabaseDriver   string `envconfig:"DATABASE_DRIVER" required:"true"`
	DatabaseHost     string `envconfig:"DATABASE_HOST" default:"postgresql.local"`
	DatabasePort     uint16 `envconfig:"DATABASE_PORT" default:"54321"`
	DatabaseUser     string `envconfig:"DATABASE_USER" secret:"true"`
	DatabasePassword string `envconfig:"DATABASE_PASSWORD" secret:"true"`
	DatabaseName     string `envconfig:"DATABASE_NAME" secret:"true"`

	// embedded structs
	*TestDbEmbeddedConfig
}

func TestVarPoolVariablesWithEmbeddedStructsAndSecrets(t *testing.T) {
	const initialDbPort uint16 = 12345
	var InitialEnvVariables = map[string]string{
		"DATABASE_DRIVER": "postgresql",
		"DATABASE_PORT":   fmt.Sprintf("%d", initialDbPort),
		"DATABASE_HOST":   "127.0.0.1",

		"EMBEDDED_FIELD_ONE": "embedded_field_one_value",
		"EMBEDDED_FIELD_TWO": "embedded_field_two_value",
	}

	var InitialSecretVariables = map[string]string{
		"DATABASE_USER":     "secret_user",
		"DATABASE_PASSWORD": "secret_password",
		"DATABASE_NAME":     "test_database",
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

	testTypeStruct := &TestDbEmbeddedConfigForPrepare{}
	expectedResult := &TestDbEmbeddedConfigForPrepare{
		DatabaseUser:     InitialSecretVariables["DATABASE_USER"],
		DatabasePassword: InitialSecretVariables["DATABASE_PASSWORD"],
		DatabaseHost:     InitialEnvVariables["DATABASE_HOST"],
		DatabaseDriver:   InitialEnvVariables["DATABASE_DRIVER"],
		DatabaseName:     InitialSecretVariables["DATABASE_NAME"],
		DatabasePort:     initialDbPort,

		TestDbEmbeddedConfig: &TestDbEmbeddedConfig{
			EmbeddedFieldOne: InitialEnvVariables["EMBEDDED_FIELD_ONE"],
			EmbeddedFieldTwo: InitialEnvVariables["EMBEDDED_FIELD_TWO"],
		},
	}

	cfgVarPool := newConfigVarsPool(MockSecretService, testTypeStruct)
	err := cfgVarPool.Process()
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if testTypeStruct.DatabaseUser != expectedResult.DatabaseUser {
		t.Errorf("not equal DatabseUser")
	}

	if testTypeStruct.DatabasePassword != expectedResult.DatabasePassword {
		t.Errorf("not equal DatabasePassword")
	}

	if testTypeStruct.DatabaseDriver != expectedResult.DatabaseDriver {
		t.Errorf("not equal DatabaseDriver")
	}

	if testTypeStruct.DatabasePort != expectedResult.DatabasePort {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.DatabaseHost != expectedResult.DatabaseHost {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.DatabaseName != expectedResult.DatabaseName {
		t.Errorf("not equal DatabasePort")
	}

	if testTypeStruct.EmbeddedFieldOne != expectedResult.EmbeddedFieldOne {
		t.Errorf("not equal EmbeddedFieldOne")
	}

	if testTypeStruct.EmbeddedFieldTwo != expectedResult.EmbeddedFieldTwo {
		t.Errorf("not equal EmbeddedFieldOne")
	}
}
