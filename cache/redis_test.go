package cache

import (
	"github.com/redis/go-redis/v9"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var CacheConnection *GRedis

func Test_Remember(t *testing.T) {

	redisCache := getTestGRedis()

	data := map[string]interface{}{
		"string": "value",
		"int":    10,
		"bool":   false,
	}

	//cachedData, err := redisCache.Remember("recordType:Membership__c.Period", SYSTEM_CACHE_TIMEOUT_MONTH*time.Second, func() interface{} {
	cachedData, err := redisCache.Remember("test.redis", SYSTEM_CACHE_TIMEOUT_MONTH*time.Second, func() (interface{}, error) {
		return data, nil
	})

	if !assert.ObjectsAreEqual(data, cachedData) {
		t.Error(err)
	}

}

func getTestGRedis() *GRedis {

	if CacheConnection != nil {
		return CacheConnection
	}

	options := redis.Options{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       1,
	}

	CacheConnection = NewGRedis(&options)

	return CacheConnection
}
