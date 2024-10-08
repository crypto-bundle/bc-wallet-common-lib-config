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
