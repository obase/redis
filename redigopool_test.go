package redis

import (
	"fmt"
	"strconv"
	"sync"
	"testing"
	"time"
)

func TestNewRedisPool(t *testing.T) {
	opt := &Config{
		Network:         "tcp",
		Address:         []string{"127.0.0.1:6379"},
		MaxConns:        16,
		MaxIdles:        16,
		InitConns:       1,
		TestIdleTimeout: 1 * time.Minute,
	}
	pool, err := newRedigoPool(mergeConfig(opt))
	if err != nil {
		t.Fatal(err)
	}
	t1 := time.Now()
	wg := new(sync.WaitGroup)
	for i := 0; i < 500; i++ {
		wg.Add(1)
		go func(id string) {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				//fmt.Println(id, " before pool.Get")
				c, err := pool.Get()
				//fmt.Println(id, " after pool.Get")
				if err != nil {
					t.Fatal(err)
				}
				//fmt.Println(id, " before time.Sleep")
				//time.Sleep(3 * time.Second)
				//fmt.Println(id, " after time.Sleep")
				//if i == 3 {
				//	err = errors.New("testing")
				//}
				//fmt.Println(id, " before pool.Put")
				pool.Put(c, &err)
				//fmt.Println(id, " after pool.Put")

			}
		}("groutine_" + strconv.Itoa(i))
	}
	wg.Wait()
	t2 := time.Now()
	fmt.Printf("used(ms):%v\n", t2.Sub(t1).Nanoseconds()/1000000)
}

func TestRedigoPool_Pi(t *testing.T) {

	pool := Get("demo")

	t1 := time.Now()
	for i := 1; i < 10000*1; i++ {
		//ret, err := pool.Do("SET", "abc", 12.3)
		//if err != nil {
		//	t.Fatal(err)
		//}
		////fmt.Println(ret)
		//_ = ret

		ret, _, err := StringMap(pool.Do("get", "abdcfxxxx"))
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(ret)
	}
	t2 := time.Now()
	fmt.Println("used(ms):", t2.Sub(t1).Nanoseconds()/1000000)
	//fmt.Println(ret)

}

func TestRedigoPool_Tx(t *testing.T) {

	p := Get("mqdb")

	p.Tx(func(op OP, args ...interface{}) (err error) {
		op.Do("SET", "key1", "12345")
		//err = errors.New("testing")
		return
	})
	v, _, _ := String(p.Do("GET", "key1"))
	fmt.Println(v)
}

func TestRedigoPool_Pub(t *testing.T) {
	opt := &Config{
		Network:         "tcp",
		Address:         []string{"127.0.0.1:6379"},
		MaxConns:        1,
		MaxIdles:        1,
		InitConns:       1,
		TestIdleTimeout: 1 * time.Minute,
	}
	p, err := newRedigoPool(mergeConfig(opt))
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	rets, err := p.Pi(func(op OP, args ...interface{}) (err error) {
		//op.Do("SET","key1", "12345")
		op.Do("GET", "key1234")
		return
	})

	fmt.Println("===================>")
	fmt.Println(rets)
	//fmt.Println(String(rets[1], err))
}

func TestRedigoPool_Sub(t *testing.T) {

}

func TestRedigoPool_Do(t *testing.T) {
	opt := &Config{
		Network:         "tcp",
		Address:         []string{"127.0.0.1:6379"},
		MaxConns:        1,
		MaxIdles:        1,
		InitConns:       1,
		TestIdleTimeout: 1 * time.Minute,
	}
	p, err := newRedigoPool(mergeConfig(opt))
	if err != nil {
		t.Fatal(err)
	}
	defer p.Close()

	if _, err = p.Do("AUTH", "KingSoft1239002nx624@a123"); err != nil {
		t.Fatal(err)
	}

	ret, _, err := String(p.Do("get", "abc"))
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ret)
}
