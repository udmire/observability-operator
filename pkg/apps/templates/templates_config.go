package templates

import "github.com/udmire/observability-operator/pkg/apps/templates/store"

type Config struct {
	Store store.Config `yaml:"store"`
}
