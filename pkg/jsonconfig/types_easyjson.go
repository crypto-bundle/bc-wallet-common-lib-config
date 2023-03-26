// Code generated by easyjson for marshaling/unmarshaling. DO NOT EDIT.

package jsonconfig

import (
	json "encoding/json"
	easyjson "github.com/mailru/easyjson"
	jlexer "github.com/mailru/easyjson/jlexer"
	jwriter "github.com/mailru/easyjson/jwriter"
)

// suppress unused package warning
var (
	_ *json.RawMessage
	_ *jlexer.Lexer
	_ *jwriter.Writer
	_ easyjson.Marshaler
)

func easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(in *jlexer.Lexer, out *SimpleJSONCase) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "int_field_one":
			out.IntFieldOne = int(in.Int())
		case "int_field_tow":
			out.IntFieldTwo = int(in.Int())
		case "int_field_three":
			out.IntFieldThree = int(in.Int())
		case "string_field":
			out.StringField = string(in.String())
		case "float_field":
			out.FloatField = float32(in.Float32())
		case "db_user":
			out.DBUser = string(in.String())
		case "db_password":
			out.DBPassword = string(in.String())
		case "db_name":
			out.DBName = string(in.String())
		case "db_port":
			out.DBPort = string(in.String())
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(out *jwriter.Writer, in SimpleJSONCase) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"int_field_one\":"
		out.RawString(prefix[1:])
		out.Int(int(in.IntFieldOne))
	}
	{
		const prefix string = ",\"int_field_tow\":"
		out.RawString(prefix)
		out.Int(int(in.IntFieldTwo))
	}
	{
		const prefix string = ",\"int_field_three\":"
		out.RawString(prefix)
		out.Int(int(in.IntFieldThree))
	}
	{
		const prefix string = ",\"string_field\":"
		out.RawString(prefix)
		out.String(string(in.StringField))
	}
	{
		const prefix string = ",\"float_field\":"
		out.RawString(prefix)
		out.Float32(float32(in.FloatField))
	}
	{
		const prefix string = ",\"db_user\":"
		out.RawString(prefix)
		out.String(string(in.DBUser))
	}
	{
		const prefix string = ",\"db_password\":"
		out.RawString(prefix)
		out.String(string(in.DBPassword))
	}
	{
		const prefix string = ",\"db_name\":"
		out.RawString(prefix)
		out.String(string(in.DBName))
	}
	{
		const prefix string = ",\"db_port\":"
		out.RawString(prefix)
		out.String(string(in.DBPort))
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v SimpleJSONCase) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v SimpleJSONCase) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *SimpleJSONCase) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *SimpleJSONCase) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig(l, v)
}
func easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(in *jlexer.Lexer, out *MixedJSONCase) {
	isTopLevel := in.IsStart()
	if in.IsNull() {
		if isTopLevel {
			in.Consumed()
		}
		in.Skip()
		return
	}
	in.Delim('{')
	for !in.IsDelim('}') {
		key := in.UnsafeFieldName(false)
		in.WantColon()
		if in.IsNull() {
			in.Skip()
			in.WantComma()
			continue
		}
		switch key {
		case "top_level_field_int":
			out.TopLevelField = uint32(in.Uint32())
		case "list":
			if in.IsNull() {
				in.Skip()
				out.List = nil
			} else {
				in.Delim('[')
				if out.List == nil {
					if !in.IsDelim(']') {
						out.List = make([]*SimpleJSONCase, 0, 8)
					} else {
						out.List = []*SimpleJSONCase{}
					}
				} else {
					out.List = (out.List)[:0]
				}
				for !in.IsDelim(']') {
					var v1 *SimpleJSONCase
					if in.IsNull() {
						in.Skip()
						v1 = nil
					} else {
						if v1 == nil {
							v1 = new(SimpleJSONCase)
						}
						(*v1).UnmarshalEasyJSON(in)
					}
					out.List = append(out.List, v1)
					in.WantComma()
				}
				in.Delim(']')
			}
		default:
			in.SkipRecursive()
		}
		in.WantComma()
	}
	in.Delim('}')
	if isTopLevel {
		in.Consumed()
	}
}
func easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(out *jwriter.Writer, in MixedJSONCase) {
	out.RawByte('{')
	first := true
	_ = first
	{
		const prefix string = ",\"top_level_field_int\":"
		out.RawString(prefix[1:])
		out.Uint32(uint32(in.TopLevelField))
	}
	{
		const prefix string = ",\"list\":"
		out.RawString(prefix)
		if in.List == nil && (out.Flags&jwriter.NilSliceAsEmpty) == 0 {
			out.RawString("null")
		} else {
			out.RawByte('[')
			for v2, v3 := range in.List {
				if v2 > 0 {
					out.RawByte(',')
				}
				if v3 == nil {
					out.RawString("null")
				} else {
					(*v3).MarshalEasyJSON(out)
				}
			}
			out.RawByte(']')
		}
	}
	out.RawByte('}')
}

// MarshalJSON supports json.Marshaler interface
func (v MixedJSONCase) MarshalJSON() ([]byte, error) {
	w := jwriter.Writer{}
	easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(&w, v)
	return w.Buffer.BuildBytes(), w.Error
}

// MarshalEasyJSON supports easyjson.Marshaler interface
func (v MixedJSONCase) MarshalEasyJSON(w *jwriter.Writer) {
	easyjson6601e8cdEncodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(w, v)
}

// UnmarshalJSON supports json.Unmarshaler interface
func (v *MixedJSONCase) UnmarshalJSON(data []byte) error {
	r := jlexer.Lexer{Data: data}
	easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(&r, v)
	return r.Error()
}

// UnmarshalEasyJSON supports easyjson.Unmarshaler interface
func (v *MixedJSONCase) UnmarshalEasyJSON(l *jlexer.Lexer) {
	easyjson6601e8cdDecodeGithubComCryptoBundleBcWalletCommonLibConfigPkgJsonconfig1(l, v)
}
