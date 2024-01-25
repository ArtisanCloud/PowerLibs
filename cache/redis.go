package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	fmt2 "github.com/ArtisanCloud/PowerLibs/v3/fmt"
	"github.com/ArtisanCloud/PowerLibs/v3/object"
	"github.com/redis/go-redis/v9"
)

type clienter interface {
	Get(ctx context.Context, key string) *redis.StringCmd
	MGet(ctx context.Context, keys ...string) *redis.SliceCmd
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.StatusCmd
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) *redis.BoolCmd
	Expire(ctx context.Context, key string, expiration time.Duration) *redis.BoolCmd
	Exists(ctx context.Context, keys ...string) *redis.IntCmd
	TxPipeline() redis.Pipeliner
	Del(ctx context.Context, keys ...string) *redis.IntCmd
	Close() error
}

type GRedis struct {
	cli               clienter
	defaultExpiration time.Duration
	lockRetries       int
}

const SYSTEM_CACHE_TIMEOUT = 60 * 60
const SYSTEM_CACHE_TIMEOUT_MINUTE = 60
const SYSTEM_CACHE_TIMEOUT_HOUR = 60 * 60
const SYSTEM_CACHE_TIMEOUT_DAY = 60 * 60 * 24
const SYSTEM_CACHE_TIMEOUT_MONTH = 60 * 60 * 24 * 30
const SYSTEM_CACHE_TIMEOUT_SEASON = 60 * 60 * 24 * 30 * 3
const SYSTEM_CACHE_TIMEOUT_YEAR = 60 * 60 * 24 * 30 * 3 * 12

const (
	defaultMaxIdle            = 5
	defaultMaxActive          = 0
	defaultTimeoutIdle        = 240
	defaultTimeoutConnect     = 10000
	defaultTimeoutRead        = 5000
	defaultTimeoutWrite       = 5000
	defaultAddr               = "localhost:6379"
	defaultProtocol           = "tcp"
	defaultRetryThreshold     = 5
	defaultIdleCheckFrequency = time.Minute
)

const SCRIPT_SETEX = `return redis.call('exists',KEYS[1])<1 and redis.call('setex',KEYS[1],ARGV[2],ARGV[1])`

var CTXRedis = context.Background()

const lockRetries = 5

func NewGRedis(opts *redis.Options) (gr *GRedis) {
	c := redis.NewClient(opts)
	gr = &GRedis{
		cli:         c,
		lockRetries: lockRetries,
	}

	return gr
}

func NewGRedisCluster(opts *redis.ClusterOptions) (gr *GRedis) {
	c := redis.NewClusterClient(opts)
	gr = &GRedis{
		cli:         c,
		lockRetries: lockRetries,
	}

	return gr
}

func (gr *GRedis) AddNX(key string, value interface{}, ttl time.Duration) bool {
	cmd := gr.cli.SetNX(CTXRedis, key, value, ttl)
	r, err := cmd.Result()
	//fmt.Printf("r:%b \r\n", r)
	if err != nil {
		fmt.Errorf("SetNX error: %+v \r\n", err)
	}
	return r
}
func (gr *GRedis) Add(key string, value interface{}, ttl time.Duration) (err error) {
	// If the store has an "add" method we will call the method on the store so it
	// has a chance to override this logic. Some drivers better support the way
	// this operation should work with a total "atomic" implementation of it.
	var obj interface{}
	obj, err = gr.Get(key, obj)
	if err == ErrCacheMiss {
		return gr.SetEx(key, value, ttl)
	} else {
		return errors.New("this value has been actually added to the cache")
	}

}

func (gr *GRedis) Set(key string, value interface{}, expires time.Duration) error {
	mValue, err := json.Marshal(value)
	//mExpire, err := json.Marshal(expires)
	if err != nil {
		return err
	}

	result := gr.cli.Set(CTXRedis, key, mValue, expires)
	return result.Err()
}

func (gr *GRedis) SetEx(key string, value interface{}, expires time.Duration) error {
	mValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	// 使用Set()方法并传递SetOptions设置过期时间
	cmd := gr.cli.Set(CTXRedis, key, mValue, expires)
	return cmd.Err()
}

func (gr *GRedis) Get(key string, ptrValue interface{}) (returnValue interface{}, err error) {
	b, err := gr.cli.Get(CTXRedis, key).Bytes()
	if err == redis.Nil {
		return nil, ErrCacheMiss
	}
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &ptrValue)
	returnValue = ptrValue
	return returnValue, err
}

func (gr *GRedis) Has(key string) bool {

	value, err := gr.Get(key, nil)
	if value != nil && err == nil {
		return true
	}

	return false
}

func (gr *GRedis) GetMulti(keys ...string) (object.HashMap, error) {
	res, err := gr.cli.MGet(CTXRedis, keys...).Result()
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
	return gr.cli.Del(CTXRedis, key).Err()
}

func (gr *GRedis) Keys() ([]string, error) {
	// if client is a redis.ClusterClient
	if cli, ok := gr.cli.(*redis.ClusterClient); ok {
		var keys []string
		slots, err := cli.ClusterSlots(CTXRedis).Result()
		if err != nil {
			return nil, err
		}
		for _, slot := range slots {
			for _, node := range slot.Nodes {
				client := redis.NewClient(&redis.Options{
					Addr: node.Addr,
				})
				nodeKeys, err := client.Keys(CTXRedis, "*").Result()
				if err != nil {
					return nil, err
				}
				keys = append(keys, nodeKeys...)
				client.Close()
			}
		}
		return keys, nil
	} else if cli, ok := gr.cli.(*redis.Client); ok {
		// code to be executed if client is a redis.Client
		return cli.Keys(CTXRedis, "*").Result()
	}

	return nil, nil
}

func (gr *GRedis) Flush() error {
	if cli, ok := gr.cli.(*redis.ClusterClient); ok {
		// 在集群环境下，FlushAll命令需要对每个主节点分别执行
		slots, err := cli.ClusterSlots(CTXRedis).Result()
		if err != nil {
			return err
		}
		for _, slot := range slots {
			for _, node := range slot.Nodes {
				client := redis.NewClient(&redis.Options{
					Addr: node.Addr,
				})
				if err := client.FlushAll(CTXRedis).Err(); err != nil {
					return err
				}
				client.Close()
			}
		}
		return nil
	} else if cli, ok := gr.cli.(*redis.Client); ok {
		// code to be executed if client is a redis.Client
		return cli.FlushAll(CTXRedis).Err()
	}

	return nil
}

/**
 * Get an item from the cache, or execute the given Closure and store the result.
 *
 * @param  string  key
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @param  \Closure  callback
 * @return mixed
 */
func (gr *GRedis) Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (obj interface{}, err error) {

	var value interface{}
	value, err = gr.Get(key, value)

	// If the item exists in the cache we will just return this immediately and if
	// not we will execute the given Closure and cache the result of that for a
	// given number of seconds so it's available for all subsequent requests.
	if err != nil && err != ErrCacheMiss {
		fmt2.Dump("error:", err.Error())
		return nil, err

	} else if value != nil {
		return value, err
	}

	value, err = callback()
	if err != nil {
		return nil, err
	}

	result := gr.Put(key, value, ttl)
	if !result {
		err = errors.New(fmt.Sprintf("remember cache put err, ttl:%d", ttl))
	}
	// ErrCacheMiss and query value from source
	return value, err
}

/**
 * Store an item in the cache.
 *
 * @param  string  key
 * @param  mixed  value
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) Put(key interface{}, value interface{}, ttl time.Duration) bool {
	// key如果是数组
	//if arrayKey, ok := key.([]interface{}) !ok {
	//	return gr.PutMany(arrayKey, value)
	//}

	//if ttl == nil {
	//	return gr.forever(key, value)
	//}

	//seconds := gr.GetSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.Delete(key)
	//}

	//result = gr.cli.Put(gr.itemKey(key), value, seconds)

	err := gr.SetEx(key.(string), value, ttl)
	if err != nil {
		fmt.Println(err)
		return false
	}

	return true

}

/**
 * Store multiple items in the cache for a given number of seconds.
 *
 * @param  array  values
 * @param  \DateTimeInterface|\DateInterval|int|null  ttl
 * @return bool
 */
func (gr *GRedis) PutMany(values object.Array, ttl time.Duration) bool {
	//if ttl == nil {
	//	return gr.PutManyForever(values)
	//}
	//
	//seconds := gr.GetSeconds(ttl)
	//
	//if seconds <= 0 {
	//	return gr.cli.Del(array_keys(values))
	//}
	//
	//gr.cli.

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
func (gr *GRedis) GetSeconds(ttl time.Duration) int {
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

func (gr *GRedis) SetByTags(key string, val interface{}, tags []string, expiry time.Duration) error {
	pipe := gr.cli.TxPipeline()
	for _, tag := range tags {
		pipe.SAdd(CTXRedis, tag, key)
		pipe.Expire(CTXRedis, tag, expiry)
	}

	pipe.Set(CTXRedis, key, val, expiry)

	_, errExec := pipe.Exec(CTXRedis)
	return errExec
}

func (gr *GRedis) Invalidate(tags []string) error {
	if cli, ok := gr.cli.(*redis.ClusterClient); ok {
		// Redis集群模式下，删除操作需要确保键在同一个节点(slot)上。
		for _, tag := range tags {
			// 首先获取集群中，tag对应的所有成员
			keys, err := cli.SMembers(CTXRedis, tag).Result()
			if err != nil {
				return err
			}
			keys = append(keys, tag) // 也需要删除tag键本身

			// 在集群模式下，删除操作需要单独对每个键进行处理
			for _, key := range keys {
				_, err := cli.Del(CTXRedis, key).Result()
				if err != nil {
					return err
				}
			}
		}
		return nil
	} else if cli, ok := gr.cli.(*redis.Client); ok {
		keys := make([]string, 0)
		for _, tag := range tags {
			k, _ := cli.SMembers(CTXRedis, tag).Result()
			keys = append(keys, tag)
			keys = append(keys, k...)
		}
		cli.Del(CTXRedis, keys...)
	}

	return nil
}
