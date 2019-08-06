package redis

import (
	"github.com/obase/conf"
	"time"
)

const CKEY = "redis"

// 对接conf.yml, 读取原redis相关配置
func init() {
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
			_select, ok := conf.ElemInt(config, "select")
			proxyips, ok := conf.ElemStringMap(config, "proxyips")
			defalt, ok := conf.ElemBool(config, "default")

			option := &Config{
				Key:             key,
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
				Select:          _select,
				Cluster:         cluster,
				Proxyips:        proxyips,
				Default:         defalt,
			}

			if err := Setup(option); err != nil {
				panic(err)
			}
		}
	}
}
