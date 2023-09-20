package providers

import (
	"context"
	"os"
	"time"

	"github.com/go-kit/log/level"
	"gopkg.in/yaml.v2"
	core_v1 "k8s.io/api/core/v1"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (w *providers) ClusterNameProvider() StringProvider {
	return func() string {
		return w.holders.clusterName
	}
}

func (w *providers) tryUpdateClusterName(_ context.Context) error {
	if len(w.holders.clusterName) > 0 {
		return nil
	}

	// try read name from env
	w.holders.clusterName = os.Getenv("KUBERNETES_CLUSTER_NAME")
	if len(w.holders.clusterName) > 0 {
		return nil
	}

	// try read name from kubeadm-config
	w.tryUpdateClusterNameFromKubeadmConfig()
	if len(w.holders.clusterName) > 0 {
		return nil
	}
	return nil
}

func (w *providers) tryUpdateClusterNameFromKubeadmConfig() {
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

	w.holders.clusterName = config.ClusterName
}
