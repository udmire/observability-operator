package remote

// Config holds configuration for accessing long-term storage.
type Config struct {
	StorageBackendConfig `yaml:",inline"`

	StoragePrefix string `yaml:"storage_prefix"`
}

type StorageBackendConfig struct {
	Backend string `yaml:"backend"`
}
