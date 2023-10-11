package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"geekgo-webook/internal/domain"
	"github.com/redis/go-redis/v9"
	"time"
)

var ErrKeyNotExist = redis.Nil

type UserCache interface {
	Get(ctx context.Context, uid int64) (domain.User, error)
	Set(ctx context.Context, u domain.User) error
}

type RedisUserCache struct {
	client     redis.Cmdable
	expiration time.Duration
}

func NewUserCache(client redis.Cmdable) UserCache {
	return &RedisUserCache{
		client:     client,
		expiration: time.Minute * 15,
	}
}

func (cache *RedisUserCache) Get(ctx context.Context, uid int64) (domain.User, error) {
	key := cache.key(uid)
	val, err := cache.client.Get(ctx, key).Bytes() // string to bytes
	if err != nil {
		return domain.User{}, err
	}
	var u domain.User
	err = json.Unmarshal(val, &u) // []byte转存到结构体 反序列化
	//json.Marshal() // 结构体 转成[]byte
	//if err != nil {
	//	return domain.User{}, err
	//}
	//return u, nil
	return u, err
}

func (cache *RedisUserCache) Set(ctx context.Context, u domain.User) error {
	val, err := json.Marshal(u) // 序列化
	if err != nil {
		return err
	}
	key := cache.key(u.Id)

	return cache.client.Set(ctx, key, val, cache.expiration).Err()
}

func (cache *RedisUserCache) key(id int64) string {
	return fmt.Sprintf("user:info:%d", id)
}
