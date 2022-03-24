package etcd

import "context"

type options struct {
	ctx    context.Context
	path   string
	prefix bool
}

type Option func(o *options)

func WithContext(ctx context.Context) Option {
	return func(o *options) {
		o.ctx = ctx
	}
}

func WithPath(path string) Option {
	return func(o *options) {
		o.path = path
	}
}
