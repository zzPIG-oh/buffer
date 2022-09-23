package buffer

import (
	"math"
	"time"
	"unsafe"
)

/*
 *	- 循环监听dict
 *		移除ttl过期的data
 *
 */

func patrol() {

	var max, min int64 = 0, 9999

	for k, v := range dict {
		// 有过期时间的
		if v.ttl > int64(-1) && time.Since(time.Unix(v.ttl, 0)) < 1 {
			dictRW.Lock()
			delete(dict, k)
			dictRW.Unlock()
			continue
		}

		// 无过期时间需要根据hot来管控内存
		if v.hot > max {
			max = v.hot
		}

		if v.hot < min {
			min = v.hot
		}

	}

	// 小于20m的内存不做冷热数据删除
	if unsafe.Sizeof(dict) < uintptr((defaultMemory * math.Pow(1024, 3))) {
		return
	}

	middle := (max - min) / 2
	for k, v := range dict {
		// 冷数据删除
		if v.hot < middle {
			dictRW.Lock()
			delete(dict, k)
			dictRW.Unlock()
		}
	}

}
