package jsonconfig

import (
	"strconv"

	_ "github.com/mailru/easyjson/gen"
)

//go:generate easyjson types.go

// easyjson:json
type MixedJSONCase struct {
	List          []*SimpleJSONCase `json:"list"`
	TopLevelField uint32            `json:"top_level_field_int"`
}

// easyjson:json
type SimpleJSONCase struct {
	e           errorFormatterService
	StringField string `json:"string_field"`

	DBUser     string `json:"db_user" secret:"true"`
	DBPassword string `json:"db_password" secret:"true"`
	DBName     string `json:"db_name" secret:"true"`
	DBPort     string `json:"db_port" secret:"true"`

	IntFieldOne   int `json:"int_field_one"`
	IntFieldTwo   int `json:"int_field_tow"`
	IntFieldThree int `json:"int_field_three"`

	FloatField float32 `json:"float_field"`

	dbPortAsInt uint32 `json:"-"`
}

func (v *SimpleJSONCase) GetPort() uint32 {
	return v.dbPortAsInt
}

// Prepare variables to static configuration...
func (v *SimpleJSONCase) Prepare() error {
	return nil
}

// PrepareWith struct by passed dependecies list ...
func (v *SimpleJSONCase) PrepareWith(dependenciesList ...interface{}) error {
	for _, cfgSrv := range dependenciesList {
		switch castedDependency := cfgSrv.(type) {
		case errorFormatterService:
			v.e = castedDependency

		default:
			continue
		}
	}

	dbPortAsInt, err := strconv.Atoi(v.DBPort)
	if err != nil {
		return v.e.ErrorOnly(err)
	}

	v.dbPortAsInt = uint32(dbPortAsInt)

	return nil
}
