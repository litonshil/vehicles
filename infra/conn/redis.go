package conn

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/redis/go-redis/v9"
	"strconv"
	"time"
	"vehicles/config"
	"vehicles/utils"
	"vehicles/utils/errutil"
)

var client *redis.Client

func ConnectCache() {
	conf := config.Cache().Redis

	//_logger.Info("connecting to cache at ", conf.Host, ":", conf.Port, "...")

	//fmt.Println("connecting to cache at ", conf.Host, ":", conf.Port, "...")
	//client = redis.NewClient(&redis.Options{
	//	Addr:     fmt.Sprintf("%s:%d", conf.Host, conf.Port),
	//	Password: conf.Password,
	//	DB:       conf.Database,
	//})
	userName := conf.Username
	pass := conf.Password
	opt, err := redis.ParseURL(fmt.Sprintf("redis://%s:%s@%s:%d/%d", userName, pass, conf.Host, conf.Port, conf.Database))
	if err != nil {
		panic(err)
	}

	client := redis.NewClient(opt)

	if _, err := client.Ping(context.Background()).Result(); err != nil {
		//_logger.Error("failed to connect cache: ", err)
		fmt.Println("failed to connect cache: ", err)
		panic(err)
	}

	//_logger.Info("cache connection successful...")
	fmt.Println("cache connection successful...")
}

type CacheClient struct{}

func NewCacheClient() *CacheClient {
	return &CacheClient{}
}

func (rc *CacheClient) Set(key string, value interface{}, ttl int) error {
	if utils.IsEmpty(key) || utils.IsEmpty(value) {
		return errutil.ErrEmptyRedisKeyValue
	}

	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(context.Background(), key, string(serializedValue), time.Duration(ttl)*time.Second).Err()
}

func (rc *CacheClient) SetString(key string, value string, ttl time.Duration) error {
	if utils.IsEmpty(key) || utils.IsEmpty(value) {
		return errutil.ErrEmptyRedisKeyValue
	}

	return client.Set(context.Background(), key, value, ttl*time.Second).Err()
}

func (rc *CacheClient) SetStruct(key string, value interface{}, ttl time.Duration) error {
	serializedValue, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return client.Set(context.Background(), key, string(serializedValue), ttl*time.Second).Err()
}

func (rc *CacheClient) Get(key string) (string, error) {
	return client.Get(context.Background(), key).Result()
}

func (rc *CacheClient) GetInt(key string) (int, error) {
	str, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(str)
}

func (rc *CacheClient) GetStruct(key string, outputStruct interface{}) error {
	serializedValue, err := client.Get(context.Background(), key).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(serializedValue), &outputStruct); err != nil {
		return err
	}

	return nil
}

func (rc *CacheClient) Del(keys ...string) error {
	return client.Del(context.Background(), keys...).Err()
}

//
//func (rc *CacheClient) DelPattern(pattern string) error {
//	iter := client.Scan(context.Background(), 0, pattern, 0).Iterator()
//
//	for iter.Next(context.Background(),) {
//		err := client.Del(iter.Val()).Err()
//		if err != nil {
//			return err
//		}
//	}
//
//	if err := iter.Err(); err != nil {
//		return err
//	}
//
//	return nil
//}

func (rc *CacheClient) Exists(key string) bool {
	exists, err := client.Exists(context.Background(), key).Result()
	if err != nil {
		return false
	}

	return exists == 1
}
