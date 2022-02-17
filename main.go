package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"strings"
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
			args := strings.SplitN(string(reply.Data), " ", 2)

			if len(args) != 2 {
				continue
			}

			to := []string{args[0]}
			message := []byte(args[1])
			mail(to, message)
		case error:
			log.Fatal(reply)
		default:
			continue
		}
	}
}

var container *Container

func init() {
	container = New()
}
