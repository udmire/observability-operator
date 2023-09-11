package store

import (
	"flag"
	"fmt"
	"strings"

	"github.com/grafana/dskit/cache"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	"github.com/udmire/observability-operator/pkg/templates/store/local"
)

var supportedCacheBackends = []string{cache.BackendMemcached, cache.BackendRedis}

// Config configures a rule store.
type Config struct {
	Local local.Config `yaml:"local"`

	// Cache holds the configuration used for the ruler storage cache.
	Cache cache.BackendConfig `yaml:"cache"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	c.Local.RegisterFlags(f)

	prefix := "templates.store."
	f.StringVar(&c.Cache.Backend, prefix+"cache.backend", "", fmt.Sprintf("Backend for template storage cache, if not empty. The cache is supported for any storage backend except %q. Supported values: %s.", local.Name, strings.Join(supportedCacheBackends, ", ")))
	c.Cache.Memcached.RegisterFlagsWithPrefix(prefix+"cache.memcached.", f)
	c.Cache.Redis.RegisterFlagsWithPrefix(prefix+"cache.redis.", f)
}

type TemplateStore interface {
	provider.TemplateProvider

	SyncTemplates()
	ListAppTemplates()
	LoadTemplate(path string) error
	UnloadTemplate(path string)
}
