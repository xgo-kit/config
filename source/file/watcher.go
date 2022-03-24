package file

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/xgo-kit/config/source"
)

var (
	_ source.Watcher = (*watcher)(nil)
)

type watcher struct {
	ctx context.Context
	fw  *fsnotify.Watcher
	f   *file
}

func (w *watcher) Next() ([]*source.KV, error) {
	select {
	case <-w.ctx.Done():
		return nil, w.ctx.Err()
	case evt := <-w.fw.Events:
		w.f.path = evt.Name
		return w.f.Load()
	case err := <-w.fw.Errors:
		return nil, err
	}
}

func (w *watcher) Stop() error {
	return w.fw.Close()
}
