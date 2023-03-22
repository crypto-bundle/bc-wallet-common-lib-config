package config

import "context"

type targetConfigWrapper struct {
	dependentCfgSrvList []interface{} `ignored:"true"`
	castedTarget        configService `ignored:"true"`

	TargetForPrepare interface{}
}

func (m *targetConfigWrapper) Prepare() error {
	if m.castedTarget == nil {
		return nil
	}

	if m.dependentCfgSrvList == nil || len(m.dependentCfgSrvList) == 0 {
		return m.castedTarget.Prepare()
	}

	return m.castedTarget.PrepareWith(m.dependentCfgSrvList...)
}

func (m *targetConfigWrapper) PrepareWith(cfgSrv ...interface{}) error {
	if m.castedTarget == nil {
		return nil
	}

	err := m.castedTarget.PrepareWith(cfgSrv...)
	if err != nil {
		return err
	}

	return m.castedTarget.Prepare()
}

type configManager struct {
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
		dependentCfgSrvList: make([]interface{}, 0),
		TargetForPrepare:    targetForPrepare,
	}

	castedCfgSrv, isPossibleToCast := targetForPrepare.(configService)
	if isPossibleToCast {
		wrappedTargetConf.castedTarget = castedCfgSrv
	}

	m.wrapperConfig = wrappedTargetConf

	return m
}

func (m *configManager) Do(_ context.Context) error {
	cfgVarPool := newConfigVarsPool(m.secretsSrv, m.wrapperConfig.TargetForPrepare,
		m.wrapperConfig.dependentCfgSrvList)
	err := cfgVarPool.Process()
	if err != nil {
		return err
	}

	err = cfgVarPool.ClearENV()
	if err != nil {
		return err
	}

	return nil
}

func NewConfigManager() *configManager {
	return &configManager{}
}
