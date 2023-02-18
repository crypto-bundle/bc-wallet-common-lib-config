package config

import (
	"errors"
	"fmt"
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
	targetConfigSrv targetConfigService
	secretsSrv      secretManagerService

	envVariablesNameCount uint16
	envVariablesNameList  []string
	envVariablesList      []field

	secretVariablesCount uint16
	secretVariablesList  []field
}

func (u *configVariablesPool) addSecretVariable(variable field) error {
	err := setField(variable.Value, variable.RfValue)
	if err != nil {
		return err
	}

	u.secretVariablesCount++
	u.secretVariablesList = append(u.envVariablesList, variable)

	return nil
}

func (u *configVariablesPool) addEnvVariable(variable field) error {
	err := setField(variable.Value, variable.RfValue)
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
	s = s.Elem()
	if s.Kind() != reflect.Struct {
		return ErrPassedStructMustBeAStructPointer
	}
	elemType := s.Type()

	// iterate over struct fields
	for i := 0; i < s.NumField(); i++ {
		fv := s.Field(i)        // reflect.RfValue
		sf := elemType.Field(i) // struct field info
		if !fv.CanSet() {
			continue
		}

		ignored, _ := strconv.ParseBool(sf.Tag.Get(tagIgnored))
		if ignored {
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

			castedField, isPossibleToCast := fv.Interface().(targetConfigService)
			if isPossibleToCast {
				prepErr := castedField.Prepare()
				if prepErr != nil {
					return prepErr
				}
			}

			continue
		}

		isSecret, err := strconv.ParseBool(sf.Tag.Get(tagSecret))
		if err != nil {
			return err
		}

		isRequired, err := strconv.ParseBool(sf.Tag.Get(tagRequired))
		if err != nil {
			return err
		}

		if isSecret {
			value, isExists := u.secretsSrv.GetByName(sf.Name)
			if !isExists && isRequired {
				return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, sf.Name)
			}

			f := field{
				Name:    sf.Name,
				RfValue: fv,
				RfTags:  sf.Tag,
				Value:   value,
			}

			err = u.addSecretVariable(f)
			if err != nil {
				return err
			}

			continue
		}

		envConfigKey := sf.Tag.Get(tagEnvconfig)
		value, isExists := os.LookupEnv(envConfigKey)
		if !isExists && isRequired {
			return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, sf.Name)
		}

		f := field{
			Name:    sf.Name,
			RfValue: fv,
			RfTags:  sf.Tag,
			Value:   value,
		}

		err = u.addEnvVariable(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func newConfigVarsPool(secretSrv secretManagerService,
	processedConfig targetConfigService,
) *configVariablesPool {
	return &configVariablesPool{
		targetConfigSrv: processedConfig,
		secretsSrv:      secretSrv,

		envVariablesNameCount: 0,
		envVariablesNameList:  make([]string, 0),
		envVariablesList:      make([]field, 0),

		secretVariablesCount: 0,
		secretVariablesList:  make([]field, 0),
	}
}
