package templates

import "github.com/udmire/observability-operator/pkg/templates/store"

type Config struct {
	Store store.Config `yaml:"store"`
}
