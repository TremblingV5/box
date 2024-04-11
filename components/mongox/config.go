package mongox

import (
	"fmt"
	"net/url"
)

type Config struct {
	Host        string   `json:"host" yaml:"host"`
	Port        int      `json:"port" yaml:"port"`
	Username    string   `json:"username" yaml:"username"`
	Password    string   `json:"password" yaml:"password"`
	Database    string   `json:"database" yaml:"database"`
	Collections []string `json:"collections" yaml:"collections"`
}

func (c *Config) SetDefault() {
	if c.Host == "" {
		c.Host = "localhost"
	}

	if c.Port == 0 {
		c.Port = 27017
	}

	if c.Database == "" {
		c.Database = "default"
	}

	if len(c.Collections) == 0 {
		c.Collections = append(c.Collections, "default")
	}
}

func (c *Config) ToDSN() string {
	dsn := "mongodb://"
	if c.Username != "" {
		dsn += url.QueryEscape(c.Username)
	}

	if c.Password != "" {
		dsn += ":" + c.Password
	}

	if c.Username != "" || c.Password != "" {
		dsn += "@"
	}

	dsn += fmt.Sprintf("%s:%d", url.QueryEscape(c.Host), c.Port)
	return dsn
}
