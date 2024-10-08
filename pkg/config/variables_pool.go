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
	"errors"
	"os"
	"reflect"
	"strconv"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

var (
	ErrPassedStructMustBeAPointer       = errors.New("must be a pointer")
	ErrPassedStructMustBeAStructPointer = errors.New("must be a struct pointer")
	ErrVariableEmptyButRequired         = errors.New("variables is empty and has required tag")
)

var _ configVariablesPoolService = (*configVariablesPool)(nil)

type configVariablesPool struct {
	e                     errorFormatterService
	targetConfigSvc       interface{}
	secretsDataSvc        secretManagerService
	dependenciesSvc       []interface{}
	envVariablesNameList  []string
	envVariablesList      []common.Field
	secretVariablesList   []common.Field
	envVariablesNameCount uint16
	secretVariablesCount  uint16
}

func (u *configVariablesPool) addSecretVariable(variable common.Field) error {
	err := common.SetField(variable.Value, variable.RfValue)
	if err != nil {
		return u.e.ErrorNoWrap(err)
	}

	u.secretVariablesCount++
	u.secretVariablesList = append(u.secretVariablesList, variable)

	return nil
}

func (u *configVariablesPool) addEnvVariable(variable common.Field) error {
	err := common.SetField(variable.Value, variable.RfValue)
	if err != nil {
		return u.e.ErrorNoWrap(err)
	}

	u.envVariablesNameCount++
	u.envVariablesList = append(u.envVariablesList, variable)

	return nil
}

func (u *configVariablesPool) Process() error {
	err := u.processFields(u.targetConfigSvc)
	if err != nil {
		return u.e.ErrorNoWrap(err)
	}

	return nil
}

// extractFields returns information of the struct fields, including nested structures
// based on https://github.com/kelseyhightower/envconfig
// TODO: refactor it - separate by sub-function
//
//nolint:funlen,gocognit,gocyclo,cyclop // it's ok. Need to refactor this function, but now - it's ok.
func (u *configVariablesPool) processFields(target interface{}) error {
	targetSource := reflect.ValueOf(target)

	// must be a pointer
	if targetSource.Kind() != reflect.Ptr {
		return u.e.ErrorOnly(ErrPassedStructMustBeAPointer)
	}

	// pointer must refer to structure
	element := targetSource.Elem()
	elemType := element.Type()

	castedInitConfigField, isPossibleToCast := element.Addr().Interface().(configInitService)
	if isPossibleToCast {
		prepErr := castedInitConfigField.InitWith(u.dependenciesSvc...)
		if prepErr != nil {
			return u.e.ErrorOnly(prepErr)
		}
	}

	// iterate over struct fields
	numFields := elemType.NumField()
	for i := range numFields {
		structFieldInfo := elemType.Field(i) // struct field info

		fieldValue := element.Field(i) // reflect.RfValue
		if !fieldValue.CanSet() {
			continue
		}

		isIgnored, _ := strconv.ParseBool(structFieldInfo.Tag.Get(common.TagIgnored))
		if isIgnored {
			continue
		}

		// unfold pointers
		for fieldValue.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				if fieldValue.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
			}

			fieldValue = fieldValue.Elem()
		}

		// recursively process nested struct
		if fieldValue.Kind() == reflect.Struct && fieldValue.CanInterface() {
			processErr := u.processFields(fieldValue.Addr().Interface())
			if processErr != nil {
				return u.e.ErrorOnly(processErr)
			}

			continue
		}

		var isSecret = false

		boolVarSrt, isTagExists := structFieldInfo.Tag.Lookup(common.TagSecret)
		if isTagExists {
			boolVar, err := strconv.ParseBool(boolVarSrt)
			if err != nil {
				return u.e.ErrorOnly(err)
			}

			isSecret = boolVar
		}

		var isRequired = false

		boolVarSrt, isTagExists = structFieldInfo.Tag.Lookup(common.TagRequired)
		if isTagExists {
			boolVar, err := strconv.ParseBool(boolVarSrt)
			if err != nil {
				return u.e.ErrorOnly(err)
			}

			isRequired = boolVar
		}

		if isSecret {
			envConfigKey := structFieldInfo.Tag.Get(common.TagEnvconfig)

			value, isExists := u.secretsDataSvc.GetByName(envConfigKey)
			if !isExists && isRequired {
				return u.e.ErrorOnly(ErrVariableEmptyButRequired, structFieldInfo.Name)
			}

			commonField := common.Field{
				Name:    structFieldInfo.Name,
				RfValue: fieldValue,
				RfTags:  structFieldInfo.Tag,
				Value:   value,
			}

			addErr := u.addSecretVariable(commonField)
			if addErr != nil {
				return u.e.ErrorOnly(addErr)
			}

			continue
		}

		envConfigKey := structFieldInfo.Tag.Get(common.TagEnvconfig)

		value, isEnvVariableExists := os.LookupEnv(envConfigKey)
		if !isEnvVariableExists && isRequired {
			return u.e.ErrorOnly(ErrVariableEmptyButRequired, structFieldInfo.Name)
		}

		defaultValue, hasDefaultValue := structFieldInfo.Tag.Lookup(common.TagDefault)
		if !isEnvVariableExists && hasDefaultValue {
			value = defaultValue
		}

		commonField := common.Field{
			Name:    structFieldInfo.Name,
			RfValue: fieldValue,
			RfTags:  structFieldInfo.Tag,
			Value:   value,
		}

		addErr := u.addEnvVariable(commonField)
		if addErr != nil {
			return u.e.ErrorOnly(addErr)
		}
	}

	castedField, isPossibleToCast := element.Addr().Interface().(dependentConfigService)
	if isPossibleToCast {
		if u.dependenciesSvc != nil {
			prepErr := castedField.PrepareWith(u.dependenciesSvc...)
			if prepErr != nil {
				return u.e.ErrorOnly(prepErr)
			}
		}

		prepErr := castedField.Prepare()
		if prepErr != nil {
			return u.e.ErrorOnly(prepErr)
		}

		return nil
	}

	castedConfigField, isPossibleToCast := element.Addr().Interface().(configService)
	if isPossibleToCast {
		prepErr := castedConfigField.Prepare()
		if prepErr != nil {
			return u.e.ErrorOnly(prepErr)
		}
	}

	return nil
}

func (u *configVariablesPool) ClearENV() error {
	for i := uint16(0); i != u.envVariablesNameCount; i++ {
		envField := u.envVariablesList[i]

		err := os.Unsetenv(envField.Name)
		if err != nil {
			return u.e.ErrorOnly(err)
		}
	}

	return nil
}

func newConfigVarsPool(errFmtSvc errorFormatterService,
	secretDataProviderSvc secretManagerService,
	processedConfig interface{},
	dependenciesSvcList []interface{},
) *configVariablesPool {
	return &configVariablesPool{
		e: errFmtSvc,

		dependenciesSvc: dependenciesSvcList,
		targetConfigSvc: processedConfig,
		secretsDataSvc:  secretDataProviderSvc,

		envVariablesNameCount: 0,
		envVariablesNameList:  make([]string, 0),
		envVariablesList:      make([]common.Field, 0),

		secretVariablesCount: 0,
		secretVariablesList:  make([]common.Field, 0),
	}
}
