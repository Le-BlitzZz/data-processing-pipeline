package config

import "path/filepath"

const datasetsDir = "datasets"

var splits = []string{"train", "test", "val"}

func (c *Config) Splits() []string {
	return splits
}

func (c *Config) SplitPath(split string) string {
	return filepath.Join(datasetsDir, split+".csv")
}
