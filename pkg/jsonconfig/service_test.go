package jsonconfig

import (
	"context"
	"os"
	"strconv"
	"testing"
)

func TestSimpleJSONStructWithSecret(t *testing.T) {
	ctx := context.Background()

	var InitialSecretVariables = map[string]string{
		"DATABASE_USER":     "secret_user_true",
		"DATABASE_PASSWORD": "secret_password_true",
		"DATABASE_NAME":     "test_database_true",
		"DATABASE_PORT":     "1234",
	}

	expectedPortNumber, err := strconv.Atoi(InitialSecretVariables["DATABASE_PORT"])
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var MockSecretSrv = &mockSecretManager{
		ValuesPool: InitialSecretVariables,
	}

	rawData, err := os.ReadFile("./simple.json")
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	unmarshaledData := &SimpleJSONCase{}

	cfgPreparer := &Service{}
	err = cfgPreparer.PrepareTo(unmarshaledData).PrepareFrom(rawData).
		With(MockSecretSrv).
		Do(ctx)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if unmarshaledData.IntFieldOne != 1 {
		t.Errorf("IntFieldOne not equal")
	}
	if unmarshaledData.IntFieldTwo != 2 {
		t.Errorf("IntFieldTwo not equal")
	}
	if unmarshaledData.IntFieldThree != 3 {
		t.Errorf("IntFieldThree not equal")
	}

	if unmarshaledData.StringField != "string_value" {
		t.Errorf("StringField not equal")
	}

	if unmarshaledData.FloatField != 4.567 {
		t.Errorf("FloatField not equal")
	}

	if unmarshaledData.DBUser != InitialSecretVariables["DATABASE_USER"] {
		t.Errorf("DBUser not equal")
	}

	if unmarshaledData.DBPassword != InitialSecretVariables["DATABASE_PASSWORD"] {
		t.Errorf("DBPassword not equal")
	}

	if unmarshaledData.DBName != InitialSecretVariables["DATABASE_NAME"] {
		t.Errorf("DBName not equal")
	}

	if unmarshaledData.DBPort != uint32(expectedPortNumber) {
		t.Errorf("DBPort not equal")
	}
}
