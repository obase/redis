package redis

import (
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"log"
	"reflect"
	"testing"
	"time"
)

const TIMES = 10000 * 100

func TestRedis(t *testing.T) {
	c, err := redigo.Dial("tcp", "120.92.144.252:7000")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	if _, err = c.Do("AUTH", "KingSoft1239002nx624@a123"); err != nil {
		log.Fatal(err)
	}

	//var sgs []interface{}
	//if sgs, err = redis.Values(c.Do("cluster", "slots")); err != nil {
	//	log.Fatal(err)
	//}
	//for _, sg := range sgs {
	//	ss := sg.([]interface{})
	//	fmt.Println("------------------------------", len(ss))
	//	fmt.Println("start: ", reflect.TypeOf(ss[0]))
	//	fmt.Println("end: ", reflect.TypeOf(ss[1]))
	//	fmt.Println("ip: ", reflect.ValueOf(ss[2]))
	//	//fmt.Println("addr: ", reflect.TypeOf(ss[4]))
	//
	//}
	if _, err = c.Do("SET", "thisisf1", "123"); err != nil {
		fmt.Printf("ERR:%#v, %v\n", err, reflect.TypeOf(err))
		//log.Fatal(err)
	}
}

////////////////////////////////////////////////////////////////////////////////////
func TestGetClusterSlots(t *testing.T) {

	Redigo()
}

////////////////////////////////////////////////////////////////////////////////////
func Redigo() {
	c, err := redigo.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()
	//
	//if _, err = c.Do("AUTH", "123"); err != nil {
	//	log.Fatal(err)
	//}

	start := time.Now()
	for i := 0; i < TIMES; i++ {
		if _, err = c.Do("SET", "abc", 123); err != nil {
			log.Fatal(err)
		}
		if _, err = c.Do("GET", "abc"); err != nil {
			log.Fatal(err)
		}
	}
	fmt.Printf("redigo used: %v\n", time.Now().Sub(start).Nanoseconds()/1000000)
}

func TestList(t *testing.T) {
	ch := make(chan *redigoConn, 2)
	c := &redigoConn{}
	ch <- c
	c.T = time.Now().Unix()
	c = <-ch
	fmt.Println(c)
}
