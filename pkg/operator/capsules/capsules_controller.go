package capsules

import (
	"context"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/capsules/reconcile"
	"github.com/udmire/observability-operator/pkg/capsules/specs"
	"github.com/udmire/observability-operator/pkg/operator/manager"
	"github.com/udmire/observability-operator/pkg/templates/provider"
)

// CapsulesReconciler reconciles a Capsule object
type CapsulesReconciler struct {
	*services.BasicService

	client.Client
	Scheme *runtime.Scheme

	mgr ctrl.Manager
	cnp manager.ClusterNameProvider

	handler       specs.CapsuleHandler
	capReconciler reconcile.CapsuleReconciler
	logger        log.Logger
}

func New(client client.Client, schema *runtime.Scheme, tp provider.TemplateProvider, logger log.Logger) *CapsulesReconciler {
	reconciler := &CapsulesReconciler{
		Client: client,
		Scheme: schema,

		handler:       specs.New(tp, logger),
		capReconciler: reconcile.New(logger, client),
		logger:        logger,
	}
	reconciler.BasicService = services.NewIdleService(func(serviceContext context.Context) error {
		return reconciler.SetupWithManager(reconciler.mgr)
	}, nil)
	return reconciler
}

func (r *CapsulesReconciler) SetManager(mgr ctrl.Manager) {
	r.mgr = mgr
}

func (r *CapsulesReconciler) SetClusterNameProvider(cnp manager.ClusterNameProvider) {
	r.cnp = cnp
}

//+kubebuilder:rbac:groups=udmire.cn,resources=capsule,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=udmire.cn,resources=capsule/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=udmire.cn,resources=capsule/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the Apps object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.15.0/pkg/reconcile
func (r *CapsulesReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	level.Info(r.logger).Log("msg", "reconciling capsule")
	defer level.Info(r.logger).Log("msg", "done reconciling capsule")

	instance := v1alpha1.Capsule{}
	if err := r.Get(ctx, req.NamespacedName, &instance); apierrors.IsNotFound(err) {
		level.Error(r.logger).Log("msg", "detected deleted capsule", "err", err)
		return ctrl.Result{}, nil
	} else if err != nil {
		level.Error(r.logger).Log("msg", "unable to get capsule", "err", err)
		return ctrl.Result{}, nil
	}

	r.normalize(&instance)
	owner := metav1.OwnerReference{
		APIVersion:         instance.APIVersion,
		BlockOwnerDeletion: pointer.Bool(true),
		Controller:         pointer.Bool(true),
		Kind:               instance.Kind,
		Name:               instance.Name,
		UID:                instance.UID,
	}

	manifest, err := r.handler.Handle(instance.Spec)
	if err != nil {
		level.Error(r.logger).Log("msg", "failed to generate manifests", "instance", instance.Name, "capsule", instance.Spec.Name, "err", err)
	}

	err = r.capReconciler.Reconcile(ctx, owner, manifest)
	if err != nil {
		level.Error(r.logger).Log("msg", "failed to apply manifests", "instance", instance.Name, "capsule", instance.Spec.Name, "err", err)
	}

	return ctrl.Result{}, nil
}

// SetupWithManager sets up the controller with the Manager.
func (r *CapsulesReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&v1alpha1.Capsule{}).
		Complete(r)
}

func (r *CapsulesReconciler) normalize(instance *v1alpha1.Capsule) {
	instance.Spec.Namespace = instance.Namespace
}
