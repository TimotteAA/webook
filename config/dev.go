//go:build !k8s

// asdsf go:build dev
// sdd go:build test
// dsf 34

// 没有k8s 这个编译标签
package config

// 本地连接
var Config = config{
	DB: DBConfig{
		// 本地连接
		DSN: "root:root@tcp(localhost:30002)/webook",
	},
	Redis: RedisConfig{
		Addr: "localhost:30033",
	},
}
