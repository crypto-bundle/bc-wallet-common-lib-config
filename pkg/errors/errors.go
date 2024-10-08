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

package errors

import "sync"

//nolint:gochecknoglobals // it's ok
var errorsFmtService errorFormatterService = nil

func InitInternalFmt(fmtSvc errorFormatterService) {
	sync.OnceFunc(func() {
		if errorsFmtService == nil {
			errorsFmtService = fmtSvc
		}
	})
}

func ErrorWithCode(err error, code int) error {
	return errorsFmtService.ErrorWithCode(err, code)
}

func ErrWithCode(err error, code int) error {
	return errorsFmtService.ErrWithCode(err, code)
}

func ErrorGetCode(err error) int {
	return errorsFmtService.ErrorGetCode(err)
}

func ErrGetCode(err error) int {
	return errorsFmtService.ErrGetCode(err)
}

func ErrorNoWrap(err error) error {
	return errorsFmtService.ErrorNoWrap(err)
}

func ErrNoWrap(err error) error {
	return errorsFmtService.ErrNoWrap(err)
}

func ErrorOnly(err error, details ...string) error {
	return errorsFmtService.ErrorOnly(err, details...)
}

func Error(err error, details ...string) error {
	return errorsFmtService.Error(err, details...)
}

func Errorf(err error, format string, args ...interface{}) error {
	return errorsFmtService.Errorf(err, format, args...)
}

func NewError(details ...string) error {
	return errorsFmtService.NewError(details...)
}

func NewErrorf(format string, args ...interface{}) error {
	return errorsFmtService.NewErrorf(format, args...)
}
