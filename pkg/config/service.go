package config

type targetConfigWrapper struct {
	baseConfigSrv    baseConfigService
	targetForPrepare targetConfigService
}

func (m *targetConfigWrapper) Prepare() error {
	if m.baseConfigSrv == nil {
		return m.targetForPrepare.Prepare()
	}

	return m.targetForPrepare.PrepareWith(m.baseConfigSrv)
}

func (m *targetConfigWrapper) PrepareWith(baseConfigSrv baseConfigService) error {
	return m.targetForPrepare.PrepareWith(baseConfigSrv)
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

func (m *configManager) PrepareFrom(targetForPrepare targetConfigService) error {
	wrappedTargetConf := &targetConfigWrapper{
		baseConfigSrv:    m.baseCfgSrv,
		targetForPrepare: targetForPrepare,
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
