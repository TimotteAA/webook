package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
	"webook/internal/domain"
)

var ErrKeyNotExist = redis.Nil

type UserCache struct {
	// 外部传入redis的实例
	client redis.Cmdable
	// 过期时间
	expirationTime time.Duration
}

// NewUserCache
// A 用到了 B，B 一定是接口 => 这个是保证面向接口
// A 用到了 B，B 一定是 A 的字段 => 规避包变量、包方法，都非常缺乏扩展性
// A 用到了 B，A 绝对不初始化 B，而是外面注入 => 保持依赖注入(DI, Dependency Injection)和依赖反转(IOC)
func NewUserCache(client redis.Cmdable) *UserCache {
	return &UserCache{
		client:         client,
		expirationTime: time.Minute * 15,
	}
}

// 定义key的格式
func (cache *UserCache) key(id int64) string {
	return fmt.Sprintf("user.info.%v", id)
}

func (cache *UserCache) Set(ctx context.Context, u domain.User) error {
	// 先序列化
	val, err := json.Marshal(u)
	if err != nil {
		return err
	}
	key := cache.key(u.Id)
	fmt.Println("set ", key, u)
	return cache.client.Set(ctx, key, val, cache.expirationTime).Err()
}

func (cache *UserCache) Get(ctx context.Context, id int64) (domain.User, error) {
	key := cache.key(id)
	fmt.Println("Get ", key)
	result, err := cache.client.Get(ctx, key).Bytes()
	// 数据不存在，err = redis.Nil
	// 可能是别的错误
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(result, &u)
	if err != nil {
		return domain.User{}, err
	}
	return u, nil
}
