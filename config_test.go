package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/xgo-kit/config/source"
	"testing"
)

const (
	testData = `
app: t1
paths: 
  - /var/test
  - /var/log
user: 
  name: a1
  password: "123"
`
)

func TestConfig(t *testing.T) {
	c := NewConfig(WithSource(&testSource{data: []byte(testData)}))
	err := c.Load()
	assert.NoError(t, err)

	v := c.Value(`app`)
	str, _ := v.String()
	assert.Equal(t, `t1`, str)

	v = c.Value(`user.name`)
	username, _ := v.String()
	assert.Equal(t, `a1`, username)

	tc := testConfig{}
	assert.NoError(t, c.Scan(&tc))
	assert.Equal(t, `t1`, tc.App)
	assert.Len(t, tc.Paths, 2)
	assert.Equal(t, tc.User.Password, `123`)
}

type testConfig struct {
	App   string
	Paths []string
	User  struct {
		Name     string
		Password string
	}
}

type testSource struct {
	data   []byte
	Update chan bool
}

func (ts *testSource) Load() ([]*source.KV, error) {
	return []*source.KV{{
		Key:    "config.yaml",
		Value:  ts.data,
		Format: `yaml`,
	}}, nil
}

func (ts *testSource) Watch() (source.Watcher, error) {
	return &testWatcher{
		s:      ts,
		update: make(chan bool),
		close:  make(chan bool),
	}, nil
}

type testWatcher struct {
	s      *testSource
	update chan bool
	close  chan bool
}

func (w *testWatcher) Next() ([]*source.KV, error) {
	select {
	case _, ok := <-w.update:
		if !ok {
			return nil, nil
		}
		return w.s.Load()
	case <-w.close:
		return nil, nil
	}
}

func (w *testWatcher) Stop() error {
	close(w.close)
	return nil
}
