package common

import (
	"errors"
	"github.com/gomodule/redigo/redis"
	"time"
)

var (
	mPool *redis.Pool = nil
)

func InitRedis(addr string) {
	mPool = &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr)
			if err != nil {
				return nil, err
			}
			return c, err
		},

		MaxActive: 5,

		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			if time.Since(t) < time.Minute {
				return nil
			}
			_, err := c.Do("PING")
			return err
		},
	}
}

/*
 */
func DeleteIP(addr string) error {
	if mPool == nil {
		return errors.New("Call InitRedis(string) first")
	}

	c := mPool.Get()
	defer c.Close()

	_, err := c.Do("DEL", addr)

	return err
}

func InsertRequest(addr string) (int, error) {
	if mPool == nil {
		return -1, errors.New("Call InitRedis(string) first")
	}

	c := mPool.Get()
	defer c.Close()
	num, err := redis.Int64(c.Do("LLEN", addr))

	// get request number from this ip
	if err != nil {
		return -1, err
	} else if int(num) == 0 {
		c.Do("RPUSH", addr, time.Now().Format(time.RFC3339))
		return 1, nil
	}

	// check if the previous requests are expired
	for {
		s, err := redis.String(c.Do("LPOP", addr))

		if err != nil {
			break
		} else {
			t, _ := time.Parse(time.RFC3339, s)
			if time.Since(t) < time.Minute {
				c.Do("LPUSH", addr, s)
				break
			}
		}
	}

	num, err = redis.Int64(c.Do("LLEN", addr))
	if err != nil {
		return -1, err
	}
	if int(num) >= 60 { // more than 60 requests in 1 minute
		return 61, errors.New("Too many requests")
	}
	num, err = redis.Int64(c.Do("RPUSH", addr, time.Now().Format(time.RFC3339)))
	if err != nil {
		return -1, nil
	}

	// insert successfully and return the number
	return int(num), nil

}
