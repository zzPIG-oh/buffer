package buffer

type Buffer interface {
	// map 存取
	//	key、field应当为string；value可以是任意值
	Hset(key, field string, value interface{})
	// Hget
	//	key、field应当为string
	Hget(key, field string) Result

	// Hdel
	//	key、field应当为string
	Hdel(key, field string)

	// 判断健存不存在
	Exist(key string) bool
}

type Result interface {
	IsEmpty() bool
	// 结果转换
	//	ToInt 将interface转换成int
	ToInt() int
	//	ToBool 将interface转换成bool
	ToBool() bool
	//	ToString 将interface转换成string
	ToString() string

	//	返回原结果
	Interface() interface{}
}

var (
	_ Buffer = &buffer{}
	_ Result = &result{}
)
