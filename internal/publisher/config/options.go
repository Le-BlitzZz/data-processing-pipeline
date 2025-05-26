package config

import "github.com/urfave/cli/v2"

type Options struct {
	BrokerServer                string
	BrokerUser                  string
	BrokerPassword              string
	BrokerRawExchange           string
	BrokerRawProcessingExchange string
	DataDir                     string
	DataSplits                  []string
}

func NewOptions(ctx *cli.Context) *Options {
	o := &Options{}

	if ctx == nil {
		return o
	}

	return o
}
