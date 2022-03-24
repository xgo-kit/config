package reader

import (
	"bytes"
	"encoding/gob"
	"github.com/bitly/go-simplejson"
	"github.com/imdario/mergo"
	"github.com/xgo-kit/config/source"
	"github.com/xgo-kit/encoding"
	"github.com/xgo-kit/encoding/json"
	"strings"
	"sync"
)

var (
	_         Reader = (*reader)(nil)
	jsonCodec        = encoding.Get(json.Name)
)

type Reader interface {
	Merge(kvs ...*source.KV) error
	Value(key string) (Value, bool)
	Scan(v interface{}) error
}

type reader struct {
	j      *simplejson.Json
	codecs map[string]encoding.Codec
	mutex  sync.Mutex
}

func (r *reader) Merge(kvs ...*source.KV) error {
	r.mutex.Lock()
	defer r.mutex.Unlock()
	merged, err := r.deepCopy()
	if err != nil {
		return err
	}

	for _, kv := range kvs {
		if kv == nil {
			continue
		}

		if len(kv.Value) == 0 {
			continue
		}

		codec := encoding.Get(kv.Format)

		if codec == nil {
			codec = jsonCodec
		}

		data := make(map[string]interface{})

		if err := codec.Unmarshal(kv.Value, &data); err != nil {
			return err
		}

		if err := mergo.Map(&merged, data, mergo.WithOverride); err != nil {
			return err
		}
	}

	b, err := jsonCodec.Marshal(merged)
	if err != nil {
		return err
	}
	return r.j.UnmarshalJSON(b)
}

func (r *reader) Value(key string) (Value, bool) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	paths := strings.Split(key, `.`)

	j := r.j.GetPath(paths...)
	if j == nil {
		return nil, false
	}
	return &value{j}, true
}

func (r *reader) deepCopy() (map[string]interface{}, error) {
	// https://gist.github.com/soroushjp/0ec92102641ddfc3ad5515ca76405f4d
	var buf bytes.Buffer

	m, err := r.j.Map()
	if err != nil {
		return nil, err
	}

	enc := gob.NewEncoder(&buf)
	dec := gob.NewDecoder(&buf)
	err = enc.Encode(m)
	if err != nil {
		return nil, err
	}
	var cp map[string]interface{}
	err = dec.Decode(&cp)
	if err != nil {
		return nil, err
	}
	return cp, nil
}

func (r *reader) Scan(v interface{}) error {
	b, err := r.j.MarshalJSON()
	if err != nil {
		return err
	}
	return jsonCodec.Unmarshal(b, v)
}

func NewReader(codecs map[string]encoding.Codec) Reader {
	return &reader{
		codecs: codecs,
		j:      simplejson.New(),
	}
}

func init() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})
}
