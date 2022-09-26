package buffer

import (
	"encoding/json"
	"io"
	"os"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
)

var (
	once sync.Once
)

// Init
//	加载指定的数据源
func init() {
	once.Do(func() {
		// maybe don`t need source
		if defaultFile == "" {
			return
		}

		fd, err := os.Open(defaultFile)
		if err != nil {
			log.DefaultLogger.Log(log.LevelFatal, "source.open error", err.Error())
			return
		}

		bt, err := io.ReadAll(fd)
		if err != nil {
			log.DefaultLogger.Log(log.LevelFatal, "source.read error", err.Error())
			return
		}

		json.Unmarshal(bt, &dict)
	})

	go func() {
		// -
		// 	每三十秒巡查一次
		//	不用担心键过期不能被及时回收,get的时候会判断触发
		for range time.Tick(30 * time.Second) {
			patrol()
		}
	}()
}
