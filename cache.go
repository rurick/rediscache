// Copyright 2021 (c) Yuriy Iovkov aka Rurick.
// iovkov@antsgames.com

// access to redis cache.
// params for connection to redis DB tacking from OS ENV:
// REDIS_User
// REDIS_Passwd
// REDIS_Addr
// REDIS_DB

// default value of cache life time is 1 minute. You can change it by calling SetCacheLifeTime

package cache

import (
	"context"
	"crypto/sha1"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	logger "github.com/sirupsen/logrus"
)

type (
	Opt struct {
		Name string
		Val  interface{}
	}

	cacheData struct {
		v interface{}
	}

	// configuration for database connection - list of addresses
	configuration map[string]string
)

var (
	c               *cache.Cache
	ctx             context.Context
	cacheExpiration time.Duration
)

func init() {
	ring := redis.NewRing(&redis.RingOptions{Addrs: getConfiguration()})

	c = cache.New(&cache.Options{
		Redis: ring,
		// LocalCache: cache.NewTinyLFU(),
	})
	ctx = context.Background()
	cacheExpiration = time.Minute // set default value of cache life time
}

// return configuration for database connection
func getConfiguration() configuration {
	_ = godotenv.Load() // Try to load env variables from .env file if one exists
	osEnvAddr := os.Getenv("REDIS_Addr")
	lst := strings.Split(osEnvAddr, ",")

	if len(lst) == 0 || osEnvAddr == "" {
		logger.WithFields(logger.Fields{
			"project":       "admin",
			"package":       "pkg/cache",
			"func":          "getConfiguration",
			"configuration": lst,
		}).Warning("It's possible configuration for redis not set")

	}
	ring := configuration{}
	for i, l := range lst {
		ring[fmt.Sprintf("server%d", i)] = l
	}
	return ring
}

// SetCacheExpiration - set default cache expire time
func SetCacheExpiration(d time.Duration) {
	if d > 0 {
		cacheExpiration = d
	}
}

// KeyGen generate 20bits hex key for v
func KeyGen(v ...interface{}) string {
	s := sha1.Sum([]byte(fmt.Sprint(v...)))
	return hex.EncodeToString(s[:])

}

func get(key string) (out interface{}, exists bool, err error) {
	err = c.Get(ctx, key, &out)
	exists = err != cache.ErrCacheMiss
	return
}

// Load object from cache to obj
func Load(key string, obj interface{}) (exists bool, err error) {
	err = c.Get(ctx, key, &obj)
	exists = err != cache.ErrCacheMiss
	return
}

// Set - set value in cache
// cache life time by default is cacheExpiration
// for set up cache life time for setting value set up opt parameter (exp: Opt{"expiration", time.Minute})
func Set(key string, val interface{}, opt ...Opt) error {
	exp := expirationFromOpt(opt...)

	return c.Set(&cache.Item{
		Ctx:   ctx,
		Key:   key,
		Value: &val,
		TTL:   exp,
	})
}

// GetInt64 - get value from cache as int64
func GetInt64(key string) (int64, bool, error) {
	r, ok, err := get(key)
	if err != nil {
		return 0, ok, err
	}

	var out int64
	switch r.(type) {
	case uint8:
		out = int64(r.(uint8))
	case int8:
		out = int64(r.(int8))
	case uint16:
		out = int64(r.(uint16))
	case int16:
		out = int64(r.(int16))
	case uint32:
		out = int64(r.(uint32))
	case int32:
		out = int64(r.(int32))
	case uint64:
		out = int64(r.(uint64))
	case int64:
		out = r.(int64)
	default:
		return 0, ok, errors.New("Value is not integer type")
	}

	return out, ok, nil
}

// Delete - delete value from cache by key
func Delete(key string) error {
	return c.Delete(ctx, key)
}

// GetString - get value from cache as string
func GetString(key string) (string, bool, error) {
	r, ok, err := get(key)
	if err != nil {
		return "", ok, err
	}
	return r.(string), ok, nil
}

func expirationFromOpt(opt ...Opt) time.Duration {
	if len(opt) > 0 {
		for _, o := range opt {
			if o.Name == "expiration" {
				return o.Val.(time.Duration)
			}
		}
	}
	return cacheExpiration
}
