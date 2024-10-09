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
	"context"
	"os"

	errfmt "github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/errors"

	"github.com/joho/godotenv"
)

type targetConfigWrapper struct {
	e                   errorFormatterService
	castedTarget        dependentConfigService `ignored:"true"`
	TargetForPrepare    interface{}
	dependentCfgSrvList []interface{} `ignored:"true"`
}

func (m *targetConfigWrapper) Prepare() error {
	if m.castedTarget == nil {
		return nil
	}

	if len(m.dependentCfgSrvList) == 0 {
		err := m.castedTarget.Prepare()
		if err != nil {
			return m.e.ErrorNoWrap(err)
		}
	}

	err := m.castedTarget.PrepareWith(m.dependentCfgSrvList...)
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	return nil
}

func (m *targetConfigWrapper) PrepareWith(cfgDependenciesSvcList ...interface{}) error {
	if m.castedTarget == nil {
		return nil
	}

	err := m.castedTarget.PrepareWith(cfgDependenciesSvcList...)
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	err = m.castedTarget.Prepare()
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	return nil
}

type configManager struct {
	e errorFormatterService

	secretsSrv secretManagerService

	wrapperConfig *targetConfigWrapper
}

func (m *configManager) With(dependenciesList ...interface{}) *configManager {
	for _, cfgSrv := range dependenciesList {
		switch castedDependency := cfgSrv.(type) {
		case secretManagerService:
			m.secretsSrv = castedDependency
		default:
			continue
		}
	}

	m.wrapperConfig.dependentCfgSrvList = append(m.wrapperConfig.dependentCfgSrvList, dependenciesList...)

	return m
}

func (m *configManager) PrepareTo(targetForPrepare interface{}) *configManager {
	wrappedTargetConf := &targetConfigWrapper{
		e:                   m.e,
		dependentCfgSrvList: make([]interface{}, 0),
		castedTarget:        nil, // will be filled later by targetForPrepare
		TargetForPrepare:    targetForPrepare,
	}

	castedCfgSrv, isPossibleToCast := targetForPrepare.(dependentConfigService)
	if isPossibleToCast {
		wrappedTargetConf.castedTarget = castedCfgSrv
	}

	m.wrapperConfig = wrappedTargetConf

	return m
}

func (m *configManager) Do(_ context.Context) error {
	cfgVarPool := newConfigVarsPool(m.e, m.secretsSrv, m.wrapperConfig.TargetForPrepare,
		m.wrapperConfig.dependentCfgSrvList)

	err := cfgVarPool.Process()
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	err = cfgVarPool.ClearENV()
	if err != nil {
		return m.e.ErrorNoWrap(err)
	}

	return nil
}

func NewConfigManager(errFmtSvc errorFormatterService) *configManager {
	return &configManager{
		e:             errFmtSvc,
		secretsSrv:    nil,
		wrapperConfig: nil,
	}
}

func LoadLocalEnvIfDev() error {
	value, isEnvVariableExists := os.LookupEnv(AppEnvironmentNameVariable)
	if !isEnvVariableExists {
		return errfmt.ErrorOnly(ErrVariableEmptyButRequired, AppEnvironmentNameVariable)
	}

	if value == EnvDev || value == EnvLocal {
		envFilePath, isExists := os.LookupEnv(AppEnvFilePathVariableName)
		if !isExists {
			return errfmt.ErrorOnly(ErrVariableEmptyButRequired, AppEnvFilePathVariableName)
		}

		loadErr := godotenv.Load(envFilePath)
		if loadErr != nil {
			return errfmt.ErrorOnly(loadErr)
		}
	}

	return nil
}

func LoadEnvFromFile(filePath string) error {
	if filePath == "" {
		return LoadLocalEnvIfDev()
	}

	loadErr := godotenv.Load(filePath)
	if loadErr != nil {
		return errfmt.ErrorOnly(loadErr)
	}

	return nil
}
