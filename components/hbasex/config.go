package hbasex

type Config struct {
	Host string `json:"host" yaml:"host"`
}

func (c *Config) SetDefault() {
	if c.Host == "" {
		c.Host = "localhost"
	}
}
