package buffer

import (
	"context"
	"strings"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-redis/redis/v8"
)

var (
	defaultFile           = "/opt/buffer.json"
	defaultTTL    int64   = -1
	defaultMemory float64 = 20
	// defaultFile           = "buffer.json"
)

type (
	config struct {
		// source
		//	souce当前仅支持本地
		//	todo:从oss拉取数据
		source   string
		hasRedis bool
		redis    *redis.Client
	}
)

var c = config{
	source: defaultFile,
}

func SetSource(file string) {
	c.source = file
}

func SetRedis(client *redis.Client) {
	c.redis = client
	val, err := c.redis.Ping(context.Background()).Result()
	if err != nil {
		log.DefaultLogger.Log(log.LevelFatal, "redis.ping error", err.Error())
		return
	}
	if strings.ToLower(val) != "pong" {
		log.DefaultLogger.Log(log.LevelFatal, "redis.instance error", val)
		return
	}
	c.hasRedis = true
}
