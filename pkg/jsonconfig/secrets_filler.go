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
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

var (
	ErrPassedStructMustBeAPointer       = errors.New("must be a pointer")
	ErrPassedStructMustBeAStructPointer = errors.New("must be a struct pointer")
	ErrVariableEmptyButRequired         = errors.New("variables is empty and has required tag")
	ErrWrongSecretStringFormat          = errors.New("wrong secret string format")
)

type secretFiller struct {
	e              errorFormatterService
	secretsDataSvc secretManagerService

	target          interface{}
	dependenciesSvc []interface{}
}

func (u *secretFiller) Process() error {
	return u.processFields(u.target)
}

// extractFields returns information of the struct fields, including nested structures
// based on https://github.com/kelseyhightower/envconfig
// TODO: refactor it - separate by sub-function
//
//nolint:funlen,gocognit,cyclop // it's ok. Need to refactor this function, but now - it's ok.
func (u *secretFiller) processFields(target interface{}) error {
	targetSource := reflect.ValueOf(target)

	// must be a pointer
	switch targetSource.Kind() {
	case reflect.Slice:
	case reflect.Ptr:
	default:
		return u.e.ErrorOnly(ErrPassedStructMustBeAPointer)
	}

	// pointer must refer to structure
	element := targetSource.Elem()
	elemType := element.Type()

	// iterate over struct fields
	numFields := elemType.NumField()
	for i := range numFields {
		structField := elemType.Field(i) // struct field info

		fieldValue := element.Field(i) // reflect.RfValue
		if !fieldValue.CanSet() {
			continue
		}

		// recursively process nested struct
		if fieldValue.Kind() == reflect.Struct && fieldValue.CanInterface() {
			processErr := u.processFields(fieldValue.Addr().Interface())
			if processErr != nil {
				return u.e.ErrorOnly(processErr)
			}

			continue
		}

		if fieldValue.Kind() == reflect.Slice && fieldValue.CanInterface() {
			for j := range fieldValue.Len() {
				item := fieldValue.Index(j)

				indirectValue := reflect.Indirect(item)
				if indirectValue.Kind() == reflect.Struct {
					processErr := u.processFields(indirectValue.Addr().Interface())
					if processErr != nil {
						return u.e.ErrorOnly(processErr)
					}

					continue
				}
			}
		}

		boolVarSrt, isTagExists := structField.Tag.Lookup(common.TagSecret)
		if !isTagExists {
			continue
		}

		boolVar, err := strconv.ParseBool(boolVarSrt)
		if err != nil {
			return u.e.ErrorOnly(err)
		}

		if !boolVar {
			continue
		}

		value := fieldValue.String()
		if !strings.HasPrefix(value, "!secret:") {
			continue
		}

		separated := strings.Split(value, ":")
		if len(separated) > 2 {
			return u.e.ErrorOnly(ErrWrongSecretStringFormat)
		}

		secretKey := separated[1]

		value, isExists := u.secretsDataSvc.GetByName(secretKey)
		if !isExists {
			return u.e.ErrorOnly(ErrVariableEmptyButRequired, structField.Name)
		}

		err = common.SetField(value, fieldValue)
		if err != nil {
			return u.e.ErrorOnly(err)
		}
	}

	castedField, isPossibleToCast := element.Addr().Interface().(configService)
	if !isPossibleToCast {
		return nil
	}

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
