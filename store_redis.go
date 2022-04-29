package base64Captcha

import (
	"context"
	"time"

	"github.com/axiaoxin-com/logging"
	"github.com/go-redis/redis/v8"
)

// RedisStore An object implementing Store interface
type RedisStore struct {
	redisClient *redis.Client
	expiration  time.Duration
	keyPrefix   string
}

// NewRedisStore return redis store
func NewRedisStore(redcli *redis.Client, expiration time.Duration, keyPrefix string) Store {
	return &RedisStore{
		redisClient: redcli,
		expiration:  expiration,
		keyPrefix:   keyPrefix,
	}
}

// Set RedisStore implementing Set method of Store interface
func (s *RedisStore) Set(id string, value string) error {
	key := s.keyPrefix + ":" + id
	return s.redisClient.Set(context.Background(), key, value, s.expiration).Err()
}

// Get RedisStore implementing Get method of Store interface
func (s *RedisStore) Get(id string, clear bool) (value string) {
	key := s.keyPrefix + ":" + id
	val, err := s.redisClient.Get(context.Background(), key).Result()
	if err != nil {
		logging.Errorf(nil, "redis get key:%v error:%v", key, err)
		return ""
	}

	if clear {
		if err := s.redisClient.Del(context.Background(), key).Err(); err != nil {
			logging.Errorf(nil, "redis del key:%v error:%v", key, err)
			return ""
		}
	}
	return val
}

// Verify RedisStore implementing Verify method of Store interface
func (s *RedisStore) Verify(id, answer string, clear bool) bool {
	if answer == "" {
		return false
	}
	val := s.Get(id, clear)
	return val == answer
}
