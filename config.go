package buffer

import (
	"context"
	"log"
	"strings"

	"github.com/go-redis/redis/v8"
)

var (
	defaultFile            = ""
	defaultTTL     int64   = -1
	defaultMemory  float64 = 20
	defaultChannle         = "sync-local"
)

type (
	config struct {
		// source
		//	souce当前仅支持本地
		//	todo:从oss拉取数据
		source   string
		limit    float64
		hasRedis bool
		redis    *redis.Client
	}
)

var c = config{
	source: defaultFile,
	limit:  defaultMemory,
}

func SetMemoryLimit(limit float64) {
	c.limit = limit
}

func SetSource(file string) {
	c.source = file
}

func SetRedis(client *redis.Client) {
	c.redis = client
	val, err := c.redis.Ping(context.Background()).Result()
	if err != nil {
		log.Println("redis.ping error", err.Error())
		return
	}

	if strings.ToLower(val) != "pong" {
		log.Println("redis.instance error", val)
		return
	}

	c.hasRedis = true
}
