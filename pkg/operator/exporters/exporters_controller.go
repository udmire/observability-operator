/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package exporters

import (
	"context"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	app_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/apps/reconcile"
	"github.com/udmire/observability-operator/pkg/apps/specs"
)

const (
	defaultConcurrency = 3
)

// ExportersReconciler reconciles a Exporters object
type ExportersReconciler struct {
	client.Client
	Scheme *runtime.Scheme

	concurrency   int
	handler       specs.AppHandler
	appReconciler reconcile.AppReconciler
	logger        log.Logger
}

//+kubebuilder:rbac:groups=udmire.cn,resources=exporters,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=udmire.cn,resources=exporters/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=udmire.cn,resources=exporters/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Exporters object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *ExportersReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	level.Info(r.logger).Log("msg", "reconciling exporters")
	defer level.Info(r.logger).Log("msg", "done reconciling exporters")

	var instance v1alpha1.Exporters
	if err := r.Get(ctx, req.NamespacedName, &instance); apierrors.IsNotFound(err) {
		level.Error(r.logger).Log("msg", "detected deleted Exporters", "err", err)
		return ctrl.Result{}, nil
	} else if err != nil {
		level.Error(r.logger).Log("msg", "unable to get Exporters", "err", err)
		return ctrl.Result{}, nil
	}

	if r.concurrency <= 0 {
		r.concurrency = defaultConcurrency
	}

	r.normalizeExporters(&instance)
	owner := metav1.OwnerReference{
		APIVersion:         instance.APIVersion,
		BlockOwnerDeletion: pointer.Bool(true),
		Controller:         pointer.Bool(true),
		Kind:               instance.Kind,
		Name:               instance.Name,
		UID:                instance.UID,
	}

	var wg sync.WaitGroup
	semaphore := make(chan struct{}, r.concurrency)
	for _, exploy := range instance.Spec.Exployments {
		wg.Add(1)
		go func(app v1alpha1.AppSpec) {
			defer wg.Done()
			semaphore <- struct{}{}

			manifest, err := r.handler.Handle(app)
			if err != nil {
				level.Error(r.logger).Log("msg", "failed to generate manifests", "instance", instance.Name, "exporter", app.Name, "err", err)
				<-semaphore
				return
			}

			err = r.appReconciler.Reconcile(owner, "exporter", app.Name, manifest)
			if err != nil {
				level.Error(r.logger).Log("msg", "failed to apply manifests", "instance", instance.Name, "exporter", app.Name, "err", err)
			}

			<-semaphore
		}(exploy)
	}
	wg.Wait()

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *ExportersReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Exporters{}).
		Owns(&core_v1.ConfigMap{}).
		Owns(&core_v1.Secret{}).
		Owns(&core_v1.ServiceAccount{}).
		Owns(&core_v1.Service{}).
		Owns(&app_v1.Deployment{}).
		Owns(&app_v1.DaemonSet{}).
		Owns(&app_v1.StatefulSet{}).
		Owns(&app_v1.ReplicaSet{}).
		Owns(&batch_v1.Job{}).
		Owns(&batch_v1.CronJob{}).
		Owns(&networking_v1.Ingress{}).
		Owns(&rbac_v1.ClusterRole{}).
		Owns(&rbac_v1.ClusterRoleBinding{}).
		Owns(&rbac_v1.RoleBinding{}).
		Owns(&rbac_v1.Role{}).
		Owns(&autoscaling_v1.HorizontalPodAutoscaler{}).
		Complete(r)
}

func (r *ExportersReconciler) normalizeExporters(instance *v1alpha1.Exporters) {
	for name, exporter := range instance.Spec.Exployments {
		exporter.Name = name
		exporter.Namespace = instance.Namespace
	}
}
