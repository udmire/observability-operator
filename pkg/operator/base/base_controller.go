package base

import (
	"context"
	"fmt"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/utils"
	util_client "github.com/udmire/observability-operator/pkg/utils/client"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type BaseReconciler struct {
	*services.BasicService

	client.Client
	Scheme *runtime.Scheme
	Logger log.Logger
}

func (r *BaseReconciler) ProcessDependencies(owner metav1.OwnerReference, ns string, template v1alpha1.Template, singleton bool, dep v1alpha1.AppDepsSpec) error {
	instanceLabels := utils.AppInstanceLabels(owner.Name, template.Name, template.Version)
	ctx := context.Background()

	for name, capsuleSpec := range dep.Capsules {
		if !singleton {
			name = fmt.Sprintf("%s-%s", owner.Name, name)
		}
		capsule := v1alpha1.Capsule{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
				OwnerReferences: []metav1.OwnerReference{
					owner,
				},
				Labels: instanceLabels,
			},
			Spec: capsuleSpec,
		}
		level.Info(r.Logger).Log("msg", "start to create dependency", "instance", owner.Name, "type", "capsule", "name", name)
		if err := util_client.CreateOrUpdateCapsule(ctx, r.Client, &capsule); err != nil {
			level.Warn(r.Logger).Log("msg", "failed to create dependency", "instance", owner.Name, "type", "capsule", "name", name, "err", err)
			return err
		}
		level.Info(r.Logger).Log("msg", "create dependency success", "instance", owner.Name, "type", "capsule", "name", name)
	}

	return nil
}
