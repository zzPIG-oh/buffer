package buffer

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/zzPIG-oh/buffer/util"
)

var (
	once sync.Once
)

// Start
//	加载指定的数据源
func Start() {

	once.Do(func() {

		// 加载数据源
		bt, err := read()
		if err != nil {
			return
		}
		write(bt)
	})

	go func() {
		// -
		// 	每三十秒巡查一次
		//	不用担心键过期不能被及时回收,get的时候会判断触发
		for range time.Tick(30 * time.Second) {
			patrol()
		}
	}()

	go syncBuffer()
}

func read() ([]byte, error) {

	fd, err := os.Open(c.source)
	if err != nil {
		log.Println("source.open error", err.Error())
		return nil, err
	}

	bt, err := io.ReadAll(fd)
	if err != nil {
		log.Println("source.read error", err.Error())
		return nil, err
	}

	return bt, err
}

func write(bt []byte) {

	// tmpMap--临时接受dict的值
	// 数据源的格式应该为map[string]map[string]interface{}
	tmpMap := map[string]map[string]interface{}{}

	err := json.Unmarshal(bt, &tmpMap)
	if err != nil {
		log.Println("json.Unmarshal error", err.Error())
		return
	}

	for k, v := range tmpMap {
		dict[k] = &inner{
			hash: v,
			ttl:  -1,
			hot:  9999, // 数据源的数据暂时全部不过期
		}
	}

}

func syncBuffer() {
	if !c.hasRedis {
		return
	}
	b := &buffer{}
	for msg := range c.redis.PSubscribe(context.Background(), c.channel).Channel() {

		log.Println("message:", msg.Payload)

		args := util.Spilt(msg.Payload)

		if len(args) < 3 {
			continue
		}

		// 不处理自己的消息
		if args[2] == c.tag {
			continue
		}

		b.Hdel(args[0], args[1])
	}
}
