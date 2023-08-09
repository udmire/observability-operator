package store

import (
	"github.com/grafana/dskit/cache"
	"github.com/udmire/observability-operator/pkg/apps/templates/provider"
	"github.com/udmire/observability-operator/pkg/apps/templates/store/local"
	"github.com/udmire/observability-operator/pkg/apps/templates/store/remote"
)

// Config configures a rule store.
type Config struct {
	remote.Config `yaml:",inline"`
	Local         local.Config `yaml:"local"`

	// Cache holds the configuration used for the ruler storage cache.
	Cache cache.BackendConfig `yaml:"cache"`
}

type TemplateStore interface {
	provider.TemplateProvider

	SyncTemplates()
	ListAppTemplates()
	LoadTemplate(path string)
	UnloadTemplate(path string)
}
