package config

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
	baseCfgSrv baseConfigService
	secretsSrv secretManagerService
}

func (m *configManager) With(basCfgSrv baseConfigService) *configManager {
	cloned := *m
	cloned.baseCfgSrv = basCfgSrv

	return &cloned
}

func (m *configManager) PrepareTo(targetForPrepare interface{}) error {
	wrappedTargetConf := &targetConfigWrapper{
		dependentCfgSrvList: nil,
		TargetForPrepare:    targetForPrepare,
	}

	castedCfgSrv, isPossibleToCast := targetForPrepare.(configService)
	if isPossibleToCast {
		wrappedTargetConf.castedTarget = castedCfgSrv
	}

	cfgVarPool := newConfigVarsPool(m.secretsSrv, targetForPrepare)
	err := cfgVarPool.Process()
	if err != nil {
		return err
	}

	return nil
}

func NewConfigManager() *configManager {
	return &configManager{}
}
