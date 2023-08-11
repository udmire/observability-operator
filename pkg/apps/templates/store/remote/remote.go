package remote

import "flag"

// Config holds configuration for accessing long-term storage.
type Config struct {
	StorageBackendConfig `yaml:",inline"`

	StoragePrefix string `yaml:"storage_prefix"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {}

type StorageBackendConfig struct {
	Backend string `yaml:"backend"`
}
