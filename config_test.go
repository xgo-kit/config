package config

import (
	"github.com/stretchr/testify/assert"
	"github.com/xgo-kit/config/reader"
	"github.com/xgo-kit/config/source"
	"sync"
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

	updatedData = `
app: t2
paths: 
  - /var/log
user: 
  name: a2
  password: "222"
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

func TestObserver(t *testing.T) {
	ts := &testSource{data: []byte(testData), update: make(chan bool)}
	c := NewConfig(WithSource(ts))
	err := c.Load()
	assert.NoError(t, err)

	var username string
	paths := make([]string, 0)

	wg := sync.WaitGroup{}
	wg.Add(2)
	err = c.Observer(`user.name`, func(key string, value reader.Value) {
		username, _ = value.String()
		wg.Done()
	})
	assert.NoError(t, err)

	err = c.Observer(`paths`, func(key string, value reader.Value) {
		_ = value.Scan(&paths)
		wg.Done()
	})

	// test observer
	ts.data = []byte(updatedData)
	ts.update <- true

	wg.Wait()
	assert.Equal(t, `a2`, username)
	assert.Len(t, paths, 1)
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
	update chan bool
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
		s:     ts,
		close: make(chan bool),
	}, nil
}

type testWatcher struct {
	s     *testSource
	close chan bool
}

func (w *testWatcher) Next() ([]*source.KV, error) {
	select {
	case _, ok := <-w.s.update:
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
