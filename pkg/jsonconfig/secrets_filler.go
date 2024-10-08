package jsonconfig

import (
	"errors"
	"reflect"
	"strconv"
	"strings"

	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
)

var (
	ErrPassedStructMustBeAPointer       = errors.New("must be a pointer")
	ErrPassedStructMustBeAStructPointer = errors.New("must be a struct pointer")
	ErrVariableEmptyButRequired         = errors.New("variables is empty and has required tag")
	ErrWrongSecretStringFormat          = errors.New("wrong secret string format")
)

type secretFiller struct {
	e              errorFormatterService
	secretsDataSvc secretManagerService

	dependenciesSvc []interface{}

	target interface{}
}

func (u *secretFiller) Process() error {
	return u.processFields(u.target)
}

// extractFields returns information of the struct fields, including nested structures
// based on https://github.com/kelseyhightower/envconfig
func (u *secretFiller) processFields(target interface{}) error {
	s := reflect.ValueOf(target)

	// must be a pointer
	switch s.Kind() {
	case reflect.Slice:
	case reflect.Ptr:
	default:
		return u.e.ErrorOnly(ErrPassedStructMustBeAPointer)
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

		// recursively process nested struct
		if fv.Kind() == reflect.Struct && fv.CanInterface() {
			processErr := u.processFields(fv.Addr().Interface())
			if processErr != nil {
				return u.e.ErrorOnly(processErr)
			}

			continue
		}

		if fv.Kind() == reflect.Slice && fv.CanInterface() {
			for j := 0; j < fv.Len(); j++ {
				item := fv.Index(j)
				v := reflect.Indirect(item)
				if v.Kind() == reflect.Struct {
					processErr := u.processFields(v.Addr().Interface())
					if processErr != nil {
						return u.e.ErrorOnly(processErr)
					}

					continue
				}
			}
		}

		var isSecret = false
		boolVarSrt, isTagExists := sf.Tag.Lookup(common.TagSecret)
		if !isTagExists {
			continue
		}

		boolVar, err := strconv.ParseBool(boolVarSrt)
		if err != nil {
			return u.e.ErrorOnly(err)
		}
		isSecret = boolVar

		if !isSecret {
			continue
		}

		value := fv.String()
		if !strings.HasPrefix(value, "!secret:") {
			continue
		}

		separated := strings.Split(value, ":")
		if len(separated) > 2 {
			return u.e.ErrorOnly(ErrWrongSecretStringFormat)
		}

		secretKey := separated[1]
		value, isExists := u.secretsDataSvc.GetByName(secretKey)
		if !isExists {
			return u.e.ErrorOnly(ErrVariableEmptyButRequired, sf.Name)
		}

		err = common.SetField(value, fv)
		if err != nil {
			return u.e.ErrorOnly(err)
		}

	}

	castedField, isPossibleToCast := element.Addr().Interface().(configService)
	if !isPossibleToCast {
		return nil
	}

	if u.dependenciesSvc != nil {
		prepErr := castedField.PrepareWith(u.dependenciesSvc...)
		if prepErr != nil {
			return u.e.ErrorOnly(prepErr)
		}
	}

	prepErr := castedField.Prepare()
	if prepErr != nil {
		return u.e.ErrorOnly(prepErr)
	}

	return nil
}
