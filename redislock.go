package redis

import (
	"errors"
	"math/rand"
	"strconv"
	"time"
)

var ErrFailed = errors.New("failed to acquire lock")

type Mutex struct {
	key   string        // redis locker key
	ex    time.Duration // 第一个参数：持有锁的最大时间，默认8s, 字符串，格式如：1s、1m
	tries int           // 第二个参数：获取锁的尝试次数，默认32次, 整数
	delay time.Duration // 第三个参数：获取锁的间隔时间，默认500ms，字符串，格式如：1s、1m
	pk    string        // 第四个参数：redis连接池的key值，如果不传，默认将使用mqdb, 字符串
	value string        // 锁有惟一键值

	redis Redis
}

func NewMutex(key string, ex time.Duration, tries int, delay time.Duration, pk string) *Mutex {
	value := strconv.FormatInt(time.Now().UnixNano(), 32) + "-" + strconv.FormatUint(rand.Uint64(), 32)
	return &Mutex{
		key:   key,
		ex:    ex,
		tries: tries,
		delay: delay,
		pk:    pk,
		value: value,
	}
}

var (
	ErrRedisNotFound = errors.New("mutext redis not found")
)

// Lock locks m. In case it returns an error on failure, you may retry to acquire the lock by calling this method again.
func (m *Mutex) Lock() error {
	r := Get(m.pk)
	if r == nil {
		return ErrRedisNotFound
	}

	var (
		pv  string
		err error
	)

	for i := 0; i < m.tries; i++ {
		if pv, _, err = String(r.Do("SET", m.key, m.value, "EX", m.ex.Seconds(), "NX")); pv != "" || err != nil {
			break
		}
		time.Sleep(m.delay)
	}

	if err != nil {
		return err
	} else if pv == "" {
		return ErrFailed
	}

	return nil
}

const UNLOCK_SCRIPT = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("DEL", KEYS[1]);else return 0;end;`

// Unlock unlocks m and returns the status of unlock. It is a run-time error if m is not locked on entry to Unlock.
func (m *Mutex) Unlock() bool {
	r := Get(m.pk)
	if r == nil {
		panic(ErrRedisNotFound)
	}
	pi, _, err := Int(r.Eval(UNLOCK_SCRIPT, 1, m.key, m.value))
	if err != nil {
		panic(err)
	}
	return pi > 0
}

const EXTEND_SCRIPT = `if redis.call("GET", KEYS[1]) == ARGV[1] then return redis.call("EXPIRE", KEYS[1], tonumber(ARGV[2]));else return 0;end;`

// Extend resets the mutex's expiry and returns the status of expiry extension. It is a run-time error if m is not locked on entry to Extend.
func (m *Mutex) Extend() bool {
	r := Get(m.pk)
	if r == nil {
		panic(ErrRedisNotFound)
	}
	pi, _, err := Int(r.Eval(EXTEND_SCRIPT, 1, m.key, m.value, m.ex.Seconds()))
	if err != nil {
		panic(err)
	}
	return pi > 0
}

// AcquireLock 获取指定名字的锁，如果此名字的所未生成，则会自动生成一个锁，
// 如果已生成，会返回现有的锁
// name: 锁的名字
// options:
//     第一个参数：持有锁的最大时间，默认8s, 字符串，格式如：1s、1m
//     第二个参数：获取锁的尝试次数，默认32次, 整数
//     第三个参数：获取锁的间隔时间，默认500ms，字符串，格式如：1s、1m
//     第四个参数：redis连接池的key值，如果不传，默认将使用mqdb, 字符串
func AcquireLock(key string, options ...interface{}) *Mutex {
	ex := 8 * time.Second
	tries := 32
	delay := 500 * time.Millisecond
	pk := "mqdb"
	switch len(options) {
	case 1:
		ex, _ = time.ParseDuration(options[0].(string))
	case 2:
		ex, _ = time.ParseDuration(options[0].(string))
		tries, _ = options[1].(int)
	case 3:
		ex, _ = time.ParseDuration(options[0].(string))
		tries, _ = options[1].(int)
		delay, _ = time.ParseDuration(options[2].(string))
	case 4:
		ex, _ = time.ParseDuration(options[0].(string))
		tries, _ = options[1].(int)
		delay, _ = time.ParseDuration(options[2].(string))
		pk = options[3].(string)
	}

	return NewMutex(key, ex, tries, delay, pk)
}
