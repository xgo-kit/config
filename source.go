package config

type Source interface {
	Load() ([]*KV, error)
	watch
}

type watch interface {
	Watch() (Watcher, error)
}

type KV struct {
	Key   string
	Value []byte
}

type Watcher interface {
	Next() ([]*KV, error)
	Stop() error
}
