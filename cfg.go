package main

import (
	"log"
	"os"
	"strings"
)

const VarCfgPath = ".cfg"

type Config struct {
	internal map[string]string
}

func (c *Config) TryGet(iden string) (string, bool) {
	val, ok := c.internal[iden]
	return val, ok
}

func (c *Config) Get(iden string) string {
	val, ok := c.internal[iden]
	if !ok {
		log.Fatalf("\"%s\" was not found in loaded configuration.", iden)
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

func LoadConfig() *Config {
	data, err := os.ReadFile(VarCfgPath)
	if err != nil {
		log.Fatal(err)
	}

	internal := parseFile(data)

	loaded := &Config{internal}

	return loaded
}

func parseFile(data []byte) map[string]string {
	lines := strings.Split(string(data), "\n")

	out := make(map[string]string, len(lines))

	for i := 0; i < len(lines); i++ {
		currMap := strings.SplitN(lines[i], "=", 2)

		if len(currMap) == 2 {
			iden := strings.TrimSpace(currMap[0])
			val := strings.TrimSpace(currMap[1])
			out[iden] = val
		}
	}

	return out
}
