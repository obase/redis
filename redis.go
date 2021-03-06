package redis

import (
	"errors"
	"strings"
	"time"
)

/*
1. 支持pass
2. 支持pool
3. 支持cluster
4. 支持translation, 事务批量, 要么全部成功, 要么全部失败.
5. 支持pipeline, 管道批量, 有可能部分成功
6. 支持pub/sub
*/

/*====================辅助方法====================*/
type ValueScore struct {
	Value string
	Score float64
}

// 用于转换结果为struct或其他类型
type Scanner interface {
	Scan(reply interface{}) (interface{}, error)
}

/*====================操作接口====================*/
type OP interface {
	Do(cmd string, keysArgs ...interface{}) (err error)
}
type SubDataFunc func(data []byte)
type SubMetaFunc func(count int)
type Batch func(op OP, keysArgs ...interface{}) (err error)

type Redis interface {
	Do(cmd string, keysArgs ...interface{}) (reply interface{}, err error)
	// 注意: 集群模式不支持
	// 管道批量, 有可能部分成功.
	Pi(bf Batch, keysArgs ...interface{}) (reply []interface{}, err error)
	// 注意: 集群模式不支持
	// 事务批量, 要么全部成功, 要么全部失败.
	Tx(bf Batch, keysArgs ...interface{}) (reply []interface{}, err error)
	// Publish
	Pub(key string, msg interface{}) (err error)
	// Subscribe, 阻塞执行sf直到返回stop或error才会结束
	Sub(key string, data SubDataFunc, meta SubMetaFunc) (err error)
	// Script, 执行Lua脚本, 集群模式只支持单个KEYS
	Eval(script string, keys int, keysArgs ...interface{}) (reply interface{}, err error)
	//关闭清除链接
	Close()
}

/*====================选项设置====================*/
type Config struct {
	Key string `json:"key" bson:"key" yaml:"key"`
	// Conn参数
	Network        string        `json:"network" bson:"network" yaml:"network"`                      // 网络类簇,默认TCP
	Address        []string      `json:"address" bson:"address" yaml:"address"`                      //连接的ip:port, 默认127.0.0.1:6379.
	Keepalive      time.Duration `json:"keepalive" bson:"keepalive" yaml:"keepalive"`                //KeepAlive的间隔, 默认0不开启keepalive
	ConnectTimeout time.Duration `json:"connectTimeout" bson:"connectTimeout" yaml:"connectTimeout"` //连接超时, 默认0不设置
	ReadTimeout    time.Duration `json:"readTimeout" bson:"readTimeout" yaml:"readTimeout"`          // 读超时, 默认0永远不超时
	WriteTimeout   time.Duration `json:"writeTimeout" bson:"writeTimeout" yaml:"writeTimeout"`       // 写超时, 默认0永远不超时
	Password       string        `json:"password" bson:"password" yaml:"password"`                   //密码
	// Pool参数
	InitConns       int           `json:"initConns" bson:"initConns" yaml:"initConns"`                   //初始链接数, 默认0
	MaxConns        int           `json:"maxConns" bson:"maxConns" yaml:"maxConns"`                      //最大链接数, 默认0永远不限制
	MaxIdles        int           `json:"maxIdles" bson:"maxIdles" yaml:"maxIdles"`                      //最大空闲数, 超出会在用完后自动关闭, 默认为InitConns
	TestIdleTimeout time.Duration `json:"testIdleTimeout" bson:"testIdleTimeout" yaml:"testIdleTimeout"` //最大空闲超时, 超出会在获取时执行PING,如果失败则舍弃重建. 默认0表示不处理. 该选项是TestOnBorrow的一种优化
	ErrExceMaxConns bool          `json:"errExceMaxConns" bson:"errExceMaxConns" yaml:"errExceMaxConns"` // 达到最大链接数, 是等待还是报错. 默认false等待
	Keyfix          string        `json:"keyfix" bson:"keyfix" yaml:"keyfix"`                            // Key的统一后缀. 兼容此前的name情况, 不建议使用
	Select          int           `json:"select" bson:"select" yaml:"select"`                            // 选择DB下标, 默认0
	Cluster         bool          `json:"cluster" bson:"cluster" yaml:"cluster"`                         //是否集群

	// cluster参数
	Proxyips map[string]string `json:"proxyips" bson:"proxyips" yaml:"proxyips"` //代理IP集合,一般用于本地测试用

	// 是否默认
	Default bool `json:"default" bson:"default" yaml:"default"`
}

func cloneConfig(opt *Config) (ret *Config) {
	ret = new(Config)
	ret.Network = opt.Network
	ret.Address = opt.Address
	ret.Keepalive = opt.Keepalive
	ret.ConnectTimeout = opt.ConnectTimeout
	ret.ReadTimeout = opt.ReadTimeout
	ret.WriteTimeout = opt.WriteTimeout
	ret.Password = opt.Password

	ret.InitConns = opt.InitConns
	ret.MaxConns = opt.MaxConns
	ret.MaxIdles = opt.MaxIdles
	ret.TestIdleTimeout = opt.TestIdleTimeout
	ret.ErrExceMaxConns = opt.ErrExceMaxConns
	ret.Keyfix = opt.Keyfix
	ret.Select = opt.Select
	ret.Cluster = opt.Cluster // 下发集群信息

	ret.Proxyips = opt.Proxyips
	return
}

func mergeConfig(opt *Config) (ret *Config) {
	if opt == nil {
		ret = &Config{}
	} else {
		ret = opt
	}
	if ret.Network == "" {
		ret.Network = "tcp"
	}
	if ret.MaxIdles == 0 {
		ret.MaxIdles = ret.InitConns
	}
	if ret.Keepalive == 0 {
		ret.Keepalive = 10 * time.Minute
	}
	return
}

/*====================操作错误====================*/

var ErrExceedMaxConns = errors.New("excced the max conns")

/*====================操作设置====================*/

var (
	//默认
	Default Redis
	//其他
	Clients map[string]Redis = make(map[string]Redis)
)

func Get(name string) Redis {
	if rt, ok := Clients[name]; ok {
		return rt
	}
	return nil
}

func Setup(opt *Config) error {

	keys := strings.Split(opt.Key, ",")
	for _, k := range keys {
		if _, ok := Clients[k]; ok {
			return errors.New("duplicate redis key " + k)
		}
	}

	var c Redis
	// 注意: 不能直接将concreate type赋值给interface, 前者为nil时, 后者也不为会nil
	if opt.Cluster {
		if cp, err := newRedigoCluster(mergeConfig(opt)); err != nil {
			return err
		} else if cp != nil {
			c = cp
		}
	} else {
		if cc, err := newRedigoPool(mergeConfig(opt)); err != nil {
			return err
		} else if cc != nil {
			c = cc
		}
	}
	if c != nil {
		for _, k := range keys {
			Clients[k] = c
		}
		if opt.Default {
			Default = c
		}
	}
	return nil
}

/*====================slots约定====================*/
const CLUSTER_SLOTS_NUMBER = 16384 //redis cluster fixed slots

var tab = [256]uint16{
	0x0000, 0x1021, 0x2042, 0x3063, 0x4084, 0x50a5, 0x60c6, 0x70e7,
	0x8108, 0x9129, 0xa14a, 0xb16b, 0xc18c, 0xd1ad, 0xe1ce, 0xf1ef,
	0x1231, 0x0210, 0x3273, 0x2252, 0x52b5, 0x4294, 0x72f7, 0x62d6,
	0x9339, 0x8318, 0xb37b, 0xa35a, 0xd3bd, 0xc39c, 0xf3ff, 0xe3de,
	0x2462, 0x3443, 0x0420, 0x1401, 0x64e6, 0x74c7, 0x44a4, 0x5485,
	0xa56a, 0xb54b, 0x8528, 0x9509, 0xe5ee, 0xf5cf, 0xc5ac, 0xd58d,
	0x3653, 0x2672, 0x1611, 0x0630, 0x76d7, 0x66f6, 0x5695, 0x46b4,
	0xb75b, 0xa77a, 0x9719, 0x8738, 0xf7df, 0xe7fe, 0xd79d, 0xc7bc,
	0x48c4, 0x58e5, 0x6886, 0x78a7, 0x0840, 0x1861, 0x2802, 0x3823,
	0xc9cc, 0xd9ed, 0xe98e, 0xf9af, 0x8948, 0x9969, 0xa90a, 0xb92b,
	0x5af5, 0x4ad4, 0x7ab7, 0x6a96, 0x1a71, 0x0a50, 0x3a33, 0x2a12,
	0xdbfd, 0xcbdc, 0xfbbf, 0xeb9e, 0x9b79, 0x8b58, 0xbb3b, 0xab1a,
	0x6ca6, 0x7c87, 0x4ce4, 0x5cc5, 0x2c22, 0x3c03, 0x0c60, 0x1c41,
	0xedae, 0xfd8f, 0xcdec, 0xddcd, 0xad2a, 0xbd0b, 0x8d68, 0x9d49,
	0x7e97, 0x6eb6, 0x5ed5, 0x4ef4, 0x3e13, 0x2e32, 0x1e51, 0x0e70,
	0xff9f, 0xefbe, 0xdfdd, 0xcffc, 0xbf1b, 0xaf3a, 0x9f59, 0x8f78,
	0x9188, 0x81a9, 0xb1ca, 0xa1eb, 0xd10c, 0xc12d, 0xf14e, 0xe16f,
	0x1080, 0x00a1, 0x30c2, 0x20e3, 0x5004, 0x4025, 0x7046, 0x6067,
	0x83b9, 0x9398, 0xa3fb, 0xb3da, 0xc33d, 0xd31c, 0xe37f, 0xf35e,
	0x02b1, 0x1290, 0x22f3, 0x32d2, 0x4235, 0x5214, 0x6277, 0x7256,
	0xb5ea, 0xa5cb, 0x95a8, 0x8589, 0xf56e, 0xe54f, 0xd52c, 0xc50d,
	0x34e2, 0x24c3, 0x14a0, 0x0481, 0x7466, 0x6447, 0x5424, 0x4405,
	0xa7db, 0xb7fa, 0x8799, 0x97b8, 0xe75f, 0xf77e, 0xc71d, 0xd73c,
	0x26d3, 0x36f2, 0x0691, 0x16b0, 0x6657, 0x7676, 0x4615, 0x5634,
	0xd94c, 0xc96d, 0xf90e, 0xe92f, 0x99c8, 0x89e9, 0xb98a, 0xa9ab,
	0x5844, 0x4865, 0x7806, 0x6827, 0x18c0, 0x08e1, 0x3882, 0x28a3,
	0xcb7d, 0xdb5c, 0xeb3f, 0xfb1e, 0x8bf9, 0x9bd8, 0xabbb, 0xbb9a,
	0x4a75, 0x5a54, 0x6a37, 0x7a16, 0x0af1, 0x1ad0, 0x2ab3, 0x3a92,
	0xfd2e, 0xed0f, 0xdd6c, 0xcd4d, 0xbdaa, 0xad8b, 0x9de8, 0x8dc9,
	0x7c26, 0x6c07, 0x5c64, 0x4c45, 0x3ca2, 0x2c83, 0x1ce0, 0x0cc1,
	0xef1f, 0xff3e, 0xcf5d, 0xdf7c, 0xaf9b, 0xbfba, 0x8fd9, 0x9ff8,
	0x6e17, 0x7e36, 0x4e55, 0x5e74, 0x2e93, 0x3eb2, 0x0ed1, 0x1ef0,
}

func Slot(key string) uint16 {
	bs := []byte(key)
	ln := len(bs)
	start, end := 0, ln
	for i := 0; i < ln; i++ {
		if bs[i] == '{' {
			for j := i + 1; j < ln; j++ {
				if bs[j] == '}' {
					start, end = i, j
					break
				}
			}
			break
		}
	}
	crc := uint16(0)
	for i := start; i < end; i++ {
		index := byte(crc>>8) ^ bs[i]
		crc = (crc << 8) ^ tab[index]
	}
	return crc % CLUSTER_SLOTS_NUMBER
}

/*
第一种形式, 使用key更新key
*/
func FillKeyfix1(fix *string, key *string) {
	*key = *key + "." + *fix
}

/*
第二种形式: 使用keysArgs[0], 更新keysArgs[0]
*/
func FillKeyfix2(fix *string, keysArgs []interface{}) {
	keysArgs[0] = keysArgs[0].(string) + "." + *fix
}

/*
第三种形式: 使用keysArgs[0] 更新key
*/
func FillKeyfix3(fix *string, key *string, keysArgs []interface{}) {
	*key = keysArgs[0].(string) + "." + *fix
}

/*
第四种形式: 使用keysArgs[0] 更新key与keysArgs[0]
*/
func FillKeyfix4(fix *string, key *string, keysArgs []interface{}) {
	*key = keysArgs[0].(string) + "." + *fix
	keysArgs[0] = *key
}
