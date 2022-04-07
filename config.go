package config

import (
	"errors"
	"github.com/xgo-kit/config/reader"
	"github.com/xgo-kit/config/source"
	"github.com/xgo-kit/encoding"
	"github.com/xgo-kit/encoding/json"
	"github.com/xgo-kit/encoding/yaml"
	"sync"
	"time"
)

var (
	_ Config = (*config)(nil)

	ErrNotFound = errors.New(`key not found`)
)

type Config interface {
	Load() error
	Value(key string) reader.Value
	Scan(v interface{}) error
	Observer(key string, o Observer) error
}

type Observer func(key string, value reader.Value)

type config struct {
	opts      options
	observers sync.Map
}

func (c *config) Load() error {
	for _, s := range c.opts.sources {
		kvs, err := s.Load()
		if err != nil {
			return err
		}
		err = c.opts.reader.Merge(kvs...)
		if err != nil {
			return err
		}
		w, err := s.Watch()
		if err != nil {
			return err
		}

		go func(w source.Watcher) {
			for {
				kvs, err := w.Next()
				if err != nil {
					return
				}

				err = c.opts.reader.Merge(kvs...)
				if err != nil {
					return
				}

				c.observers.Range(func(key, value any) bool {
					k := key.(string)
					v := value.(Observer)

					cv := c.Value(k)
					if _, ok := cv.(*reader.NullValue); ok {
						return true
					}

					v(k, cv)

					return true
				})

				time.Sleep(time.Second)
			}
		}(w)
	}
	return nil
}

func (c *config) Value(key string) reader.Value {
	v, ok := c.opts.reader.Value(key)
	if !ok {
		return &reader.NullValue{}
	}
	return v
}

func (c *config) Observer(key string, o Observer) error {
	v := c.Value(key)
	if _, ok := v.(*reader.NullValue); ok {
		return ErrNotFound
	}
	c.observers.Store(key, o)
	return nil
}

func (c *config) Scan(v interface{}) error {
	return c.opts.reader.Scan(v)
}

func NewConfig(opts ...Option) Config {
	o := options{
		reader: defaultReader(),
	}

	for _, opt := range opts {
		opt(&o)
	}

	return &config{
		opts: o,
	}
}

func defaultReader() reader.Reader {
	return reader.NewReader(defaultCodecs())
}

func defaultCodecs() map[string]encoding.Codec {
	return map[string]encoding.Codec{
		yaml.Name: encoding.Get(yaml.Name),
		json.Name: encoding.Get(json.Name),
	}
}
