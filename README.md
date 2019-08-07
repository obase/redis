# package redis
redis客户端封装

# Installation
- go get 
```
go get -u github.com/gomodule/redigo
go get -u github.com/obase/redis
```
- go mod
```
go mod edit -require=github.com/obase/redis@latest
```
# Configuration
```
# conf.yml 中的redis配置文件样例
redis:
  -
    # 引用的key(必需)
    key:
    # 地址(必需). 多值用逗号分隔
    address: "127.0.0.1:6379"
    # 是否集群(必需)
    cluster: false
    # 密码(可选)
    password:
    # keepalive间隔(可选). 默认空不设置
    keepalive: "1m"
    # 连接超时(可选). 默认空不设置
    connectTimeout: "1m"
    # 读超时(可选). 默认空不设置
    readTimeout: "1m"
    # 写超时(可选): 默认空不设置
    writeTimeout: "1m"
    # 连接池初始数量(可选). 默认为0
    initConns: 4
    # 连接池最大数量(可选). 默认没有限制
    maxConns: 256
    # 连接池最大空闲数量. 默认为initConns
    maxIdles:
    # 连接池测试空闲超时. 处理空闲的连接若超时会执行PING测试是否可用.
    testIdleTimeout: "20m"
    # 连接池达到最大链接数量立即报错还是阻塞等待
    errExceMaxConns: false
    # 统一后缀. 默认为空, 一般用于多个业务共用Redis集群的情况
    keyfix:
    # 支持Database下标, 默认0
    select: 0
    # 代理IP. 默认为空, 一般用于网关集群测试,自动将cluster slots的内网IP替换为外网IP.
    proxyips: {"127.0.0.1","192.168.2.21"}
```
# Index
- Constants
- Variables
- type OP
```
type OP interface {
	Do(cmd string, keysArgs ...interface{}) (err error)
}
```
- type SubDataFunc
```
type SubDataFunc func(data []byte)
```
- type SubMetaFunc
```
type SubMetaFunc func(count int)
```
- type Batch 
```
type Batch func(op OP, keysArgs ...interface{}) (err error)
```
- type Redis
```
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

```
- func Bool
```
func Bool(val interface{}, err error) (bool, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```
- func Int
```
func Int(val interface{}, err error) (int, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```
- func Int64
```
func Int64(val interface{}, err error) (int64, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Float64
```
func Float64(val interface{}, err error) (float64, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Bytes
```
func Bytes(val interface{}, err error) ([]byte, bool, error)  
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Bytes
```
func Bytes(val interface{}, err error) ([]byte, bool, error)  
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func String
```
func String(val interface{}, err error) (string, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func IntSlice
```
func IntSlice(reply interface{}, err error) ([]int, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Int64Slice
```
func Int64Slice(reply interface{}, err error) ([]int64, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Float64Slice
```
func Float64Slice(reply interface{}, err error) ([]float64, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func StringSlice
```
func StringSlice(reply interface{}, err error) ([]string, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func BytesSlice
```
func BytesSlice(reply interface{}, err error) ([][]byte, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func ValueSlice
```
func ValueSlice(reply interface{}, err error) ([]interface{}, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func IntMap
```
func IntMap(reply interface{}, err error) (map[string]int, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Int64Map
```
func Int64Map(reply interface{}, err error) (map[string]int64, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Float64Map
```
func Float64Map(reply interface{}, err error) (map[string]float64, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func StringMap
```
func StringMap(reply interface{}, err error) (map[string]string, bool, error)
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func ValueMap
```
func ValueMap(reply interface{}, err error) (map[string]interface{}, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func ValueMap
```
func ValueScoreSlice(reply interface{}, err error) ([]*ValueScore, bool, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

- func Scan
```
func Scan(sc Scanner, src interface{}, err error) (interface{}, error) 
```
返回对应类型结果, 各参数意义:
```
- val: 原始值
- err: 原始错误
- ret: 返回结果
- ok: 在redis是否存在? true表示存在, false表示不存在
- err: 返回错误
```

# Examples
```
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

func TestRedisSub(t *testing.T) {
	r := Get("demo")
	r.Sub("bcd", func(data []byte) {
		fmt.Println("data ", string(data))
	}, func(count int) {
		fmt.Println("meta ", count)
	})
}

func TestRedisPub(t *testing.T) {
	r := Get("demo")
	r.Pub("bcd", "123456")
}

```