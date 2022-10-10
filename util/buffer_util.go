package util

import "strings"

func String(str ...string) string {
	return strings.Join(str, "%")
}

func Spilt(old string) (key, field string) {
	r := strings.Split(old, "%")
	return r[0], r[1]
}
