package jsonconfig

import (
	"context"
	"os"

	"github.com/mailru/easyjson"
	"github.com/mailru/easyjson/jlexer"
)

type targetConfigWrapper struct {
	DependentCfgSrvList []interface{}                 `ignored:"true"`
	castedTarget        easyjson.MarshalerUnmarshaler `ignored:"true"`
	sourceData          []byte                        `ignored:"true"`
	sourceFIlePath      *string                       `ignored:"true"`

	TargetForPrepare interface{}
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
	if m.wrapperConfig.sourceFIlePath != nil {
		rawData, err := os.ReadFile(*m.wrapperConfig.sourceFIlePath)
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
