package database

import "fmt"

type Config struct {
	User     string
	Password string
	Server   string
	Name     string
	Timeout  int
}

func NewConfig(user, password, server, name string, timeout int) *Config {
	return &Config{
		User:     user,
		Password: password,
		Server:   server,
		Name:     name,
		Timeout:  timeout,
	}
}

func (conf *Config) DatabaseDsn() string {
	dbServer := fmt.Sprintf("tcp(%s)", conf.Server)

	return fmt.Sprintf(
		"%s:%s@%s/%s?charset=utf8mb4,utf8&collation=utf8mb4_unicode_ci&parseTime=true&timeout=%ds",
		conf.User,
		conf.Password,
		dbServer,
		conf.Name,
		conf.Timeout,
	)
}
