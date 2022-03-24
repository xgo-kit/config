package config

import (
	"github.com/xgo-kit/config/reader"
	"github.com/xgo-kit/config/source"
)

type options struct {
	sources []source.Source
	reader  reader.Reader
}

type Option func(o *options)

func WithSource(sources ...source.Source) Option {
	return func(o *options) {
		o.sources = sources
	}
}

func WithReader(r reader.Reader) Option {
	return func(o *options) {
		o.reader = r
	}
}
