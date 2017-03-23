package util

import (
	"log"
	"time"

	"github.com/garyburd/redigo/redis"
)

// Pool ...
var Pool *redis.Pool

// NewRedis ...
func NewRedis(connection string) {
	Pool = &redis.Pool{
		MaxIdle:     3,
		IdleTimeout: 10 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", connection)
			if err != nil {
				log.Println("Redis connect error ! connection name ->", connection, ", error log ->", err)
				return nil, err
			}
			return c, nil
		},
	}
}
