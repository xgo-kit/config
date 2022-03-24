package etcd

import (
	"context"
	source2 "github.com/xgo-kit/config/source"
	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	_ source2.Watcher = (*watcher)(nil)
)

type watcher struct {
	s           *source
	watcherChan clientv3.WatchChan
}

func (w *watcher) Next() ([]*source2.KV, error) {
	select {
	case <-w.s.options.ctx.Done():
		return nil, w.s.options.ctx.Err()
	case <-w.watcherChan:
		return w.s.Load()
	}
}

func (w *watcher) Stop() error {
	_, cancel := context.WithCancel(w.s.options.ctx)
	cancel()
	return nil
}
