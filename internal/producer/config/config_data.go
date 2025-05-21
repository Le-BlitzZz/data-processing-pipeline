package config

import "path/filepath"

const datasetsDir = "datasets"

var splits = []string{"train", "test", "val"}

func (c *Config) Splits() []string {
	if c.options.Splits == nil {
		return splits
	}
	return c.options.Splits
}

func (c *Config) DatasetsDir() string {
	if c.options.DatasetsDir == "" {
		return datasetsDir
	}
	return c.options.DatasetsDir
}

func (c *Config) SplitPath(split string) string {
	return filepath.Join(c.DatasetsDir(), split+".csv")
}

func (c *Config) SplitPathMap() map[string]string {
	m := make(map[string]string, len(c.Splits()))
	for _, split := range splits {
		m[split] = c.SplitPath(split)
	}
	return m
}
