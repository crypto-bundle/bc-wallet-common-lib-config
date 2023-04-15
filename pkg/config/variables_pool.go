package config

import (
	"errors"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
	"os"
	"reflect"
	"strconv"
)

var (
	ErrPassedStructMustBeAPointer       = errors.New("must be a pointer")
	ErrPassedStructMustBeAStructPointer = errors.New("must be a struct pointer")
	ErrVariableEmptyButRequired         = errors.New("variables is empty and has required tag")
)

type configVariablesPool struct {
	targetConfigSrv interface{}
	dependenciesSrv []interface{}
	secretsSrv      secretManagerService

	envVariablesNameCount uint16
	envVariablesNameList  []string
	envVariablesList      []common.Field

	secretVariablesCount uint16
	secretVariablesList  []common.Field
}

func (u *configVariablesPool) addSecretVariable(variable common.Field) error {
	err := common.SetField(variable.Value, variable.RfValue)
	if err != nil {
		return err
	}

	u.secretVariablesCount++
	u.secretVariablesList = append(u.envVariablesList, variable)

	return nil
}

func (u *configVariablesPool) addEnvVariable(variable common.Field) error {
	err := common.SetField(variable.Value, variable.RfValue)
	if err != nil {
		return err
	}

	u.envVariablesNameCount++
	u.envVariablesList = append(u.envVariablesList, variable)

	return nil
}

func (u *configVariablesPool) Process() error {
	err := u.processFields(u.targetConfigSrv)
	if err != nil {
		return err
	}

	return nil
}

// extractFields returns information of the struct fields, including nested structures
// based on https://github.com/kelseyhightower/envconfig
func (u *configVariablesPool) processFields(target interface{}) error {
	s := reflect.ValueOf(target)

	// must be a pointer
	if s.Kind() != reflect.Ptr {
		return ErrPassedStructMustBeAPointer
	}

	// pointer must refer to structure
	element := s.Elem()
	elemType := element.Type()

	// iterate over struct fields
	numFields := elemType.NumField()
	for i := 0; i < numFields; i++ {
		fv := element.Field(i)  // reflect.RfValue
		sf := elemType.Field(i) // struct field info
		if !fv.CanSet() {
			continue
		}

		isIgnored, _ := strconv.ParseBool(sf.Tag.Get(common.TagIgnored))
		if isIgnored {
			continue
		}

		// unfold pointers
		for fv.Kind() == reflect.Ptr {
			if fv.IsNil() {
				if fv.Type().Elem().Kind() != reflect.Struct {
					// nil pointer to a non-struct: leave it alone
					break
				}
				// nil pointer to struct: create a zero instance
				fv.Set(reflect.New(fv.Type().Elem()))
			}
			fv = fv.Elem()
		}

		// recursively process nested struct
		if fv.Kind() == reflect.Struct && fv.CanInterface() {
			processErr := u.processFields(fv.Addr().Interface())
			if processErr != nil {
				return processErr
			}

			continue
		}

		var isSecret = false
		boolVarSrt, isTagExists := sf.Tag.Lookup(common.TagSecret)
		if isTagExists {
			boolVar, err := strconv.ParseBool(boolVarSrt)
			if err != nil {
				return err
			}

			isSecret = boolVar
		}

		var isRequired = false
		boolVarSrt, isTagExists = sf.Tag.Lookup(common.TagRequired)
		if isTagExists {
			boolVar, err := strconv.ParseBool(boolVarSrt)
			if err != nil {
				return err
			}

			isRequired = boolVar
		}

		if isSecret {
			envConfigKey := sf.Tag.Get(common.TagEnvconfig)
			value, isExists := u.secretsSrv.GetByName(envConfigKey)
			if !isExists && isRequired {
				return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, sf.Name)
			}

			f := common.Field{
				Name:    sf.Name,
				RfValue: fv,
				RfTags:  sf.Tag,
				Value:   value,
			}

			addErr := u.addSecretVariable(f)
			if addErr != nil {
				return addErr
			}

			continue
		}

		envConfigKey := sf.Tag.Get(common.TagEnvconfig)
		value, isEnvVariableExists := os.LookupEnv(envConfigKey)
		if !isEnvVariableExists && isRequired {
			return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, sf.Name)
		}

		defaultValue, hasDefaultValue := sf.Tag.Lookup(common.TagDefault)
		if !isEnvVariableExists && hasDefaultValue {
			value = defaultValue
		}

		f := common.Field{
			Name:    sf.Name,
			RfValue: fv,
			RfTags:  sf.Tag,
			Value:   value,
		}

		addErr := u.addEnvVariable(f)
		if addErr != nil {
			return addErr
		}
	}

	castedField, isPossibleToCast := element.Addr().Interface().(dependentConfigService)
	if isPossibleToCast {
		if u.dependenciesSrv != nil {
			prepErr := castedField.PrepareWith(u.dependenciesSrv...)
			if prepErr != nil {
				return prepErr
			}
		}

		prepErr := castedField.Prepare()
		if prepErr != nil {
			return prepErr
		}

		return nil
	}

	castedConfigField, isPossibleToCast := element.Addr().Interface().(configService)
	if isPossibleToCast {
		prepErr := castedConfigField.Prepare()
		if prepErr != nil {
			return prepErr
		}
	}

	return nil
}

func (u *configVariablesPool) ClearENV() error {
	for i := uint16(0); i != u.envVariablesNameCount; i++ {
		envField := u.envVariablesList[i]
		err := os.Unsetenv(envField.Name)
		if err != nil {
			return err
		}
	}

	return nil
}

func newConfigVarsPool(secretSrv secretManagerService,
	processedConfig interface{},
	dependenciesSrvList []interface{},
) *configVariablesPool {
	return &configVariablesPool{
		dependenciesSrv: dependenciesSrvList,
		targetConfigSrv: processedConfig,
		secretsSrv:      secretSrv,

		envVariablesNameCount: 0,
		envVariablesNameList:  make([]string, 0),
		envVariablesList:      make([]common.Field, 0),

		secretVariablesCount: 0,
		secretVariablesList:  make([]common.Field, 0),
	}
}
