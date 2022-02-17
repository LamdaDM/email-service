package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
)

type Container struct {
	redisPool *redis.Pool
	emailOpts *EmailOpts
	config    *Config
}

func New() *Container {
	config := load()
	emailOpts := getServiceData()
	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Get("REDIS:PORT"))
			if err != nil {
				log.Fatal(err)
			}
			return c, err
		},
		MaxIdle:     8,
		MaxActive:   120,
		IdleTimeout: 600,
		Wait:        true,
	}

	return &Container{
		redisPool,
		&emailOpts,
		config,
	}
}

type EmailOpts struct {
	from     string
	password string
	mailHost string
	mailPort string
	template []byte
}

func getServiceData() EmailOpts {
	const (
		NFROM     = "FROM"
		NPASSWORD = "PASSWORD"
		NHOST     = "HOST"
		NPORT     = "PORT"
	)

	cfg := container.config.GetSection("EMAIL_PROVIDER")

	from := cfg.Get(NFROM)
	password := cfg.Get(NPASSWORD)
	mailHost := cfg.Get(NHOST)
	mailPort := cfg.Get(NPORT)
	template := template()

	return EmailOpts{
		from,
		password,
		mailHost,
		mailPort,
		template,
	}
}

func template() []byte {
	const TemplatePath = ".email.template"
	template, err := os.ReadFile(TemplatePath)
	if err != nil {
		log.Fatal(err)
	}
	return template
}
