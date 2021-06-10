package cache

import (
	"github.com/patrickmn/go-cache"
	"time"
)

const DEFAULT_EXPIRES_IN = 5
const DEFAULT_PURGE_EXPIRES_IN_PERIOD = 10

type FileCache struct {
	*cache.Cache
	cacheFile string
}

func NewFileCache(namespace string, defaultLifeTime time.Duration, directory string) *FileCache {

	if defaultLifeTime <= 0 {
		defaultLifeTime = time.Duration(DEFAULT_EXPIRES_IN) * time.Minute
	}
	defaultPurgePeriod := time.Duration(DEFAULT_PURGE_EXPIRES_IN_PERIOD) * time.Minute

	fileCache := &FileCache{
		cache.NewFrom(
			defaultLifeTime,
			defaultPurgePeriod,
			nil,
		),
		directory,
	}
	return fileCache
}

func (cache *FileCache) Get(key string, defaultValue interface{}) error {
	return nil
}

func (cache *FileCache) Set(key string, value interface{}, expires time.Duration) error {
	return nil
}

func (cache *FileCache) Has(key string) bool {
	return false
}

func (cache *FileCache) AddNX(key string, value interface{}, ttl time.Duration) bool {
	return false
}

func (cache *FileCache) Add(key string, value interface{}, ttl time.Duration) (err error) {
	return nil
}

func (cache *FileCache) Remember(key string, ttl time.Duration, callback func() interface{}) (obj interface{}, err error) {
	return nil, nil
}
