package etcd

import (
	"context"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
	"testing"
	"time"
)

func TestSource(t *testing.T) {
	c, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{`127.0.0.1:23790`},
		Context:     context.Background(),
		DialTimeout: 10 * time.Second,
	})
	assert.NoError(t, err)
	defer c.Close()

	s := NewSource(c, WithPath(`/api/config.json`))
	_, err = c.Put(context.Background(), `/api/config.json`, `test1`)
	assert.NoError(t, err)

	kv, err := s.Load()
	assert.NoError(t, err)
	assert.Equal(t, kv[0].Key, `/api/config.json`)
	assert.Equal(t, string(kv[0].Value), `test1`)
	assert.Equal(t, kv[0].Format, `json`)

	w, err := s.Watch()
	assert.NoError(t, err)

	defer func() {
		w.Stop()
	}()

	_, err = c.Put(context.Background(), `/api/config.json`, `test2`)
	assert.NoError(t, err)

	kv, err = w.Next()
	assert.NoError(t, err)

	assert.Equal(t, kv[0].Key, `/api/config.json`)
	assert.Equal(t, string(kv[0].Value), `test2`)
}
