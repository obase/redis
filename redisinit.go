package redis

import (
	"github.com/obase/conf"
)

const CKEY = "redis"

// 对接conf.yml, 读取原redis相关配置
func init() {
	var configs []*Option
	if ok := conf.Scan(CKEY, &configs); !ok || len(configs) == 0 {
		return
	}

	for _, config := range configs {
		if config.Key != "" {
			// 设置默认值
			if config.Network == "" {
				config.Network = "tcp"
			}
			if err := Setup(config); err != nil {
				panic(err)
			}
		}
	}
}
