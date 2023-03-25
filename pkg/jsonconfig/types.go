package jsonconfig

import _ "github.com/mailru/easyjson/gen"

//go:generate easyjson types.go

// easyjson:json
type SimpleJSONCase struct {
	IntFieldOne   int     `json:"int_field_one"`
	IntFieldTwo   int     `json:"int_field_tow"`
	IntFieldThree int     `json:"int_field_three"`
	StringField   string  `json:"string_field"`
	FloatField    float32 `json:"float_field"`
	DBUser        string  `json:"db_user" secret:"true" secret_name:"DATABASE_USER"`
	DBPassword    string  `json:"db_password" secret:"true" secret_name:"DATABASE_PASSWORD"`
	DBName        string  `json:"db_name" secret:"true" secret_name:"DATABASE_NAME"`
	DBPort        uint32  `json:"db_port" secret:"true" secret_name:"DATABASE_PORT"`
}
