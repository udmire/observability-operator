package reconcile

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	util_client "github.com/udmire/observability-operator/pkg/utils/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type AppReconciler interface {
	Reconcile(owner metav1.OwnerReference, appType, name string, manifest *manifest.AppManifests) error
	CleanClusterLayerResources(uid types.UID, selector labels.Selector) error
}

func New(logger log.Logger, client client.Client) AppReconciler {
	return &reconciler{
		logger: logger,
		client: client,
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

func (r *reconciler) CleanClusterLayerResources(uid types.UID, selector labels.Selector) error {
	level.Info(r.logger).Log("msg", "start to clean cluster layer resources")
	ctx := context.Background()
	err := util_client.CleanClusterRoleBindings(ctx, r.client, uid, selector)
	if err != nil {
		return err
	}
	err = util_client.CleanClusterRoles(ctx, r.client, uid, selector)
	if err != nil {
		return err
	}
	level.Info(r.logger).Log("msg", "clean cluster layer resources success")
	return nil
}

func (r *reconciler) reconcile(ctx context.Context, owner metav1.OwnerReference, appType, name string, manifest *manifest.Manifests) error {
	if manifest.ServiceAccount != nil {
		manifest.ServiceAccount.OwnerReferences = append(manifest.ServiceAccount.OwnerReferences, owner)
		err := util_client.CreateOrUpdateServiceAccount(ctx, r.client, manifest.ServiceAccount)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create sa", appType, name, "err", err)
			return err
		}
	}

	if manifest.ClusterRole != nil {
		manifest.ClusterRole.OwnerReferences = append(manifest.ClusterRole.OwnerReferences, owner)
		err := util_client.CreateOrUpdateClusterRole(ctx, r.client, manifest.ClusterRole)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create clusterRole", appType, name, "err", err)
			return err
		}
	}

	if manifest.ClusterRoleBinding != nil {
		manifest.ClusterRole.OwnerReferences = append(manifest.ClusterRole.OwnerReferences, owner)
		err := util_client.CreateOrUpdateClusterRoleBinding(ctx, r.client, manifest.ClusterRoleBinding)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create clusterRoleBinding", appType, name, "err", err)
			return err
		}
	}

	if manifest.Role != nil {
		manifest.Role.OwnerReferences = append(manifest.Role.OwnerReferences, owner)
		err := util_client.CreateOrUpdateRole(ctx, r.client, manifest.Role)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create role", appType, name, "err", err)
			return err
		}
	}

	if manifest.RoleBinding != nil {
		manifest.RoleBinding.OwnerReferences = append(manifest.RoleBinding.OwnerReferences, owner)
		err := util_client.CreateOrUpdateRoleBinding(ctx, r.client, manifest.RoleBinding)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create roleBinding", appType, name, "err", err)
			return err
		}
	}

	if manifest.Ingress != nil {
		manifest.Ingress.OwnerReferences = append(manifest.Ingress.OwnerReferences, owner)
		err := util_client.CreateOrUpdateIngress(ctx, r.client, manifest.Ingress)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create ingress", appType, name, "err", err)
			return err
		}
	}

	for _, secret := range manifest.Secrets {
		secret.OwnerReferences = append(secret.OwnerReferences, owner)
		err := util_client.CreateOrUpdateSecret(ctx, r.client, secret)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create secret", appType, name, "err", err)
			return err
		}
	}

	for _, cm := range manifest.ConfigMaps {
		cm.OwnerReferences = append(cm.OwnerReferences, owner)
		err := util_client.CreateOrUpdateConfigMap(ctx, r.client, cm)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create configmap", appType, name, "err", err)
			return err
		}
	}

	for _, svc := range manifest.Services {
		svc.OwnerReferences = append(svc.OwnerReferences, owner)
		err := util_client.CreateOrUpdateService(ctx, r.client, svc)
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
		manifest.Deployment.OwnerReferences = append(manifest.Deployment.OwnerReferences, owner)
		err := util_client.CreateOrUpdateDeployment(ctx, r.client, manifest.Deployment)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create deployment workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.DaemonSet != nil {
		manifest.DaemonSet.OwnerReferences = append(manifest.DaemonSet.OwnerReferences, owner)
		err := util_client.CreateOrUpdateDaemonSet(ctx, r.client, manifest.DaemonSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create daemonset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.StatefulSet != nil {
		manifest.StatefulSet.OwnerReferences = append(manifest.StatefulSet.OwnerReferences, owner)
		err := util_client.CreateOrUpdateStatefulSet(ctx, r.client, manifest.StatefulSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create statefulset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.ReplicaSet != nil {
		manifest.ReplicaSet.OwnerReferences = append(manifest.ReplicaSet.OwnerReferences, owner)
		err := util_client.CreateOrUpdateReplicaSet(ctx, r.client, manifest.ReplicaSet)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create replicaset workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.Job != nil {
		manifest.Job.OwnerReferences = append(manifest.Job.OwnerReferences, owner)
		err := util_client.CreateOrUpdateJob(ctx, r.client, manifest.Job)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create job workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.CronJob != nil {
		manifest.ClusterRole.OwnerReferences = append(manifest.ClusterRole.OwnerReferences, owner)
		err := util_client.CreateOrUpdateCronJob(ctx, r.client, manifest.CronJob)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create cronjob workload", appType, name, "err", err)
			return err
		}
	}

	if manifest.HPA != nil {
		manifest.HPA.OwnerReferences = append(manifest.HPA.OwnerReferences, owner)
		err := util_client.CreateOrUpdateHPA(ctx, r.client, manifest.HPA)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create hpa", appType, name, "err", err)
			return err
		}
	}

	return nil
}
