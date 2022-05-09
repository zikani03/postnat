package postnat

import (
	"fmt"
	"io/ioutil"

	"github.com/BurntSushi/toml"
)

type Config struct {
	DB     PostgresConfig `toml:"postgres"`
	Nats   NatsConfig     `toml:"nats"`
	Topics TopicsConfig   `toml:"topics"`
}

type PostgresConfig struct {
	Username string `toml:"username"`
	Password string `toml:"password"`
	Database string `toml:"database"`
	Host     string `toml:"host"`
	Port     int    `toml:"port"`
	SSLMode  string `toml:"sslmode"`
}

type NatsConfig struct {
	URL           string `toml:"url"`
	MaxReconnects int    `toml:"max_reconnects"`
}

type TopicsConfig struct {
	Prefix            string   `toml:"prefix"`
	ListenFor         []string `toml:"listen_for"`
	ReplaceUnderscore bool     `toml:"replace_underscore_with_dot"`
}

func ParseConfig(path string) (*Config, error) {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}
	config := Config{}

	err = toml.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config: %v", err)
	}

	return &config, nil
}

func (c *Config) dbConnStr() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s database=%s sslmode=%s",
		c.DB.Host,
		c.DB.Port,
		c.DB.Username,
		c.DB.Password,
		c.DB.Database,
		c.DB.SSLMode,
	)
}

func (c *Config) natsConnStr() string {
	return c.Nats.URL
}
