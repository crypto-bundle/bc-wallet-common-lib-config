package jsonconfig

import (
	_ "github.com/mailru/easyjson/gen"
	"strconv"
)

//go:generate easyjson types.go

// easyjson:json
type MixedJSONCase struct {
	TopLevelField uint32            `json:"top_level_field_int"`
	List          []*SimpleJSONCase `json:"list"`
}

// easyjson:json
type SimpleJSONCase struct {
	IntFieldOne   int     `json:"int_field_one"`
	IntFieldTwo   int     `json:"int_field_tow"`
	IntFieldThree int     `json:"int_field_three"`
	StringField   string  `json:"string_field"`
	FloatField    float32 `json:"float_field"`
	DBUser        string  `json:"db_user" secret:"true"`
	DBPassword    string  `json:"db_password" secret:"true"`
	DBName        string  `json:"db_name" secret:"true"`
	DBPort        string  `json:"db_port" secret:"true"`
	dbPortAsInt   uint32  `json:"-"`
}

func (v *SimpleJSONCase) GetPort() uint32 {
	return v.dbPortAsInt
}

// Prepare variables to static configuration
func (v *SimpleJSONCase) Prepare() error {
	return nil
}

func (v *SimpleJSONCase) PrepareWith(cfgSrvList ...interface{}) error {
	dbPortAsInt, err := strconv.Atoi(v.DBPort)
	if err != nil {
		return err
	}
	v.dbPortAsInt = uint32(dbPortAsInt)

	return nil
}
