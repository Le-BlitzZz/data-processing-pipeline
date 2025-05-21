package config

import "github.com/urfave/cli/v2"

type Options struct {
	BrokerServer            string
	BrokerUser              string
	BrokerPassword          string
	BrokerRawExchange       string
	BrokerProcessedExchange string
	BrokerRawQueue          string
	BrokerProcessedQueue    string
	DatabaseUser            string
	DatabasePassword        string
	DatabaseServer          string
	DatabaseName            string
	DatabaseTimeout         int
	HttpHost                string
	HttpPort                int
}

func NewOptions(ctx *cli.Context) *Options {
	o := &Options{}

	if ctx == nil {
		return o
	}

	return o
}
