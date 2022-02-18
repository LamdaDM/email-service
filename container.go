package main

import (
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"regexp"
)

type Container struct {
	redisPool *redis.Pool
	emailOpts *EmailOpts
	config    *Config
}

func LoadContainer() *Container {
	config := LoadConfig()
	emailOpts := getServiceData(config)
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

const TemplatePath = ".email.template"

func getServiceData(config *Config) EmailOpts {
	const (
		NFROM     = "FROM"
		NPASSWORD = "PASSWORD"
		NHOST     = "HOST"
		NPORT     = "PORT"
	)

	cfg := config.GetSection("EMAIL_PROVIDER")

	from := cfg.Get(NFROM)
	password := cfg.Get(NPASSWORD)
	mailHost := cfg.Get(NHOST)
	mailPort := cfg.Get(NPORT)
	template := template(TemplatePath)

	pattern, err := regexp.Compile("([^\r]|)([\n])")
	if err != nil {
		log.Fatal(err)
	}

	template = pattern.ReplaceAll(template, []byte("$1\r\n"))

	return EmailOpts{
		from,
		password,
		mailHost,
		mailPort,
		template,
	}
}

func template(templatePath string) []byte {
	template, err := os.ReadFile(templatePath)
	if err != nil {
		log.Fatal(err)
	}
	return template
}
