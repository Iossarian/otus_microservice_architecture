package build

import "github.com/bradfitz/gomemcache/memcache"

func (b *Builder) memcache() (*memcache.Client, error) {
	return memcache.New(b.config.Memcache.Host), nil
}
