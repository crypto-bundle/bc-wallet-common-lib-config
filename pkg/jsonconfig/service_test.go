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

	rawData, err := os.ReadFile("./service_single_object_test_data.json")
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

	if unmarshaledData.GetPort() != uint32(expectedPortNumber) {
		t.Errorf("GetPort not equal")
	}
}

func TestArrayJSONStructWithSecret(t *testing.T) {
	ctx := context.Background()

	var InitialSecretVariables = map[string]string{
		"DATABASE_USER_ONE":     "first_secret_user_true",
		"DATABASE_PASSWORD_ONE": "first_secret_password_true",
		"DATABASE_NAME_ONE":     "first_test_database_true",
		"DATABASE_PORT_ONE":     "1234",

		"DATABASE_USER_TWO":     "second_secret_user_true",
		"DATABASE_PASSWORD_TWO": "second_secret_password_true",
		"DATABASE_NAME_TWO":     "second_test_database_true",
		"DATABASE_PORT_TWO":     "5678",
	}

	expectedPortOneNumber, err := strconv.Atoi(InitialSecretVariables["DATABASE_PORT_ONE"])
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	expectedPortTwoNumber, err := strconv.Atoi(InitialSecretVariables["DATABASE_PORT_TWO"])
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	var MockSecretSrv = &mockSecretManager{
		ValuesPool: InitialSecretVariables,
	}

	rawData, err := os.ReadFile("./service_array_test_data.json")
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	unmarshaledData := &MixedJSONCase{}

	cfgPreparer := &Service{}
	err = cfgPreparer.PrepareTo(unmarshaledData).PrepareFrom(rawData).
		With(MockSecretSrv).
		Do(ctx)
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	if unmarshaledData.TopLevelField != 100500 {
		t.Errorf("DBPort not equal")
	}

	if len(unmarshaledData.List) != 2 {
		t.Errorf("wrong count in list of struct")
	}

	if unmarshaledData.List[0].IntFieldOne != 1 {
		t.Errorf("IntFieldOne not equal")
	}
	if unmarshaledData.List[0].IntFieldTwo != 2 {
		t.Errorf("IntFieldTwo not equal")
	}
	if unmarshaledData.List[0].IntFieldThree != 3 {
		t.Errorf("IntFieldThree not equal")
	}

	if unmarshaledData.List[0].StringField != "string_value_one" {
		t.Errorf("StringField not equal")
	}

	if unmarshaledData.List[0].FloatField != 4.567 {
		t.Errorf("FloatField not equal")
	}

	if unmarshaledData.List[0].DBUser != InitialSecretVariables["DATABASE_USER_ONE"] {
		t.Errorf("DBUser not equal")
	}

	if unmarshaledData.List[0].DBPassword != InitialSecretVariables["DATABASE_PASSWORD_ONE"] {
		t.Errorf("DBPassword not equal")
	}

	if unmarshaledData.List[0].DBName != InitialSecretVariables["DATABASE_NAME_ONE"] {
		t.Errorf("DBName not equal")
	}

	if unmarshaledData.List[0].DBPort != InitialSecretVariables["DATABASE_PORT_ONE"] {
		t.Errorf("DBPort not equal")
	}

	if unmarshaledData.List[0].GetPort() != uint32(expectedPortOneNumber) {
		t.Errorf("GetPort not equal")
	}

	if unmarshaledData.List[1].GetPort() != uint32(expectedPortTwoNumber) {
		t.Errorf("GetPort not equal")
	}
}
