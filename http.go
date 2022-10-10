package buffer

import (
	"encoding/json"
	"net/http"
)

/**
 * 使用http读取本地缓存的数据
 */

func probe(w http.ResponseWriter, r *http.Request) {
	data, _ := json.Marshal(dict)
	w.Write(data)
}
