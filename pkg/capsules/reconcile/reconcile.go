package reconcile

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/capsules/manifest"
	util_client "github.com/udmire/observability-operator/pkg/utils/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type CapsuleReconciler interface {
	Reconcile(ctx context.Context, owner metav1.OwnerReference, manifest *manifest.CapsuleManifests) error
}

type capsuleReconciler struct {
	logger log.Logger
	client client.Client
}

func New(logger log.Logger, client client.Client) CapsuleReconciler {
	return &capsuleReconciler{
		logger: logger,
		client: client,
	}
}

func (r *capsuleReconciler) Reconcile(ctx context.Context, owner metav1.OwnerReference, manifest *manifest.CapsuleManifests) error {
	for _, secret := range manifest.Secrets {
		secret.OwnerReferences = append(secret.OwnerReferences, owner)
		err := util_client.CreateOrUpdateSecret(ctx, r.client, secret)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create secret", "name", secret.Name, "err", err)
			return err
		}
	}

	for _, cm := range manifest.ConfigMaps {
		cm.OwnerReferences = append(cm.OwnerReferences, owner)
		err := util_client.CreateOrUpdateConfigMap(ctx, r.client, cm)
		if err != nil {
			level.Warn(r.logger).Log("msg", "reconcile manifests failed to create configmap", "name", cm.Name, "err", err)
			return err
		}
	}

	return nil
}
