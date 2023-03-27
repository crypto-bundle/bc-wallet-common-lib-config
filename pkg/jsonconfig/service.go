package jsonconfig

import (
	"context"
	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
	"os"
)

type targetConfigWrapper struct {
	DependentCfgSrvList []interface{}                 `ignored:"true"`
	castedTarget        easyjson.MarshalerUnmarshaler `ignored:"true"`
	sourceData          []byte                        `ignored:"true"`
	sourceFIlePath      *string                       `ignored:"true"`

	TargetForPrepare interface{}
}

type Service struct {
	secretsSrv secretManagerService

	wrapperConfig *targetConfigWrapper
}

func (m *Service) PrepareFrom(rawJSONData []byte) *Service {
	m.wrapperConfig.sourceData = rawJSONData

	return m
}

func (m *Service) PrepareFromFile(fileDataPath string) *Service {
	m.wrapperConfig.sourceFIlePath = &fileDataPath

	return m
}

func (m *Service) PrepareTo(targetForPrepare interface{}) *Service {
	wrappedTargetConf := &targetConfigWrapper{
		DependentCfgSrvList: make([]interface{}, 0),
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
		default:
			continue
		}
	}

	m.wrapperConfig.DependentCfgSrvList = append(m.wrapperConfig.DependentCfgSrvList, dependenciesList...)

	return m
}

func (m *Service) Do(_ context.Context) error {
	if m.wrapperConfig.sourceFIlePath != nil {
		rawData, err := os.ReadFile(*m.wrapperConfig.sourceFIlePath)
		if err != nil {
			return err
		}
		
		m.wrapperConfig.sourceData = rawData
	}

	r := jlexer.Lexer{Data: m.wrapperConfig.sourceData}

	m.wrapperConfig.castedTarget.UnmarshalEasyJSON(&r)
	err := r.Error()
	if err != nil {
		return err
	}

	secretFillerSrv := &secretFiller{
		dependenciesSrv: m.wrapperConfig.DependentCfgSrvList,
		secretsSrv:      m.secretsSrv,
		target:          m.wrapperConfig.TargetForPrepare,
	}

	err = secretFillerSrv.Process()
	if err != nil {
		return err
	}

	return nil
}
