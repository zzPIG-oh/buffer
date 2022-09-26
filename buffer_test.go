package buffer

import (
	"fmt"
	"testing"
)

func TestBuffer(t *testing.T) {
	bc := NewBufferClient()
	bl := bc.Hget("keyword", "show").IsEmpty()
	rt := bc.Hget("keyword", "show").Interface()
	fmt.Println(bl, rt)
}

func TestBufferSet(t *testing.T) {
	defaultFile = "../../../buffer.json"
	bc := NewBufferClient()

	bl1 := bc.Hget("keyword", "zhangzhen").IsEmpty()
	bc.Hset("keyword", "zhangzhen", Empty)
	bl2 := bc.Hget("keyword", "zhangzhen").IsEmpty()
	fmt.Println(bl1, bl2)
}
