package util

import "strings"

func String(str ...string) string {
	return strings.Join(str, "%")
}

func Spilt(old string) []string {
	return strings.Split(old, "%")
}
