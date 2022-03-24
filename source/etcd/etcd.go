package etcd

import (
	"context"
	source2 "github.com/xgo-kit/config/source"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.etcd.io/etcd/client/v3"
	"path/filepath"
	"strings"
)

var (
	_ source2.Source = (*source)(nil)
)

type source struct {
	client  *clientv3.Client
	options options
}

func (s *source) Load() ([]*source2.KV, error) {
	var o []clientv3.OpOption

	if s.options.prefix {
		o = append(o, clientv3.WithPrefix())
	}

	resp, err := s.client.Get(s.options.ctx, s.options.path, o...)
	if err != nil {
		return nil, err
	}

	return s.load(resp.Kvs)
}

func (s *source) load(mkvs []*mvccpb.KeyValue) ([]*source2.KV, error) {
	kvs := make([]*source2.KV, 0, len(mkvs))
	for _, kv := range mkvs {
		kvs = append(kvs, &source2.KV{
			Key:    string(kv.Key),
			Value:  kv.Value,
			Format: strings.TrimPrefix(filepath.Ext(string(kv.Key)), `.`),
		})
	}

	return kvs, nil
}

func (s *source) Watch() (source2.Watcher, error) {
	w := &watcher{
		s: s,
	}
	// 这里必须提前拿到etcd 的watch事件，否则watcher执行的时候拿不到通知
	w.watcherChan = s.client.Watch(w.s.options.ctx, w.s.options.path)
	return w, nil
}

func NewSource(c *clientv3.Client, opts ...Option) source2.Source {
	o := options{
		ctx: context.Background(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &source{client: c, options: o}
}
