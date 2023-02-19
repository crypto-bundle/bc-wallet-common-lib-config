package config

type targetConfigWrapper struct {
	dependentCfgSrvList []ConfigService

	targetForPrepare ConfigService
}

func (m *targetConfigWrapper) Prepare() error {
	if m.dependentCfgSrvList == nil || len(m.dependentCfgSrvList) == 0 {
		return m.targetForPrepare.Prepare()
	}

	return m.targetForPrepare.PrepareWith(m.dependentCfgSrvList...)
}

func (m *targetConfigWrapper) PrepareWith(cfgSrv ...ConfigService) error {
	return m.targetForPrepare.PrepareWith(cfgSrv...)
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

func (m *configManager) PrepareTo(targetForPrepare ConfigService) error {
	wrappedTargetConf := &targetConfigWrapper{
		dependentCfgSrvList: nil,
		targetForPrepare:    targetForPrepare,
	}

	cfgVarPool := newConfigVarsPool(m.secretsSrv, wrappedTargetConf)
	err := cfgVarPool.Process()
	if err != nil {
		return err
	}

	return nil
}

func NewConfigManager(
	secretsSrv secretManagerService,
) *configManager {
	return &configManager{
		secretsSrv: secretsSrv,
	}
}
