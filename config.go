package buffer

import (
	"context"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/go-redis/redis/v8"
)

var (
	defaultFile                = ""
	defaultTTL         int64   = -1
	defaultMemory      float64 = 20
	defaultChannel             = "sync-local"
	defaultRedisExpire         = 48 * time.Hour
)

type (
	config struct {
		// source
		//	souce当前仅支持本地
		//	todo:从oss拉取数据
		source      string
		limit       float64
		hasRedis    bool
		channel     string
		redis       *redis.Client
		tag         string
		redisExpire time.Duration
	}
)

var c = config{
	source:      defaultFile,
	limit:       defaultMemory,
	channel:     defaultChannel,
	tag:         tag(),
	redisExpire: defaultRedisExpire,
}

func tag() string {
	rand.Seed(time.Now().UnixMilli())
	nonce := rand.Int63n(int64(math.Pow(2, 32)))
	return strconv.FormatInt(nonce, 10)
}

func SetRedisExpire(ex time.Duration) {
	c.redisExpire = ex
}

func SetSyncChannel(channel string) {
	c.channel = channel
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
