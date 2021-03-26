package cache

import (
	"encoding/json"
	"github.com/ArtisanCloud/go-libs/object"
	"github.com/go-redis/redis"
	"time"
)

type GRedis struct {
	Pool              *redis.Client
	defaultExpiration time.Duration
	lockRetries       int
}

const (
	defaultMaxIdle        = 5
	defaultMaxActive      = 0
	defaultTimeoutIdle    = 240
	defaultTimeoutConnect = 10000
	defaultTimeoutRead    = 5000
	defaultTimeoutWrite   = 5000
	defaultHost           = "localhost:6379"
	defaultProtocol       = "tcp"
	defaultRetryThreshold = 5
)

type RedisOptions struct {
	MaxIdle        int
	MaxActive      int
	Protocol       string
	Host           string
	Password       string
	DB             int
	Expiration     time.Duration
	TimeoutConnect int
	TimeoutRead    int
	TimeoutWrite   int
	TimeoutIdle    int
}

const lockRetries = 5

func NewGRedis(opts interface{}) (gr *GRedis) {

	if options, ok := opts.(*RedisOptions); ok {
		options = options.initDefaults()
		toD := time.Millisecond * time.Duration(options.TimeoutConnect)
		toR := time.Millisecond * time.Duration(options.TimeoutRead)
		toW := time.Millisecond * time.Duration(options.TimeoutWrite)
		toI := time.Duration(options.TimeoutIdle) * time.Second
		option := &redis.Options{
			Addr:               options.Host,
			DB:                 options.DB,
			DialTimeout:        toD,
			ReadTimeout:        toR,
			WriteTimeout:       toW,
			PoolSize:           options.MaxActive,
			PoolTimeout:        30 * time.Second,
			IdleTimeout:        toI,
			Password:           options.Password,
			IdleCheckFrequency: 500 * time.Millisecond,
		}

		c := redis.NewClient(option)
		gr = &GRedis{
			Pool:        c,
			lockRetries: lockRetries,
		}
		return gr

	}

	return gr

}

func (r *RedisOptions) initDefaults() *RedisOptions {
	if r.MaxIdle == 0 {
		r.MaxIdle = defaultMaxIdle
	}

	if r.MaxActive == 0 {
		r.MaxActive = defaultMaxActive
	}

	if r.TimeoutConnect == 0 {
		r.TimeoutConnect = defaultTimeoutConnect
	}

	if r.TimeoutIdle == 0 {
		r.TimeoutIdle = defaultTimeoutIdle
	}

	if r.TimeoutRead == 0 {
		r.TimeoutRead = defaultTimeoutRead
	}

	if r.TimeoutWrite == 0 {
		r.TimeoutWrite = defaultTimeoutWrite
	}

	if r.Host == "" {
		r.Host = defaultHost
	}

	if r.Protocol == "" {
		r.Protocol = defaultProtocol
	}

	return r
}

func (gr *GRedis) Set(key string, value interface{}, expires time.Duration) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return gr.Pool.Set(key, b, expires).Err()
}

func (gr *GRedis) Get(key string, ptrValue interface{}) error {
	b, err := gr.Pool.Get(key).Bytes()
	if err == redis.Nil {
		return ErrCacheMiss
	}

	if err != nil {
		return err
	}

	return json.Unmarshal(b, ptrValue)
}

func (gr *GRedis) GetMulti(keys ...string) (object.HashMap, error) {
	res, err := gr.Pool.MGet(keys...).Result()
	if err != nil {
		return nil, err
	}

	if len(res) == 0 {
		return nil, ErrCacheMiss
	}

	m := make(object.HashMap)
	for ix, key := range keys {
		m[key] = res[ix].(string)
	}
	return m, nil
}

func (gr *GRedis) Delete(key string) error {
	return gr.Pool.Del(key).Err()
}

func (gr *GRedis) Keys() ([]string, error) {
	return gr.Pool.Keys("*").Result()
}

func (gr *GRedis) Flush() error {
	return gr.Pool.FlushAll().Err()
}

/**
 * Get an item from the cache, or execute the given Closure and store the result.
 *
 * @param  string  key
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @param  \Closure  callback
 * @return mixed
 */
func (gr *GRedis) Remember(key string, ttl interface{}, callback func() interface{}) interface{} {

	var value interface{}
	err := gr.Get(key, &value)

	// If the item exists in the cache we will just return this immediately and if
	// not we will execute the given Closure and cache the result of that for a
	// given number of seconds so it's available for all subsequent requests.
	if err != nil {
		return value
	}
	value = callback()
	gr.Put(key, value, ttl)

	return value
}

/**
 * Store an item in the cache.
 *
 * @param  string  key
 * @param  mixed  value
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) Put(key interface{}, value interface{}, ttl interface{}) bool {
	// key如果是数组
	if arrayKey, ok := key.([]interface{}); !ok {
		return gr.PutMany(arrayKey, value)
	}
	//
	//if ttl == nil {
	//	return gr.forever(key, value)
	//}
	//
	//seconds = gr.getSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.forget(key)
	//}
	//
	//result = gr.store- > put(gr.itemKey(key), value, seconds)
	//
	//if result {
	//	gr.event(new
	//	KeyWritten(key, value, seconds))
	//}

	//return result
	return false
}

/**
 * Store multiple items in the cache for a given number of seconds.
 *
 * @param  array  values
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) PutMany(values object.Array, ttl interface{}) bool {
	//if ttl == nil {
	//	return gr.PutManyForever(values)
	//}
	//
	//seconds := gr.GetSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.Pool.Del(array_keys(values))
	//}
	//
	//gr.Pool.

	return false
}

/**
 * Store multiple items in the cache indefinitely.
 *
 * @param  array  values
 * @return bool
 */
func (gr *GRedis) PutManyForever(values []interface{}) bool {
	result := true

	//for key, value := range values {
	//
	//	if !gr.Forever(key, value) {
	//		result = false
	//	}
	//}

	return result
}

/**
 * Calculate the number of seconds for the given TTL.
 *
 * @param  \DateTimeInterface|\DateInterval|int  ttl
 * @return int
 */
func (gr *GRedis) GetSeconds(ttl interface{}) int {
	//duration := gr.ParseDateInterval(ttl)
	//
	//if reflect.Type(duration).Kind() == DateTimeInterface {
	//	duration = carbon.Now().diffInRealSeconds(duration, false)
	//}
	//
	//if duration > 0 {
	//	return duration
	//} else {
	//	return 0
	//}
	return 0
}
