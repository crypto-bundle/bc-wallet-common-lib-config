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

	return m.castedTarget.PrepareWith(cfgSrv...)
}

type configManager struct {
	secretsSrv secretManagerService

	wrapperConfig *targetConfigWrapper
}

func (m *configManager) With(cfgSrvList ...interface{}) *configManager {
	cloned := *m
	cloned.wrapperConfig.dependentCfgSrvList = append(cloned.wrapperConfig.dependentCfgSrvList, cfgSrvList...)

	return &cloned
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

	return nil
}

func (m *configManager) Do(ctx context.Context) error {
	cfgVarPool := newConfigVarsPool(m.secretsSrv, m.wrapperConfig)
	err := cfgVarPool.Process()
	if err != nil {
		return err
	}

	return nil
}

func NewConfigManager() *configManager {
	return &configManager{}
}
