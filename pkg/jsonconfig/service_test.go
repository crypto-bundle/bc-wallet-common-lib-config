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

package jsonconfig

import (
	"context"
	"os"
	"strconv"
	"testing"

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

	var MockSecretDataSvc = &mockSecretManager{
		ValuesPool: InitialSecretVariables,
	}

	var MockErrorFormatterSvc = common.NewMockErrFormatter()

	rawData, err := os.ReadFile("./service_single_object_test_data.json")
	if err != nil {
		t.Errorf("%s", err)
		return
	}

	unmarshaledData := &SimpleJSONCase{}

	cfgPreparer := &Service{}
	err = cfgPreparer.PrepareTo(unmarshaledData).PrepareFrom(rawData).
		With(MockSecretDataSvc, MockErrorFormatterSvc).
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
