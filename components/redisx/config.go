package redisx

import "fmt"

type Config struct {
	Host     string          `json:"host" yaml:"host"`
	Port     int             `json:"port" yaml:"port"`
	Password string          `json:"password" yaml:"password"`
	DBList   []*DBListConfig `json:"db_list" yaml:"db_list"`
}

type DBListConfig struct {
	Name   string `json:"name" yaml:"name"`
	Number int    `json:"number" yaml:"number"`
}

func (c *Config) SetDefault() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 6379
	}

	if len(c.DBList) == 0 {
		c.DBList = append(c.DBList, &DBListConfig{
			Name:   "default",
			Number: 0,
		})
	}
}

func (c *Config) ToDSN() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
