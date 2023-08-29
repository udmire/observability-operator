package manager

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/yaml"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type clusterInfoProvider struct {
	*services.BasicService

	cli client.Client

	cfg    *Config
	logger log.Logger
}

func newClusterInfoProvider(cfg *Config, cli client.Client, logger log.Logger) *clusterInfoProvider {
	provider := &clusterInfoProvider{
		cfg:    cfg,
		logger: logger,
		cli:    cli,
	}

	provider.BasicService = services.NewBasicService(nil, provider.tryUpdateClusterName, nil)

	return provider
}
func (w *clusterInfoProvider) tryUpdateClusterName(_ context.Context) error {
	if len(w.cfg.ClusterName) > 0 {
		return nil
	}

	// try read name from env
	w.cfg.ClusterName = os.Getenv("KUBERNETES_CLUSTER_NAME")
	if len(w.cfg.ClusterName) > 0 {
		return nil
	}

	// try read name from kubeadm-config
	w.tryUpdateClusterNameFromKubeadmConfig()
	if len(w.cfg.ClusterName) > 0 {
		return nil
	}
	return nil
}

func (w *clusterInfoProvider) tryUpdateClusterNameFromKubeadmConfig() {

	cm := &core_v1.ConfigMap{}
	err := w.cli.Get(context.Background(), client.ObjectKey{Namespace: "kube-system", Name: "kubeadm-config"}, cm)
	if err != nil {
		_, ok := err.(*cache.ErrCacheNotStarted)
		for ok {
			level.Warn(w.logger).Log("msg", err.Error())
			time.Sleep(0)
			err = w.cli.Get(context.Background(), client.ObjectKey{Namespace: "kube-system", Name: "kubeadm-config"}, cm)
			_, ok = err.(*cache.ErrCacheNotStarted)
		}
		if !ok && err != nil {
			level.Warn(w.logger).Log("msg", "failed to load kubeadm-config")
			return
		}
	}
	cc := cm.Data["ClusterConfiguration"]

	type clusterConfig struct {
		ClusterName string `yaml:"clusterName"`
	}

	config := &clusterConfig{}

	err = yaml.Unmarshal([]byte(cc), config)
	if err != nil {
		level.Warn(w.logger).Log("msg", "failed to unmarshal kubeadm-config")
		return
	}

	w.cfg.ClusterName = config.ClusterName
}
