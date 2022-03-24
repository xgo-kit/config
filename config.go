package config

import (
	"github.com/xgo-kit/config/reader"
	"github.com/xgo-kit/config/source"
	"github.com/xgo-kit/encoding"
	"github.com/xgo-kit/encoding/json"
	"github.com/xgo-kit/encoding/yaml"
	"time"
)

var (
	_ Config = (*config)(nil)
)

type Config interface {
	Load() error
	Value(key string) reader.Value
	Watch() (source.Watcher, error)
	Scan(v interface{}) error
}

type config struct {
	opts options
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

func (c *config) Watch() (source.Watcher, error) {
	//TODO implement me
	panic("implement me")
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
