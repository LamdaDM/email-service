package main

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Config struct {
	internal map[string]string
}

func (c *Config) TryGet(key string) (string, bool) {
	val, ok := c.internal[key]
	return val, ok
}

func (c *Config) Get(key string) string {
	val, ok := c.internal[key]
	if !ok {
		fmt.Printf("\"%s\" was not found in loaded configuration.\n", key)
	}

	return val
}

func (c *Config) GetSection(section string) *Config {
	if section[len(section)-1] != ':' {
		section = section + ":"
	}
	internal := make(map[string]string, len(c.internal))
	for key, value := range c.internal {
		if strings.HasPrefix(key, section) {
			internal[strings.TrimPrefix(key, section)] = value
		}
	}

	return &Config{internal}
}

func LoadConfig(cfgPath string) *Config {
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		log.Fatal(err)
	}

	internal := parseFile(data)

	return &Config{internal}
}

func parseFile(data []byte) map[string]string {
	lines := strings.Split(string(data), "\n")

	out := make(map[string]string, len(lines))

	for i := 0; i < len(lines); i++ {
		currMap := strings.SplitN(lines[i], "=", 2)

		if len(currMap) == 2 {
			key := strings.TrimSpace(currMap[0])
			val := strings.TrimSpace(currMap[1])
			out[key] = val
		}
	}

	return out
}
