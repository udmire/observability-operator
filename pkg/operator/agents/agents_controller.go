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

package agents

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
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
	"github.com/udmire/observability-operator/pkg/apps/templates/provider"
)

// AgentsReconciler reconciles a Agents object
type AgentsReconciler struct {
	*services.BasicService

	client.Client
	Scheme *runtime.Scheme

	mgr ctrl.Manager

	handler       specs.AppHandler
	appReconciler reconcile.AppReconciler
	logger        log.Logger
}

func New(client client.Client, schema *runtime.Scheme, tp provider.TemplateProvider, logger log.Logger) *AgentsReconciler {
	reconciler := &AgentsReconciler{
		Client: client,
		Scheme: schema,

		handler:       specs.New(tp, logger),
		appReconciler: reconcile.New(logger),
		logger:        logger,
	}
	reconciler.BasicService = services.NewIdleService(func(serviceContext context.Context) error {
		return reconciler.SetupWithManager(reconciler.mgr)
	}, nil)
	return reconciler
}

func (r *AgentsReconciler) SetManager(mgr ctrl.Manager) {
	r.mgr = mgr
}

//+kubebuilder:rbac:groups=udmire.cn,resources=agents,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=udmire.cn,resources=agents/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=udmire.cn,resources=agents/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Agents object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.14.4/pkg/reconcile
func (r *AgentsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	level.Info(r.logger).Log("msg", "reconciling agent")
	defer level.Info(r.logger).Log("msg", "done reconciling agent")

	var instance v1alpha1.Agents
	if err := r.Get(ctx, req.NamespacedName, &instance); apierrors.IsNotFound(err) {
		level.Error(r.logger).Log("msg", "detected deleted Agents", "err", err)
		return ctrl.Result{}, nil
	} else if err != nil {
		level.Error(r.logger).Log("msg", "unable to get Agents", "err", err)
		return ctrl.Result{}, nil
	}

	r.normalizeInstance(&instance)
	owner := metav1.OwnerReference{
		APIVersion:         instance.APIVersion,
		BlockOwnerDeletion: pointer.Bool(true),
		Controller:         pointer.Bool(true),
		Kind:               instance.Kind,
		Name:               instance.Name,
		UID:                instance.UID,
	}

	manifest, err := r.handler.Handle(instance.Spec.AppSpec)
	if err != nil {
		level.Error(r.logger).Log("msg", "failed to generate manifests", "instance", instance.Name, "err", err)
		return ctrl.Result{}, err
	}

	err = r.appReconciler.Reconcile(owner, "agents", instance.Name, manifest)
	if err != nil {
		level.Error(r.logger).Log("msg", "failed to apply manifests", "instance", instance.Name, "err", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Agents{}).
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

func (r *AgentsReconciler) normalizeInstance(instance *v1alpha1.Agents) {
	instance.Spec.Name = instance.Name
	instance.Spec.Namespace = instance.Namespace
}
