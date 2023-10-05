package cache

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/allegro/bigcache/v3"
	"sync"
	"time"
)

type localCodeCache struct {
	lock  *sync.Mutex
	cache *bigcache.BigCache
}

func NewLocalCodeCache() CodeCache {
	// 验证码10分组过期
	bigCache, err := bigcache.NewBigCache(bigcache.DefaultConfig(10 * time.Minute))
	if err != nil {
		panic(err)
	}
	return &localCodeCache{lock: &sync.Mutex{}, cache: bigCache}
}

func (l *localCodeCache) Set(ctx context.Context, biz string, phone string, inputCode string) error {
	// check-do-something，先加锁
	l.lock.Lock()
	defer l.lock.Unlock()
	key := l.key(biz, phone)
	// 先判断key是不是存在
	data, err := l.cache.Get(key)
	if err == bigcache.ErrEntryNotFound {
		// 新发送手机短信，设置一下：发送的验证码、已验证次数
		codeData := CodeData{
			code:    inputCode,
			count:   3,
			setTime: time.Now().UnixMilli(),
		}
		encoded, err := l.encodeData(codeData)
		if err != nil {
			return ErrUnknownForCode
		}
		err = l.cache.Set(key, encoded)
		if err != nil {
			return ErrUnknownForCode
		}
		return nil
	}
	if err != nil {
		return ErrUnknownForCode
	}

	// 设置过，判断一下过期时间？
	c, err := l.decode(data)
	if err != nil {
		return ErrUnknownForCode
	}
	if time.Since(time.Unix(c.setTime, 0)) >= time.Minute {
		// 当前时间和设置时间超了1分钟，可以发送，重新set一下
		codeData := CodeData{
			code:    inputCode,
			count:   3,
			setTime: time.Now().UnixMilli(),
		}
		encoded, err := l.encodeData(codeData)
		if err != nil {
			return ErrUnknownForCode
		}
		err = l.cache.Set(key, encoded)
		if err != nil {
			return ErrUnknownForCode
		}
		return nil
	}
	return ErrCodeSendTooMany
}

func (l *localCodeCache) Verify(ctx context.Context, biz string, phone, inputCode string) (bool, error) {
	l.lock.Lock()
	defer l.lock.Unlock()

	key := l.key(biz, phone)
	cached, err := l.cache.Get(key)
	// 缓存不存在，或者别的问题，都认为别人在搞我
	if err != nil {
		return false, ErrUnknownForCode
	}
	codeData, err := l.decode(cached)
	if err != nil {
		return false, ErrUnknownForCode
	}

	// 用户已经使用了所有的尝试次数，有人在搞我
	if codeData.count <= 0 {
		return false, ErrCodeVerifyTooMany
	}

	// 验证码匹配
	if codeData.code == inputCode {
		codeData.count = -1
		encoded, err := l.encodeData(codeData)
		if err != nil {
			return false, ErrUnknownForCode
		}
		err = l.cache.Set(key, encoded)
		if err != nil {
			return false, ErrUnknownForCode
		}
		return true, nil
	}

	// 验证码不匹配
	codeData.count -= 1
	encoded, err := l.encodeData(codeData)
	if err != nil {
		return false, ErrUnknownForCode
	}
	err = l.cache.Set(key, encoded)
	if err != nil {
		return false, ErrUnknownForCode
	}
	return false, nil
}

func (l *localCodeCache) key(biz, phone string) string {
	return fmt.Sprintf("phone_code:%s:%s", biz, phone)
}

func (l *localCodeCache) encodeData(data CodeData) ([]byte, error) {
	var buffer bytes.Buffer
	encoder := json.NewEncoder(&buffer)
	err := encoder.Encode(data)
	if err != nil {
		return []byte{}, err
	}
	return buffer.Bytes(), nil
}

func (l *localCodeCache) decode(data []byte) (CodeData, error) {

	var buffer bytes.Buffer
	buffer.Write(data)

	decoder := json.NewDecoder(&buffer)
	var c CodeData
	err := decoder.Decode(&c)
	if err != nil {
		return CodeData{}, err
	}
	return c, nil
}

type CodeData struct {
	// 验证码
	code string
	// 验证次数
	count int
	// 设置的毫秒时间戳
	setTime int64
}
