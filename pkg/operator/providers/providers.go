package providers

import (
	"github.com/go-kit/log"
	"github.com/grafana/dskit/services"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type StringProvider func() string

type Providers interface {
	ClusterNameProvider() StringProvider
}

type infoHolder struct {
	clusterName string
}

type providers struct {
	*services.BasicService

	holders *infoHolder

	cli    client.Client
	logger log.Logger
}

func NewProviders(cli client.Client, logger log.Logger) *providers {
	provider := &providers{
		logger:  logger,
		cli:     cli,
		holders: &infoHolder{},
	}

	provider.BasicService = services.NewIdleService(provider.tryUpdateClusterName, nil)

	return provider
}
