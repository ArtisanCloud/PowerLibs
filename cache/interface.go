package cache

import (
	"time"
)

var ACCache CacheInterface

type CacheInterface interface {

	//SetOptions(opts interface{}) error
	//
	///**
	//* Fetches a value from the cache.
	//*
	//* @param string key     The unique key of this item in the cache.
	//* @param mixed  default Default value to return if the key does not exist.
	//*
	//* @return mixed The value of the item from the cache, or default in case of cache miss.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if the key string,is not a legal value.
	//*/
	Get(key string, defaultValue interface{}) (ptrValue interface{}, err error)
	//
	///**
	//* Persists data in the cache, uniquely referenced by a key with an optional expiration TTL time.
	//*
	//* @param string                 key   The key of the item to store.
	//* @param mixed                  value The value of the item to store, must be serializable.
	//* @param null|int|\DateInterval ttl   Optional. The TTL value of this item. If no value is sent and
	//*                                      the driver supports TTL then the library may set a default value
	//*                                      for it or let the driver take care of that.
	//*
	//* @return bool True on success and false on failure.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if the key string,is not a legal value.
	//*/
	Set(key string, value interface{}, expires time.Duration) error
	//
	///**
	//* Delete an item from the cache by its unique key.
	//*
	//* @param string key The unique cache key of the item to delete.
	//*
	//* @return bool True if the item was successfully removed. False if there was an error.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if the key string,is not a legal value.
	//*/
	//Delete(key string) bool
	//
	///**
	//* Wipes clean the entire cache's keys.
	//*
	//* @return bool True on success and false on failure.
	//*/
	//Clear() bool
	//
	///**
	//* Obtains multiple cache items by their unique keys.
	//*
	//* @param iterable keys    A list of keys that can obtained in a single operation.
	//* @param mixed    default Default value to return for keys that do not exist.
	//*
	//* @return iterable A list of key => value pairs. Cache keys that do not exist or are stale will have default as value.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if keys is neither an array nor a Traversable,
	//*   or if any of the keys are not a legal value.
	//*/
	//GetMultiple(keys string, defaultValue object.HashMap) interface{}
	//
	///**
	//* Persists a set of key => value pairs in the cache, with an optional TTL.
	//*
	//* @param iterable               values A list of key => value pairs for a multiple-set operation.
	//* @param null|int|\DateInterval ttl    Optional. The TTL value of this item. If no value is sent and
	//*                                       the driver supports TTL then the library may set a default value
	//*                                       for it or let the driver take care of that.
	//*
	//* @return bool True on success and false on failure.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if values is neither an array nor a Traversable,
	//*   or if any of the values are not a legal value.
	//*/
	//SetMultiple(values string, ttl int) bool
	//
	///**
	//* Deletes multiple cache items in a single operation.
	//*
	//* @param iterable keys A list of string-based keys to be deleted.
	//*
	//* @return bool True if the items were successfully removed. False if there was an error.
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if keys is neither an array nor a Traversable,
	//*   or if any of the keys are not a legal value.
	//*/
	//DeleteMultiple(keys []string) bool
	//
	///**
	//* Determines whether an item is present in the cache.
	//*
	//* NOTE: It is recommended that has() is only to be used for cache warming type purposes
	//* and not to be used within your live applications operations for get/set, as this method
	//* is subject to a race condition where your has() will return true and immediately after,
	//* another script can remove it making the state of your app out of date.
	//*
	//* @param string key The cache item key.
	//*
	//* @return bool
	//*
	//* @panic InvalidArgumentException
	//*   MUST be thrown if the key string,is not a legal value.
	//*/
	Has(key string) bool
	//
	///**
	//* Retrieve an item from the cache and delete it.
	//*
	//* @param  string  key
	//* @param  mixed  default
	//* @return mixed
	//*/
	//Pull(key string, defaultValue object.HashMap) interface{}
	//
	///**
	//* Store an item in the cache.
	//*
	//* @param  string  key
	//* @param  mixed  value
	//* @param  \DateTimeInterface|\DateInterval|int|null  ttl
	//* @return bool
	//*/
	//Put(key string, value, ttl int) bool
	//

	AddNX(key string, value interface{}, ttl time.Duration) bool

	///**
	//* Store an item in the cache if the key does not exist.
	//*
	//* @param  string  key
	//* @param  mixed  value
	//* @param  \DateTimeInterface|\DateInterval|int|null  ttl
	//* @return bool
	//*/
	Add(key string, value interface{}, ttl time.Duration) (err error)
	//
	///**
	//* Increment the value of an item in the cache.
	//*
	//* @param  string  key
	//* @param  mixed  value
	//* @return int|bool
	//*/
	//Increment(key string, value object.HashMap)(int, bool)
	//
	///**
	//* Decrement the value of an item in the cache.
	//*
	//* @param  string  key
	//* @param  mixed  value
	//* @return int|bool
	//*/
	//Decrement(key string, value object.HashMap) (int, bool)
	//
	///**
	//* Store an item in the cache indefinitely.
	//*
	//* @param  string  key
	//* @param  mixed  value
	//* @return bool
	//*/
	//Forever(key string, value object.HashMap) bool
	//
	///**
	//* Get an item from the cache, or execute the given Closure and store the result.
	//*
	//* @param  string  key
	//* @param  \DateTimeInterface|\DateInterval|int|null  ttl
	//* @param  \Closure  callback
	//* @return mixed
	//*/
	Remember(key string, ttl time.Duration, callback func() (interface{}, error)) (obj interface{}, err error)
	//
	///**
	//* Get an item from the cache, or execute the given Closure and store the result forever.
	//*
	//* @param  string  key
	//* @param  \Closure  callback
	//* @return mixed
	//*/
	//Sear(key string, callback func()) interface{}
	//
	///**
	//* Get an item from the cache, or execute the given Closure and store the result forever.
	//*
	//* @param  string  key
	//* @param  \Closure  callback
	//* @return mixed
	//*/
	//RememberForever(key string, callback func()) interface{}
	//
	///**
	//* Remove an item from the cache.
	//*
	//* @param  string  key
	//* @return bool
	//*/
	//Forget(key string) bool
	//
	///**
	//* Get the cache store implementation.
	//*
	//* @return \Illuminate\Contracts\Cache\Store
	//*/
	//GetStore()

}
