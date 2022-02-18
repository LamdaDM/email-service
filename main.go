package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
)

func main() {

	client := container.redisPool.Get()

	sub := redis.PubSubConn{Conn: client}

	err := sub.Subscribe(container.config.Get("REDIS:CHANNEL_NAME"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		switch reply := sub.Receive().(type) {
		case redis.Message:
			Mail(reply.Data)
		case error:
			log.Fatal(reply)
		default:
			continue
		}
	}
}

var container *Container

func init() {
	container = LoadContainer()
}
