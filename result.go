package buffer

import "strconv"

type result struct{ result interface{} }

func (r *result) ToInt() int {
	switch val := r.result.(type) {
	case int:
		return val
	case string:
		i, e := strconv.ParseInt(val, 10, 64)
		if e != nil {
			return 0
		}
		return int(i)
	default:
		return 0
	}
}

func (r *result) ToBool() bool {
	switch val := r.result.(type) {
	case bool:
		return val
	case string:
		i, e := strconv.ParseBool(val)
		if e != nil {
			return false
		}
		return i
	default:
		return false
	}
}

// IsEmpty
//	true:不存在所要搜索的kv
//	false:存在需要的kv
func (r *result) IsEmpty() bool {
	return r.result == nil
}

func (r *result) ToString() string {
	switch val := r.result.(type) {
	case string:
		return val
	default:
		return ""
	}
}

func (r *result) Interface() interface{} {
	return r.result
}
