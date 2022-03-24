package file

import (
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

const (
	testYaml = `
app:
  name: foo
db:
  driver: mysql
  addr: 127.0.0.1
`

	updateYaml = `
app:
  name: bar
db:
  driver: pgsql
  addr: localhost
`
)

func TestFile(t *testing.T) {
	var (
		path       = `./test.yaml`
		testdata   = []byte(testYaml)
		updatedata = []byte(updateYaml)
	)

	defer os.Remove(path)

	err := os.WriteFile(path, testdata, os.ModePerm)
	assert.NoError(t, err)

	f := NewSource(path)
	kvs, err := f.Load()
	assert.NoError(t, err)
	assert.NotNil(t, kvs[0])
	assert.Equal(t, `test.yaml`, kvs[0].Key)
	assert.Equal(t, testdata, kvs[0].Value)
	assert.Equal(t, `yaml`, kvs[0].Format)

	w, err := f.Watch()
	assert.NoError(t, err)

	// test watcher
	fp, err := os.OpenFile(path, os.O_RDWR, os.ModePerm)
	assert.NoError(t, err)
	defer fp.Close()

	_, err = fp.Write(updatedata)
	assert.NoError(t, err)

	kvs, err = w.Next()
	assert.NoError(t, err)
	assert.Equal(t, updatedata, kvs[0].Value)
}
