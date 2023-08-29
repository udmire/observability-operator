package specs

import (
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	core_v1 "k8s.io/api/core/v1"
)

const (
	clusterNameEnv string = "K8S_CLUSTER_NAME"
)

func (h *appHandler) Decorate(manifest *manifest.AppManifests, decorators ...Decorator) {
	for _, decorator := range decorators {
		decorator(manifest)
	}
}

func ClusterNameEnvDecorator(valueFunc func() string) Decorator {
	containerEnvProcessor := func(c *core_v1.Container, value func() string) {
		for _, env := range c.Env {
			if env.Name == clusterNameEnv {
				break
			}
		}
		c.Env = append(c.Env, core_v1.EnvVar{
			Name:  clusterNameEnv,
			Value: value(),
		})
	}
	envProcessor := func(spec *core_v1.PodSpec, value func() string) {
		if spec == nil {
			return
		}

		var containers []core_v1.Container
		for _, container := range spec.Containers {
			containerEnvProcessor(&container, value)
			containers = append(containers, container)
		}
		spec.Containers = containers

		containers = make([]core_v1.Container, len(spec.InitContainers))
		for _, container := range spec.InitContainers {
			containerEnvProcessor(&container, value)
			containers = append(containers, container)
		}
		spec.InitContainers = containers
	}
	return func(manifest *manifest.AppManifests) {
		for _, comp := range manifest.CompsMenifests {
			if comp.Deployment != nil {
				envProcessor(&comp.Deployment.Spec.Template.Spec, valueFunc)
			}

			if comp.DaemonSet != nil {
				envProcessor(&comp.DaemonSet.Spec.Template.Spec, valueFunc)
			}

			if comp.StatefulSet != nil {
				envProcessor(&comp.StatefulSet.Spec.Template.Spec, valueFunc)
			}

			if comp.ReplicaSet != nil {
				envProcessor(&comp.ReplicaSet.Spec.Template.Spec, valueFunc)
			}

			if comp.Job != nil {
				envProcessor(&comp.Job.Spec.Template.Spec, valueFunc)
			}

			if comp.CronJob != nil {
				envProcessor(&comp.CronJob.Spec.JobTemplate.Spec.Template.Spec, valueFunc)
			}
		}
	}
}
