package file

import (
	"context"
	"github.com/fsnotify/fsnotify"
	"github.com/xgo-kit/config/source"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

var (
	_ source.Source = (*file)(nil)
)

func NewSource(path string) source.Source {
	return &file{path: path}
}

type file struct {
	path string
}

func (f *file) Load() ([]*source.KV, error) {
	path, err := filepath.Abs(f.path)
	if err != nil {
		return nil, err
	}
	fp, err := os.Open(path)

	if err != nil {
		return nil, err
	}
	defer fp.Close()

	data, err := ioutil.ReadAll(fp)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	kv := &source.KV{
		Key:    fi.Name(),
		Value:  data,
		Format: strings.TrimPrefix(filepath.Ext(path), `.`),
	}

	return []*source.KV{kv}, nil
}

func (f *file) Watch() (source.Watcher, error) {
	fw, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	err = fw.Add(f.path)
	if err != nil {
		return nil, err
	}
	return &watcher{fw: fw, f: f, ctx: context.Background()}, nil
}
