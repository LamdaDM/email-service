package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
)

func main() {
	IM()

	client := container.redisPool.Get()

	ch := container.config.Get("REDIS:CHANNEL_NAME")

	sub := redis.PubSubConn{Conn: client}

	err := sub.Subscribe(ch)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Listening on channel: %s\n", ch)

	for {
		switch reply := sub.Receive().(type) {
		case redis.Message:
			fmt.Printf("Received message: %s\n", string(reply.Data))
			go Mail(reply.Data)
		case error:
			log.Fatal(reply)
		default:
			continue
		}
	}
}

func IM() {
	opts := args()
	container = LoadContainer(opts.cfgPath, opts.templatePath)
}

var container *Container

type ContainerOpts struct {
	cfgPath      string
	templatePath string
}

func args() ContainerOpts {
	cfgPath, templatePath := ".cfg", ".email"

	path, found := os.LookupEnv("CFG_PATH")
	if found {
		cfgPath = path
	}

	path, found = os.LookupEnv("TEMPLATE_PATH")
	if found {
		templatePath = path
	}

	return ContainerOpts{
		cfgPath,
		templatePath,
	}
}
