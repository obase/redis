package redis

import (
	"fmt"
	"testing"
)

func TestRedisDo(t *testing.T) {
	r := Get("demo")
	vl, ok, err := String(r.Do("SET", "abc", "123"))
	fmt.Println(vl, ok, err)
}
