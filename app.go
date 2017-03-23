package main

import (
	"log"

	"github.com/garyburd/redigo/redis"
	"github.com/moeghifar/golang-workshop/src/util"
)

func init() {
	util.NewRedis("localhost:6379")
}

func main() {
	log.Println("hella world!")
	log.Println("============")
	redisPingRoutine()
}

func redisPingRoutine() {
	res, err := redisPing()
	if err != nil {
		log.Println(err)
	} else {
		log.Println(res)
	}
}

func redisPing() (string, error) {
	var err error

	RedisPool := util.Pool.Get()
	defer RedisPool.Close()

	resVal, err := RedisPool.Do("PING")
	returnValue, err := redis.String(resVal, err)
	return returnValue, err
}
