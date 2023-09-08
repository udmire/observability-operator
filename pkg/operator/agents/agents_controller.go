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
	v1 "k8s.io/api/core/v1"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/apps/reconcile"
	"github.com/udmire/observability-operator/pkg/apps/specs"
	"github.com/udmire/observability-operator/pkg/operator/base"
	"github.com/udmire/observability-operator/pkg/operator/manager"
	"github.com/udmire/observability-operator/pkg/templates/provider"
)

// AgentsReconciler reconciles a Agents object
type AgentsReconciler struct {
	base.BaseReconciler

	mgr ctrl.Manager
	cnp manager.ClusterNameProvider

	handler       specs.AppHandler
	appReconciler reconcile.AppReconciler
}

func New(client client.Client, schema *runtime.Scheme, tp provider.TemplateProvider, logger log.Logger) *AgentsReconciler {
	reconciler := &AgentsReconciler{
		BaseReconciler: base.BaseReconciler{
			Client: client,
			Scheme: schema,
			Logger: logger,
		},

		handler:       specs.New(tp, logger),
		appReconciler: reconcile.New(logger, client),
	}
	reconciler.BasicService = services.NewIdleService(func(serviceContext context.Context) error {
		return reconciler.SetupWithManager(reconciler.mgr)
	}, nil)
	return reconciler
}

func (r *AgentsReconciler) SetManager(mgr ctrl.Manager) {
	r.mgr = mgr
}

func (r *AgentsReconciler) SetClusterNameProvider(cnp manager.ClusterNameProvider) {
	r.cnp = cnp
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
	level.Info(r.Logger).Log("msg", "reconciling agent")
	defer level.Info(r.Logger).Log("msg", "done reconciling agent")

	instance := &v1alpha1.Agents{}
	if err := r.Get(ctx, req.NamespacedName, instance); apierrors.IsNotFound(err) {
		level.Error(r.Logger).Log("msg", "detected deleted Agents", "err", err)
		return ctrl.Result{}, nil
	} else if err != nil {
		level.Error(r.Logger).Log("msg", "unable to get Agents", "err", err)
		return ctrl.Result{}, nil
	}

	r.normalizeInstance(instance)

	finalizerName := "agents.udmire.cn/finalizer"
	if instance.ObjectMeta.DeletionTimestamp.IsZero() {
		if !controllerutil.ContainsFinalizer(instance, finalizerName) {
			controllerutil.AddFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}
	} else {
		// The object is being deleted
		if controllerutil.ContainsFinalizer(instance, finalizerName) {
			selector := r.handler.Selector(instance.Spec.AppSpec)
			if err := r.appReconciler.CleanClusterLayerResources(instance.UID, selector); err != nil {
				return ctrl.Result{}, err
			}

			// remove our finalizer from the list and update it.
			controllerutil.RemoveFinalizer(instance, finalizerName)
			if err := r.Update(ctx, instance); err != nil {
				return ctrl.Result{}, err
			}
		}

		// Stop reconciliation as the item is being deleted
		return ctrl.Result{}, nil
	}

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
		level.Error(r.Logger).Log("msg", "failed to generate manifests", "instance", instance.Name, "err", err)
		return ctrl.Result{}, err
	}

	err = r.ProcessDependencies(owner, instance.Namespace, instance.Spec.Template, instance.Spec.Singleton, instance.Spec.Dependencies)
	if err != nil {
		level.Error(r.Logger).Log("msg", "failed to create dependencies", "instance", instance.Name, "err", err)
		return ctrl.Result{}, err
	}

	r.handler.Decorate(manifest, specs.ClusterNameEnvDecorator(r.cnp))

	err = r.appReconciler.Reconcile(owner, "agents", instance.Name, manifest)
	if err != nil {
		level.Error(r.Logger).Log("msg", "failed to apply manifests", "instance", instance.Name, "err", err)
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *AgentsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Agents{}).
		Owns(&v1alpha1.Exporters{}).
		Owns(&v1alpha1.Apps{}).
		Owns(&v1alpha1.Capsule{}).
		Complete(r)
}

func (r *AgentsReconciler) normalizeInstance(instance *v1alpha1.Agents) {
	instance.Spec.Name = instance.Name
	instance.Spec.Namespace = instance.Namespace

	for _, comp := range instance.Spec.Components {
		if comp.DaemonSet != nil {
			r.updatePodTemplateEnv(comp.DaemonSet.Template)
		}
		if comp.Deployment != nil {
			r.updatePodTemplateEnv(comp.Deployment.Template)
		}
		if comp.StatefulSet != nil {
			r.updatePodTemplateEnv(comp.StatefulSet.Template)
		}
	}
}
func (r *AgentsReconciler) updatePodTemplateEnv(spec *v1alpha1.PodTemplateSpec) {
	const clusterNameEnv = "K8S_CLUSTER_NAME"
	for _, container := range spec.Spec.InitContainers {
		for _, env := range container.Env {
			if env.Name == clusterNameEnv {
				return
			}
		}
		container.Env = append(container.Env, v1.EnvVar{
			Name:  clusterNameEnv,
			Value: r.cnp(),
		})
	}

	for _, container := range spec.Spec.Containers {
		for _, env := range container.Env {
			if env.Name == clusterNameEnv {
				return
			}
		}
		container.Env = append(container.Env, v1.EnvVar{
			Name:  clusterNameEnv,
			Value: r.cnp(),
		})
	}
}
