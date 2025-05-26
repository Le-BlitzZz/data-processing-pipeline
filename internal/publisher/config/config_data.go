package config

import "path/filepath"

const dataDir = "datasets"

var dataSplits = []string{"train", "test", "val"}

func (c *Config) DataSplits() []string {
	if c.options.DataSplits == nil {
		return dataSplits
	}
	return c.options.DataSplits
}

func (c *Config) DataDir() string {
	if c.options.DataDir == "" {
		return dataDir
	}
	return c.options.DataDir
}

func (c *Config) DataSplitPath(split string) string {
	return filepath.Join(c.DataDir(), split+".csv")
}

func (c *Config) DataSplitPathMap() map[string]string {
	m := make(map[string]string, len(c.DataSplits()))
	for _, split := range dataSplits {
		m[split] = c.DataSplitPath(split)
	}
	return m
}
