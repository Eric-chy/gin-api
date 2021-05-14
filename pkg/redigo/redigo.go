package redigo

import (
	"encoding/json"
	"fmt"
	"ginpro/common/global"
	"ginpro/config"
	"github.com/gomodule/redigo/redis"
	"time"
)

func Init() {
	cfg := config.Conf.Redis
	global.RedisPool = &redis.Pool{
		MaxIdle:     cfg.MaxIdle,
		MaxActive:   cfg.MaxActive,
		IdleTimeout: cfg.IdleTimeout,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", cfg.Host+":"+cfg.Port)
			if err != nil {
				return nil, err
			}
			if cfg.Password != "" {
				if _, err := c.Do("AUTH", cfg.Password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
	conn := global.RedisPool.Get()
	defer conn.Close()
	if _, err := conn.Do("ping"); err != nil {
		fmt.Println("redis connect fail:", err)
		panic(err)
	}
}

func Set(key string, data interface{}, time int) error {
	conn := global.RedisPool.Get()
	defer conn.Close()

	value, err := json.Marshal(data)
	if err != nil {
		return err
	}

	_, err = conn.Do("SET", key, value, "EX", time)
	if err != nil {
		return err
	}

	//_, err = conn.Do("EXPIRE", key, time)
	//if err != nil {
	//	return err
	//}

	return nil
}

func Exists(key string) bool {
	conn := global.RedisPool.Get()
	defer conn.Close()

	exists, err := redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		return false
	}

	return exists
}

func GetString(key string) (string, error) {
	conn := global.RedisPool.Get()
	defer conn.Close()

	reply, err := redis.String(conn.Do("GET", key))
	if err != nil {
		return "", err
	}

	return reply, nil
}

func GetNum(key string) (int, error) {
	conn := global.RedisPool.Get()
	defer conn.Close()
	num, err := redis.Int(conn.Do("GET", key))
	if err != nil {
		return 0, err
	}
	return num, nil
}

func IncrBy(key string, num int) (int, error) {
	conn := global.RedisPool.Get()
	defer conn.Close()
	reply, err := redis.Int(conn.Do("INCRBY", key, num))
	if err != nil {
		return 0, err
	}

	return reply, nil
}

func Delete(key string) (bool, error) {
	conn := global.RedisPool.Get()
	defer conn.Close()

	return redis.Bool(conn.Do("DEL", key))
}

func LikeDeletes(key string) error {
	conn := global.RedisPool.Get()
	defer conn.Close()

	keys, err := redis.Strings(conn.Do("KEYS", "*"+key+"*"))
	if err != nil {
		return err
	}

	for _, key := range keys {
		_, err = Delete(key)
		if err != nil {
			return err
		}
	}

	return nil
}
