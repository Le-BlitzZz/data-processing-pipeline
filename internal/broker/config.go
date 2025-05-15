package broker

import "fmt"

type Config struct {
	user     string
	password string
	server   string
	dsn      string
}

func NewConfig(user, password, server string) *Config {
	return &Config{
		user:     user,
		password: password,
		server:   server,
	}
}

func (conf *Config) BrokerDsn() string {
	if conf.dsn != "" {
		return conf.dsn
	}

	dsn := fmt.Sprintf(
		"amqp://%s:%s@%s/",
		conf.user,
		conf.password,
		conf.server,
	)
	conf.dsn = dsn

	return conf.dsn
}
