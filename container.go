package main

import (
	"context"
	"github.com/gomodule/redigo/redis"
	"log"
	"os"
	"regexp"
	"strconv"
)

type Container struct {
	redisPool *redis.Pool
	emailOpts *EmailOpts
	config    *Config
	logger    Logger
}

func LoadContainer(cfgPath string, templatePath string) *Container {
	ctx := context.Background()

	config := LoadConfig(cfgPath)

	emailOpts := getServiceData(config, templatePath)

	redisPool := &redis.Pool{
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", config.Get("REDIS:ADDRESS"))
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

	trunkPrefixLen, _ := strconv.Atoi(config.Get("TRUNK:PREFIX_LEN"))
	trunkAddr := config.Get("TRUNK:ADDRESS")
	logger := TrunkLoggerInit(trunkAddr, ctx, os.Stderr, trunkPrefixLen)

	return &Container{
		redisPool,
		&emailOpts,
		config,
		logger,
	}
}

type EmailOpts struct {
	from     string
	password string
	mailHost string
	mailPort string
	template []byte
}

func getServiceData(config *Config, templatePath string) EmailOpts {
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
	template := template(templatePath)

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
