package reader

import (
	"github.com/bitly/go-simplejson"
	"github.com/xgo-kit/encoding"
	"github.com/xgo-kit/encoding/json"
)

var (
	_ Value = (*value)(nil)
)

type value struct {
	j *simplejson.Json
}

func (v *value) Bytes() ([]byte, error) {
	return v.j.Bytes()
}

func (v *value) Int64() (int64, error) {
	return v.j.Int64()
}

func (v *value) String() (string, error) {
	return v.j.String()
}

func (v *value) Float64() (float64, error) {
	return v.j.Float64()
}

func (v *value) Bool() (bool, error) {
	return v.j.Bool()
}

func (v *value) Scan(val interface{}) error {
	b, err := v.j.MarshalJSON()
	if err != nil {
		return err
	}
	return encoding.Get(json.Name).Unmarshal(b, val)
}
