package broker

import "fmt"

type Config struct {
	User     string
	Password string
	Server   string
}

func NewConfig(user, password, server string) *Config {
	return &Config{
		User:     user,
		Password: password,
		Server:   server,
	}
}

func (conf *Config) BrokerDsn() string {
	return fmt.Sprintf(
		"amqp://%s:%s@%s/",
		conf.User,
		conf.Password,
		conf.Server,
	)
}
