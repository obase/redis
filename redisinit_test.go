package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
)

func TestRedisDo(t *testing.T) {
	r := Get("demo")
	vl, ok, err := String(r.Do("DEL", "abc"))
	fmt.Println(vl, ok, err)
}

func TestRedisTx(t *testing.T) {
	r := Get("demo")
	for i := 0; i < 10; i++ {
		vl, err := redis.Values(r.Pi(func(op OP, keysArgs ...interface{}) (err error) {
			op.Do("SMEMBERS", keysArgs...)
			return
		}, "abc"))
		fmt.Println(redis.Strings(vl[0], nil))
		fmt.Println(err)
	}

}
