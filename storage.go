package buffer

import (
	"context"
	"log"
	"sync"
	"sync/atomic"
	"time"
)

// - inner
// 	本地缓存k:v；v的具体结构
type inner struct {
	// -
	// 	该键内的值
	hash map[string]interface{}

	// -
	// 	time to live 该键可以存活的时间
	// 	仅支持秒
	//	这里有一个special的存在 即-1;使用-1后,键的删除根据当前hot值
	ttl int64

	// -
	// 	hot是系统使用该键下的频率统计
	//	每三分钟做一次数据热度排行;减少本地缓存的数据、减轻内存压力
	hot int64

	// -
	// 	控制hash的并发锁
	rw sync.RWMutex
}

// dict
//	本地缓存字典
var (
	dict   = map[string]*inner{}
	dictRW = sync.RWMutex{}
)

type buffer struct{}

// empty
//	允许插入一条value为struct的map
var Empty = struct{}{}

func NewBufferClient() Buffer {
	return &buffer{}
}

func (b *buffer) Hset(key, field string, value interface{}) {
	// 任何非0的数字与运算后都大于0;
	//	不会报错;直接返回
	if len(key) < 1 || len(field) < 1 {
		return
	}

	dictRW.Lock()
	kv, ok := dict[key]
	if !ok {
		dict[key] = &inner{
			hash: map[string]interface{}{
				field: value,
			},
			rw:  sync.RWMutex{},
			ttl: defaultTTL,
		}
		dictRW.Unlock()
		return
	}
	dictRW.Unlock()

	kv.rw.Lock()
	kv.hash[field] = value
	kv.rw.Unlock()

	if c.hasRedis {
		c.redis.HSet(context.Background(), key, field, value)
	}

}

func (b *buffer) Hget(key, field string) (r Result) {

	defer func() {

		if r.IsEmpty() && c.hasRedis {

			val, err := c.redis.HGet(context.Background(), key, field).Result()
			if err != nil {
				log.Println("Hget.error", err.Error())
				return
			}

			r = &result{result: val}
			b.Hset(key, field, val)

		}

		if !r.IsEmpty() {
			// 每取一次增加一次热度
			atomic.AddInt64(&dict[key].hot, 1)
		}

	}()

	// 加入读锁--防止此时并发出现map的更改
	dictRW.RLock()
	kv, ok := dict[key]
	if !ok {
		dictRW.RUnlock()
		return &result{result: nil}
	}
	dictRW.RUnlock()

	if kv.ttl > -1 && time.Since(time.Unix(kv.ttl, 0)) < 1 {
		return &result{result: nil}
	}

	// 不考虑惰性删除--惰性删除此处还要加锁处理
	kv.rw.RLock()
	ptr := &result{result: kv.hash[field]}
	kv.rw.RUnlock()

	return ptr
}

func (b *buffer) Hdel(key, field string) {

	dictRW.RLock()
	kv, ok := dict[key]
	if !ok {
		return
	}
	dictRW.RUnlock()

	kv.rw.Lock()
	delete(kv.hash, field)
	kv.rw.Unlock()

	if c.hasRedis {
		c.redis.HDel(context.Background(), key, field)
	}

}

func (b *buffer) Exist(key string) bool {
	_, ok := dict[key]
	return ok
}
