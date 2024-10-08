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

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type targetConfigWrapper struct {
	castedTarget easyjson.MarshalerUnmarshaler `ignored:"true"`

	TargetForPrepare    interface{}
	sourceFilePath      *string       `ignored:"true"`
	DependentCfgSrvList []interface{} `ignored:"true"`
	sourceData          []byte        `ignored:"true"`
}

type Service struct {
	e          errorFormatterService
	secretsSrv secretManagerService

	wrapperConfig *targetConfigWrapper
}

func (m *Service) PrepareFrom(rawJSONData []byte) *Service {
	m.wrapperConfig.sourceData = rawJSONData

	return m
}

func (m *Service) PrepareFromFile(fileDataPath string) *Service {
	m.wrapperConfig.sourceFilePath = &fileDataPath

	return m
}

func (m *Service) PrepareTo(targetForPrepare interface{}) *Service {
	wrappedTargetConf := &targetConfigWrapper{
		DependentCfgSrvList: make([]interface{}, 0),
		castedTarget:        nil,
		sourceData:          nil,
		sourceFilePath:      nil,
		TargetForPrepare:    targetForPrepare,
	}

	castedCfgSrv, isPossibleToCast := targetForPrepare.(easyjson.MarshalerUnmarshaler)
	if isPossibleToCast {
		wrappedTargetConf.castedTarget = castedCfgSrv
	}

	m.wrapperConfig = wrappedTargetConf

	return m
}

func (m *Service) With(dependenciesList ...interface{}) *Service {
	for _, cfgSrv := range dependenciesList {
		switch castedDependency := cfgSrv.(type) {
		case secretManagerService:
			m.secretsSrv = castedDependency
		case errorFormatterService:
			m.e = castedDependency

		default:
			continue
		}
	}

	m.wrapperConfig.DependentCfgSrvList = append(m.wrapperConfig.DependentCfgSrvList, dependenciesList...)

	return m
}

func (m *Service) Do(_ context.Context) error {
	if m.wrapperConfig.sourceFilePath != nil {
		rawData, err := os.ReadFile(*m.wrapperConfig.sourceFilePath)
		if err != nil {
			return m.e.ErrorOnly(err)
		}

		m.wrapperConfig.sourceData = rawData
	}

	JSONLexer := jlexer.Lexer{
		Data:              m.wrapperConfig.sourceData,
		UseMultipleErrors: false,
	}

	m.wrapperConfig.castedTarget.UnmarshalEasyJSON(&JSONLexer)

	err := JSONLexer.Error()
	if err != nil {
		return m.e.ErrorOnly(err)
	}

	secretDataFillerSvc := &secretFiller{
		e:               m.e,
		dependenciesSvc: m.wrapperConfig.DependentCfgSrvList,
		secretsDataSvc:  m.secretsSrv,
		target:          m.wrapperConfig.TargetForPrepare,
	}

	err = secretDataFillerSvc.Process()
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	return nil
}
