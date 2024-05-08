package redisx

import "fmt"

type Config struct {
	Host     string         `json:"host" yaml:"host"`
	Port     int            `json:"port" yaml:"port"`
	Password string         `json:"password" yaml:"password"`
	DBList   map[string]int `json:"db_list" yaml:"db_list"`
}

func (c *Config) SetDefault() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 6379
	}

	if len(c.DBList) == 0 {
		c.DBList["default"] = 0
	}
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
