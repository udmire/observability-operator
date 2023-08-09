package reconcile

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppReconciler interface {
	Reconcile(owner metav1.OwnerReference, appType, name string, manifest *manifest.AppManifests) error
}

func New(logger log.Logger) AppReconciler {
	return &reconciler{
		logger: logger,
	}
}

type reconciler struct {
	logger log.Logger
	client client.Client
}

func (r *reconciler) Reconcile(owner metav1.OwnerReference, appType, name string, manifest *manifest.AppManifests) error {
	level.Info(r.logger).Log("msg", "start to reconcile", appType, name)
	cxt := context.Background()

	err := r.reconcile(cxt, owner, appType, name, &manifest.Manifests)
	if err != nil {
		level.Warn(r.logger).Log("msg", "reconcile manifests failed", appType, name, "err", err)
		return err
	}

	for _, component := range manifest.CompsMenifests {
		err = r.reconcileComponent(cxt, owner, appType, name, component)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed", appType, name, "comp", component.Name, "err", err)
			return err
		}
	}

	level.Info(r.logger).Log("msg", "reconcile success", appType, name)

	return nil
}

func (r *reconciler) reconcile(ctx context.Context, owner metav1.OwnerReference, appType, name string, manifest *manifest.Manifests) error {
	if manifest.ServiceAccount != nil {
		err := CreateOrUpdateServiceAccount(ctx, r.client, manifest.ServiceAccount)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create sa", appType, name, "err", err)
			return err
		}
	}

	if manifest.ClusterRole != nil {
		err := CreateOrUpdateClusterRole(ctx, r.client, manifest.ClusterRole)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create clusterRole", appType, name, "err", err)
			return err
		}
	}

	if manifest.ClusterRoleBinding != nil {
		err := CreateOrUpdateClusterRoleBinding(ctx, r.client, manifest.ClusterRoleBinding)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create clusterRoleBinding", appType, name, "err", err)
			return err
		}
	}

	if manifest.Role != nil {
		err := CreateOrUpdateRole(ctx, r.client, manifest.Role)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create role", appType, name, "err", err)
			return err
		}
	}

	if manifest.RoleBinding != nil {
		err := CreateOrUpdateRoleBinding(ctx, r.client, manifest.RoleBinding)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create roleBinding", appType, name, "err", err)
			return err
		}
	}

	if manifest.Ingress != nil {
		err := CreateOrUpdateIngress(ctx, r.client, manifest.Ingress)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create ingress", appType, name, "err", err)
			return err
		}
	}

	for _, secret := range manifest.Secrets {
		err := CreateOrUpdateSecret(ctx, r.client, secret)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create secret", appType, name, "err", err)
			return err
		}
	}

	for _, cm := range manifest.ConfigMaps {
		err := CreateOrUpdateConfigMap(ctx, r.client, cm)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create configmap", appType, name, "err", err)
			return err
		}
	}

	for _, svc := range manifest.Services {
		err := CreateOrUpdateService(ctx, r.client, svc)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create service", appType, name, "err", err)
			return err
		}
	}

	return nil
}

func (r *reconciler) reconcileComponent(ctx context.Context, owner metav1.OwnerReference, appType, name string, manifest *manifest.CompManifests) error {
	err := r.reconcile(ctx, owner, appType, name, &manifest.Manifests)
	if err != nil {
		level.Warn(r.logger).Log("msg", "reconcile manifests failed to create component", appType, name, "component", manifest.Name, "err", err)
		return err
	}

	if manifest.Deployment != nil {
		err := CreateOrUpdateDeployment(ctx, r.client, manifest.Deployment)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create deployment workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.DaemonSet != nil {
		err := CreateOrUpdateDaemonSet(ctx, r.client, manifest.DaemonSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create daemonset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.StatefulSet != nil {
		err := CreateOrUpdateStatefulSet(ctx, r.client, manifest.StatefulSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create statefulset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.ReplicaSet != nil {
		err := CreateOrUpdateReplicaSet(ctx, r.client, manifest.ReplicaSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create replicaset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.Job != nil {
		err := CreateOrUpdateJob(ctx, r.client, manifest.Job)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create job workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.CronJob != nil {
		err := CreateOrUpdateCronJob(ctx, r.client, manifest.CronJob)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create cronjob workload", appType, name, "err", err)
			return err
		}
	}

	return nil
}
