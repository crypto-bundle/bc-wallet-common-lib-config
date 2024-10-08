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

package common

import (
	"reflect"
	"strconv"
	"strings"
	"time"

	errfmt "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/errors"
)

// SetField - function for case value in struct by field name and reflect value...
// TODO: refactor it - separate by sub-function and move to separated service-component...
//
//nolint:funlen,gocognit,cyclop // it's ok. Need to refactor this function, but now - it's ok.
func SetField(value string, field reflect.Value) error {
	typ := field.Type()

	if typ.Kind() == reflect.Ptr {
		typ = typ.Elem()
		if field.IsNil() {
			field.Set(reflect.New(typ))
		}

		field = field.Elem()
	}

	switch typ.Kind() {
	case reflect.String:
		field.SetString(value)

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		var (
			val int64
			err error
		)

		if field.Kind() == reflect.Int64 && typ.PkgPath() == "time" && typ.Name() == "Duration" {
			var d time.Duration
			d, err = time.ParseDuration(value)
			val = int64(d)
		} else {
			val, err = strconv.ParseInt(value, 0, typ.Bits())
		}

		if err != nil {
			return errfmt.ErrorNoWrap(err)
		}

		field.SetInt(val)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		val, err := strconv.ParseUint(value, 0, typ.Bits())
		if err != nil {
			return errfmt.ErrorNoWrap(err)
		}

		field.SetUint(val)
	case reflect.Bool:
		val, err := strconv.ParseBool(value)
		if err != nil {
			return errfmt.ErrorNoWrap(err)
		}

		field.SetBool(val)

	case reflect.Float32, reflect.Float64:
		val, err := strconv.ParseFloat(value, typ.Bits())
		if err != nil {
			return errfmt.ErrorNoWrap(err)
		}

		field.SetFloat(val)

	case reflect.Slice:
		sliceField := reflect.MakeSlice(typ, 0, 0)
		if typ.Elem().Kind() == reflect.Uint8 {
			sliceField = reflect.ValueOf([]byte(value))
		} else if len(strings.TrimSpace(value)) != 0 {
			vals := strings.Split(value, ",")
			sliceField = reflect.MakeSlice(typ, len(vals), len(vals))

			for i, val := range vals {
				err := SetField(val, sliceField.Index(i))
				if err != nil {
					return errfmt.ErrorNoWrap(err)
				}
			}
		}

		field.Set(sliceField)

	case reflect.Map:
		mapField := reflect.MakeMap(typ)

		if len(strings.TrimSpace(value)) != 0 {
			pairs := strings.Split(value, ",")
			for _, pair := range pairs {
				kvpair := strings.Split(pair, ":")
				if len(kvpair) != 2 {
					return errfmt.NewErrorf("invalid map item: %q", pair)
				}

				pairKey := reflect.New(typ.Key()).Elem()

				err := SetField(kvpair[0], pairKey)
				if err != nil {
					return errfmt.ErrorNoWrap(err)
				}

				elementValue := reflect.New(typ.Elem()).Elem()

				err = SetField(kvpair[1], elementValue)
				if err != nil {
					return errfmt.ErrorNoWrap(err)
				}

				mapField.SetMapIndex(pairKey, elementValue)
			}
		}

		field.Set(mapField)

	default:
		return nil
	}

	return nil
}
