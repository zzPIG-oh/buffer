package buffer

import (
	"context"
	"encoding/json"
	"fast-buffer/util"
	"io"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

var (
	once sync.Once
)

// Start
//	加载指定的数据源
func Start() {

	once.Do(func() {

		// 加载数据源
		bt, err := read(defaultFile)
		if err != nil {
			return
		}
		write(bt)

		// 初始化http
		if len(c.addr) > 0 {
			http.HandleFunc("/fast-buffer/probe", probe)
			http.ListenAndServe(c.addr, nil)
		}

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

func read(file string) ([]byte, error) {

	fd, err := os.Open(defaultFile)
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
	for msg := range c.redis.PSubscribe(context.Background(), defaultChannle).Channel() {
		k, f := util.Spilt(msg.Payload)
		b.Hdel(k, f)
	}
}
