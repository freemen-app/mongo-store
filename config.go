package mongoStore

import (
	"fmt"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Config struct {
	Host        string        `json:"host"`
	Port        string        `json:"port"`
	DB          string        `config:"database" json:"db"`
	Username    string        `json:"username"`
	Password    string        `json:"password"`
	ConnTimeout time.Duration `config:"conn_timeout" json:"conn_timeout"`
	PoolSize    uint64        `config:"pool_size" json:"pool_size"`
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"mongodb://%s:%s@%s:%s/%s?ssl=false&authSource=admin",
		c.Username, c.Password, c.Host, c.Port, c.DB,
	)
}

func (c *Config) Validate() error {
	return validation.ValidateStruct(
		c,
		validation.Field(&c.Host, validation.Required, is.Host),
		validation.Field(&c.Port, validation.Required, is.Port),
		validation.Field(&c.DB, validation.Required),
		validation.Field(&c.Username, validation.Required),
		validation.Field(&c.Password, validation.Required),
		validation.Field(&c.ConnTimeout, validation.Required),
		validation.Field(&c.PoolSize, validation.Required),
	)
}
