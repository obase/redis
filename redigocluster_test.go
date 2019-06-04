package redis

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"testing"
	"time"
)

func TestRedigoCluster_Do(t *testing.T) {
	opt := &Option{
		Proxyips: map[string]string{
			"172.31.0.63": "120.92.144.252",
		},
		Address: []string{
			"120.92.144.252:7000",
			"120.92.144.252:7001",
			"120.92.144.252:7002",
		},
		Password:  "xxx@a123",
		MaxConns:  2,
		MaxIdles:  2,
		InitConns: 1,
	}
	cls, err := newRedigoCluster(MergeOption(opt))
	if err != nil {
		t.Fatal(err)
	}
	defer cls.Close()

	start := time.Now()
	var ps *string
	for i := 0; i < 1000*1; i++ {
		ps, err = String(cls.Do("get", "abc"))
		if err != nil {
			t.Fatal(err)
		}
	}
	if ps == nil {
		fmt.Println("nil....")
	} else {
		fmt.Printf("%p==%v\n", ps, *ps)
	}
	end := time.Now()
	fmt.Println("used(ms):", end.Sub(start).Nanoseconds()/1000000)
}

func TestRedigoCluster_Eval(t *testing.T) {
	opt := &Option{
		Address: []string{
			"192.168.2.11:6379",
		},
		Password:  "123456",
		MaxConns:  2,
		MaxIdles:  2,
		InitConns: 1,
	}
	cls, err := newRedigoCluster(MergeOption(opt))
	if err != nil {
		t.Fatal(err)
	}
	defer cls.Close()

	fmt.Println(redis.Int(cls.Eval("return redis.call('GET',KEYS[1])", 1, "abc", "refe")))
}

func TestKeyfix(t *testing.T) {
	c := Get("mqdb")
	f, err := c.Do("SET", "key", "don't do it")
	if err != nil {
		t.Fatal(f)
	}
	fmt.Printf("f is : %v\n", f)
	v, _ := String(c.Do("GET", "key"))
	fmt.Println(*v)
}