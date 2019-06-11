package redis

import (
	"github.com/obase/conf"
	"sync"
	"time"
)

const CKEY = "redis"

// 对接conf.yml, 读取原redis相关配置
var once sync.Once

func Init() {
	once.Do(func() {
		configs, ok := conf.GetSlice(CKEY)
		if !ok || len(configs) == 0 {
			return
		}

		for _, config := range configs {
			if key, ok := conf.ElemString(config, "key"); ok {
				address, ok := conf.ElemStringSlice(config, "address")
				cluster, ok := conf.ElemBool(config, "cluster")
				password, ok := conf.ElemString(config, "password")
				keepalive, ok := conf.ElemDuration(config, "keepalive")
				if !ok {
					keepalive = time.Minute
				}
				connectTimeout, ok := conf.ElemDuration(config, "connectTimeout")
				if !ok {
					connectTimeout = 30 * time.Second
				}
				readTimeout, ok := conf.ElemDuration(config, "readTimeout")
				if !ok {
					readTimeout = 30 * time.Second
				}
				writeTimeout, ok := conf.ElemDuration(config, "writeTimeout")
				if !ok {
					writeTimeout = 30 * time.Second
				}
				initConns, ok := conf.ElemInt(config, "initConns")
				maxConns, ok := conf.ElemInt(config, "maxConns")
				if !ok {
					maxConns = 16
				}
				maxIdles, ok := conf.ElemInt(config, "maxIdles")
				if !ok {
					maxIdles = 16
				}
				testIdleTimeout, ok := conf.ElemDuration(config, "testIdleTimeout")
				errExceMaxConns, ok := conf.ElemBool(config, "errExceMaxConns")
				if !ok {
					errExceMaxConns = false
				}
				keyfix, ok := conf.ElemString(config, "keyfix")
				proxyips, ok := conf.ElemStringMap(config, "proxyips")
				defalt, ok := conf.ElemBool(config, "default")

				option := &Option{
					Network:         "tcp",
					Address:         address,
					Keepalive:       keepalive,
					ConnectTimeout:  connectTimeout,
					ReadTimeout:     readTimeout,
					WriteTimeout:    writeTimeout,
					Password:        password,
					InitConns:       initConns,
					MaxConns:        maxConns,
					MaxIdles:        maxIdles,
					TestIdleTimeout: testIdleTimeout,
					ErrExceMaxConns: errExceMaxConns,
					Keyfix:          keyfix,
					Cluster:         cluster,
					Proxyips:        proxyips,
				}

				if cluster {
					if err := SetupCluster(key, option, defalt); err != nil {
						panic(err)
					}
				} else {
					if err := SetupPool(key, option, defalt); err != nil {
						panic(err)
					}
				}
			}
		}
	})
}
