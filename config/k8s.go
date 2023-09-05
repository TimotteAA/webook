//go:build k8s

// 带有k8s的tag时，才编译
package config

var Config = config{
	DB:    DBConfig{DSN: "root:root@tcp(webook-mysql:11309)/webook"},
	Redis: RedisConfig{Addr: "webook-redis:6379"},
}
