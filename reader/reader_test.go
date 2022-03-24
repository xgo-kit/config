package reader

import (
	"github.com/stretchr/testify/assert"
	"github.com/xgo-kit/config/source"
	"github.com/xgo-kit/encoding"
	"github.com/xgo-kit/encoding/json"
	"github.com/xgo-kit/encoding/yaml"
	"testing"
)

const (
	jsonData = `{
"app": "test",
"user": {
  "name": "foo"
},
"path": ["/var/log", "/var/www"]
}`

	yamlData = `app: test2
user:
  name: bar
path:
  - /var/log
`
)

func TestReader(t *testing.T) {
	r := NewReader(map[string]encoding.Codec{
		json.Name: encoding.Get(json.Name),
		yaml.Name: encoding.Get(yaml.Name),
	})

	jsonKv := &source.KV{
		Key:    "json",
		Value:  []byte(jsonData),
		Format: "json",
	}
	err := r.Merge(jsonKv)
	assert.NoError(t, err)

	v, ok := r.Value(`app`)
	assert.True(t, ok)
	str, _ := v.String()
	assert.Equal(t, `test`, str)

	v, ok = r.Value(`user`)
	assert.True(t, ok)
	user := struct {
		Name string `json:"name"`
	}{}
	assert.NoError(t, v.Scan(&user))
	assert.Equal(t, `foo`, user.Name)

	v, ok = r.Value(`path`)
	assert.True(t, ok)
	paths := make([]string, 0)
	assert.NoError(t, v.Scan(&paths))
	assert.Len(t, paths, 2)
	assert.Equal(t, `/var/log`, paths[0])
	assert.Equal(t, `/var/www`, paths[1])

	// test merge
	yamlKv := &source.KV{
		Key:    "yaml",
		Value:  []byte(yamlData),
		Format: "yaml",
	}

	assert.NoError(t, r.Merge(yamlKv))

	v, ok = r.Value(`app`)
	assert.True(t, ok)
	str, _ = v.String()
	assert.Equal(t, `test2`, str)

	v, ok = r.Value(`user`)
	assert.True(t, ok)
	assert.NoError(t, v.Scan(&user))
	assert.Equal(t, `bar`, user.Name)

	v, ok = r.Value(`path`)
	assert.True(t, ok)
	paths = make([]string, 0)
	assert.NoError(t, v.Scan(&paths))
	assert.Len(t, paths, 1)
	assert.Equal(t, `/var/log`, paths[0])
}
