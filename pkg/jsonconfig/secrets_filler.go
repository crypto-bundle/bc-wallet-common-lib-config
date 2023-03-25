package jsonconfig

import (
	"errors"
	"fmt"
	"github.com/crypto-bundle/bc-wallet-common-lib-config/pkg/common"
	"reflect"
	"strconv"
)

var (
	ErrPassedStructMustBeAPointer       = errors.New("must be a pointer")
	ErrPassedStructMustBeAStructPointer = errors.New("must be a struct pointer")
	ErrVariableEmptyButRequired         = errors.New("variables is empty and has required tag")
)

type secretFiller struct {
	secretsSrv      secretManagerService
	dependenciesSrv []interface{}

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
		if !isTagExists {
			continue
		}

		boolVar, err := strconv.ParseBool(boolVarSrt)
		if err != nil {
			return err
		}
		isSecret = boolVar

		if !isSecret {
			continue
		}

		secretKey := sf.Tag.Get(common.TagSecretName)
		value, isExists := u.secretsSrv.GetByName(secretKey)
		if !isExists {
			return fmt.Errorf("%w: %s", ErrVariableEmptyButRequired, sf.Name)
		}

		err = common.SetField(value, fv)
		if err != nil {
			return err
		}

	}

	castedField, isPossibleToCast := element.Addr().Interface().(configService)
	if !isPossibleToCast {
		return nil
	}

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
