package redis

import (
	"errors"
	"fmt"
	redigo "github.com/gomodule/redigo/redis"
	"net"
	"strings"
	"sync"
)

var (
	ErrInvalidClusterSlots = errors.New("Invalid cluster slots")
	ErrUnsupportOperation  = errors.New("unsupport operation")
	ErrArgumentException   = errors.New("argument exception")
	ErrUnsupportKeyCount   = errors.New("unsupport script key count, which must be 1 on cluster mode")
)

type redigoCluster struct {
	*Config
	sync.RWMutex
	Slots []*SlotInfo
	Pools []*redigoPool
	Index []*redigoPool
}

func newRedigoCluster(opt *Config) (*redigoCluster, error) {
	ret := &redigoCluster{
		Config: opt,
	}
	err := ret.UpdateClusterIndexes()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// 指定写锁. 中间可能会更新cluster indexes
func (rc *redigoCluster) UpdateClusterIndexes() (err error) {

	// 获取最新的slots并比较是否发生变化
	slots, err := QueryClusterSlots(rc.Config)
	if err != nil || !IsSlotsChanged(rc.Slots, slots) {
		return
	}

	rc.RWMutex.Lock()
	// 清除旧的连接
	rc.Close()
	// 重用旧的内存
	slen := len(slots)
	if len(rc.Pools) != slen {
		rc.Pools = make([]*redigoPool, slen)
	}
	if len(rc.Index) != CLUSTER_SLOTS_NUMBER {
		rc.Index = make([]*redigoPool, CLUSTER_SLOTS_NUMBER)
	}
	for i, s := range slots {

		o := cloneConfig(rc.Config)
		o.Address = []string{s.Address}

		rc.Pools[i], err = newRedigoPool(o)
		if err != nil {
			// 如果发生错误,需要级联清除已经创建的其他连接池
			for j := 0; j < i; j++ {
				rc.Pools[i].Close()
			}
			rc.RWMutex.Unlock()
			return
		}

		for j := s.Start; j <= s.End; j++ {
			rc.Index[j] = rc.Pools[i]
		}
	}

	rc.RWMutex.Unlock()
	return
}

// 指定读锁. 中间可能会更新cluster indexes
func (rc *redigoCluster) indexRedis(key string) *redigoPool {
	sl := Slot(key)
	rc.RWMutex.RLock()
	ret := rc.Index[sl]
	rc.RWMutex.RUnlock()
	return ret
}

/*--------------------------接口方法----------------------------------*/
func (rc *redigoCluster) Do(cmd string, keysArgs ...interface{}) (reply interface{}, err error) {

	if len(keysArgs) == 0 {
		return nil, ErrArgumentException
	}

	var key string
	if rc.Keyfix != "" {
		FillKeyfix4(&rc.Keyfix, &key, keysArgs)
	} else {
		key = keysArgs[0].(string)
	}
	reply, err = rc.indexRedis(key).do(cmd, keysArgs)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		reply, err = rc.indexRedis(key).do(cmd, keysArgs)
	}
	return

}

// 管道批量, 有可能部分成功.
func (rc *redigoCluster) Pi(bf Batch, keysArgs ...interface{}) (reply []interface{}, err error) {

	if len(keysArgs) == 0 {
		return nil, ErrArgumentException
	}
	var key string
	if rc.Keyfix != "" {
		FillKeyfix3(&rc.Keyfix, &key, keysArgs)
	} else {
		key = keysArgs[0].(string)
	}
	reply, err = rc.indexRedis(key).Pi(bf, keysArgs...)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		reply, err = rc.indexRedis(key).Pi(bf, keysArgs...)
	}
	return
}

// 事务批量, 要么全部成功, 要么全部失败.
func (rc *redigoCluster) Tx(bf Batch, keysArgs ...interface{}) (reply []interface{}, err error) {
	if len(keysArgs) == 0 {
		return nil, ErrArgumentException
	}
	var key string
	if rc.Keyfix != "" {
		FillKeyfix3(&rc.Keyfix, &key, keysArgs)
	} else {
		key = keysArgs[0].(string)
	}
	reply, err = rc.indexRedis(key).Tx(bf, keysArgs...)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		reply, err = rc.indexRedis(key).Tx(bf, keysArgs...)
	}
	return
}

// Publish
func (rc *redigoCluster) Pub(key string, msg interface{}) (err error) {
	if rc.Keyfix != "" {
		FillKeyfix1(&rc.Keyfix, &key)
	}
	err = rc.indexRedis(key).pub(key, msg)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		err = rc.indexRedis(key).pub(key, msg)
	}
	return

}

// Subscribe, 阻塞执行sf直到返回stop或error才会结束
func (rc *redigoCluster) Sub(key string, data SubDataFunc, meta SubMetaFunc) (err error) {
	if rc.Keyfix != "" {
		FillKeyfix1(&rc.Keyfix, &key)
	}
	err = rc.indexRedis(key).sub(key, data, meta)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		err = rc.indexRedis(key).sub(key, data, meta)
	}
	return

}

func (rc *redigoCluster) Eval(script string, keyCount int, keysArgs ...interface{}) (reply interface{}, err error) {
	if len(keysArgs) == 0 {
		return nil, ErrArgumentException
	}
	var key string
	if rc.Keyfix != "" {
		FillKeyfix4(&rc.Keyfix, &key, keysArgs)
	} else {
		key = keysArgs[0].(string)
	}
	reply, err = rc.indexRedis(key).eval(script, keyCount, keysArgs)
	if err != nil && IsSlotsError(err) {
		rc.UpdateClusterIndexes()
		reply, err = rc.indexRedis(key).eval(script, keyCount, keysArgs)
	}
	return
}

func (rc *redigoCluster) Close() {
	// 关闭操作不用锁,避免等待影响
	if len(rc.Pools) > 0 {
		for _, p := range rc.Pools {
			if p != nil {
				p.Close()
			}
		}
	}
}

/*--------------------------辅助方法----------------------------------*/
type SlotInfo struct {
	Start   int
	End     int
	Address string
}

func QueryClusterSlots(opt *Config) ([]*SlotInfo, error) {

	var rc redigo.Conn
	var err error
	var tcp net.Conn
	for _, address := range opt.Address {
		if tcp, err = net.DialTimeout(opt.Network, address, opt.ConnectTimeout); err == nil {
			rc = redigo.NewConn(tcp, opt.ReadTimeout, opt.WriteTimeout)
			break
		}
	}
	if err != nil {
		return nil, err
	}
	defer rc.Close()

	if opt.Password != "" {
		if _, err = rc.Do("AUTH", opt.Password); err != nil {
			return nil, err
		}
	}

	infos, err := redigo.Values(rc.Do("CLUSTER", "SLOTS"))
	if err != nil {
		return nil, err
	}

	/*-------------------------------------------------------------
	172.31.0.63:7000> CLUSTER SLOTS
	1) 1) (integer) 0
	   2) (integer) 5460
	   3) 1) "172.31.0.63"
		  2) (integer) 7000
		  3) "a585e144d73c9ca1c72e0dc14ba13b18eddddf61"
	   4) 1) "172.31.0.63"
		  2) (integer) 7003
		  3) "c079a3b1385faf1d1447b38d43941f75f2411f2b"
	2) 1) (integer) 5461
	   2) (integer) 10922
	   3) 1) "172.31.0.63"
		  2) (integer) 7001
		  3) "f0e5fd569ce7eaa63ab71174b7d4ae3cb34452b9"
	   4) 1) "172.31.0.63"
		  2) (integer) 7004
		  3) "b46b3907b9cbde9060715260684b64fe3fdd7729"
	3) 1) (integer) 10923
	   2) (integer) 16383
	   3) 1) "172.31.0.63"
		  2) (integer) 7002
		  3) "abfe332e0468b67cc63e4693d94273c4a135b448"
	   4) 1) "172.31.0.63"
		  2) (integer) 7005
		  3) "afb69dec265247d7d51496a7302922b0823a08fb"
	 --------------------------------------------------------------*/
	ret := make([]*SlotInfo, len(infos))
	plen := len(opt.Proxyips)
	for i, info := range infos {
		data := info.([]interface{})
		start, _, _ := Int(data[0], nil)
		end, _, _ := Int(data[1], nil)
		addrs := data[2].([]interface{})
		host, _, _ := String(addrs[0], nil)
		port, _, _ := Int(addrs[1], nil)

		// 替换成代理IP
		if plen > 0 {
			vhost := opt.Proxyips[host]
			if vhost != "" {
				host = vhost
			}
		}
		ret[i] = &SlotInfo{
			Start:   start,
			End:     end,
			Address: fmt.Sprintf("%s:%d", host, port),
		}
	}

	return ret, nil
}

func IsSlotsError(err error) bool {
	if rerr, ok := err.(redigo.Error); ok {
		msg := rerr.Error()
		if strings.HasPrefix(msg, "MOVED") || strings.HasPrefix(msg, "ASK") {
			return true
		}
	} else if _, ok := err.(*net.OpError); ok {
		return true
	}
	return false
}

// 比较新旧slots是否发生变化,避免全局更新索引影响全面
func IsSlotsChanged(slot1, slot2 []*SlotInfo) bool {
	slen := 0
	if slen = len(slot1); slen != len(slot2) {
		return true
	}

	flags := make([]bool, slen)
	for _, s1 := range slot1 {
		found := false
		for j, s2 := range slot2 {
			if !flags[j] && s1.Start == s2.Start {
				flags[j], found = true, true
				if s1.End != s2.End || s1.Address != s2.Address {
					return true
				}
			}
		}
		if !found {
			return true
		}
	}
	return false
}
