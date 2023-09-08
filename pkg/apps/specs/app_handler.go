package specs

import (
	"fmt"

	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	"github.com/udmire/observability-operator/pkg/templates/template"
	"github.com/udmire/observability-operator/pkg/utils"
)

type Decorator func(manifest *manifest.AppManifests)

type AppHandler interface {
	Handle(app v1alpha1.AppSpec) (*manifest.AppManifests, error)
	Selector(app v1alpha1.AppSpec) labels.Selector
	Decorate(manifest *manifest.AppManifests, decorators ...Decorator)
}

type appHandler struct {
	logger log.Logger

	provider provider.TemplateProvider
}

func New(provider provider.TemplateProvider, logger log.Logger) AppHandler {
	return &appHandler{
		logger:   logger,
		provider: provider,
	}
}

func (h *appHandler) Handle(app v1alpha1.AppSpec) (*manifest.AppManifests, error) {
	template := h.getTemplate(app.Template.Name, app.Template.Version)
	if template == nil {
		version := "latest"
		if len(app.Template.Version) > 0 {
			version = app.Template.Version
		}
		level.Warn(h.logger).Log("msg", "template not found", "name", app.Template.Name, "version", version)
		return nil, fmt.Errorf("template %s:%s not found", app.Template.Name, version)
	}

	manifest := manifest.NewTemplateBuilder(template).Build()
	return h.customerizeApp(manifest, app)
}

func (h *appHandler) Selector(app v1alpha1.AppSpec) labels.Selector {
	instanceLabels := utils.AppInstanceLabels(app.Name, app.Template.Name, app.Template.Version)
	delete(instanceLabels, utils.AppLabel)
	return labels.SelectorFromSet(labels.Set(instanceLabels))
}

func (h *appHandler) customerizeApp(manifest *manifest.AppManifests, app v1alpha1.AppSpec) (*manifest.AppManifests, error) {
	prefix := ""
	if !app.Singleton {
		prefix = fmt.Sprintf("%s-", app.Name)
	}
	instanceLabels := utils.AppInstanceLabels(app.Name, app.Template.Name, app.Template.Version)
	namespace := app.Namespace

	h.customerize(&manifest.Manifests, app.CommonSpec, prefix, app.Name, namespace, instanceLabels)

	if len(manifest.CompsMenifests) > 0 {
		for _, component := range manifest.CompsMenifests {
			componentName := component.Name
			var err error
			if compSpec, exists := app.Components[componentName]; exists {
				err = h.customerizeComponent(component, compSpec, app.Template, prefix, app.Name, componentName, namespace)
			} else {
				err = h.customerizeComponent(component, v1alpha1.ComponentSpec{}, app.Template, prefix, app.Name, componentName, namespace)
			}

			if err != nil {
				level.Error(h.logger).Log("msg", "failed to customerize component", "name", componentName, "err", err)
				return nil, err
			}
		}
	}

	return manifest, nil
}

func (h *appHandler) customerizeComponent(manifest *manifest.CompManifests, spec v1alpha1.ComponentSpec, template v1alpha1.Template, prefix, instance, component, namespace string) error {
	compLabels := utils.ComponentLabels(instance, template.Name, template.Version, component)

	h.customerize(&manifest.Manifests, spec.CommonSpec, prefix, component, namespace, compLabels)

	var err error
	if manifest.Deployment != nil {
		err = mergeDeployment(manifest.Deployment, spec.Deployment, prefix, namespace, compLabels)
	}
	if manifest.DaemonSet != nil {
		err = mergeDaemonset(manifest.DaemonSet, spec.DaemonSet, prefix, namespace, compLabels)
	}
	if manifest.StatefulSet != nil {
		err = mergeStatefulSet(manifest.StatefulSet, spec.StatefulSet, prefix, namespace, compLabels)
	}
	if manifest.ReplicaSet != nil {
		err = mergeReplicaSet(manifest.ReplicaSet, spec.ReplicaSet, prefix, namespace, compLabels)
	}
	if manifest.Job != nil {
		err = mergeJob(manifest.Job, spec.Job, prefix, namespace, compLabels)
	}
	if manifest.CronJob != nil {
		err = mergeCronJob(manifest.CronJob, spec.CronJob, prefix, namespace, compLabels)
	}

	return err
}

func (h *appHandler) getTemplate(name, version string) *template.AppTemplate {
	var template *template.AppTemplate
	if len(version) == 0 {
		template = h.provider.GetLatestTemplate(name)
	} else {
		template = h.provider.GetTemplate(name, version)
	}

	return template
}

func (h *appHandler) customerize(manifest *manifest.Manifests, app v1alpha1.CommonSpec, prefix, name, namespace string, labels map[string]string) error {
	configMapsCustom(manifest, app.ConfigMaps, prefix, namespace, labels)
	secretsCustom(manifest, app.Secrets, prefix, namespace, labels)
	servicesCustom(manifest, app.Services, prefix, namespace, labels)

	if manifest.ServiceAccount != nil {
		mergeServiceAccount(manifest.ServiceAccount, app.ServiceAccount, prefix, namespace, labels)
	} else if app.ServiceAccount != nil {
		manifest.ServiceAccount = &core_v1.ServiceAccount{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Secrets:                      app.ServiceAccount.Secrets,
			AutomountServiceAccountToken: app.ServiceAccount.AutomountServiceAccountToken,
			ImagePullSecrets:             app.ServiceAccount.ImagePullSecrets,
		}
	}

	if manifest.ClusterRole != nil {
		mergeClusterRole(manifest.ClusterRole, app.ClusterRole, prefix, namespace, labels)
	} else if app.ClusterRole != nil {
		manifest.ClusterRole = &rbac_v1.ClusterRole{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Rules:           app.ClusterRole.Rules,
			AggregationRule: app.ClusterRole.AggregationRule,
		}
	}

	if manifest.ClusterRoleBinding != nil {
		mergeClusterRoleBinding(manifest.ClusterRoleBinding, app.ClusterRoleBinding, prefix, namespace, labels)
	} else if app.ClusterRoleBinding != nil {
		manifest.ClusterRoleBinding = &rbac_v1.ClusterRoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Subjects: app.ClusterRoleBinding.Subjects,
			RoleRef:  *app.ClusterRoleBinding.RoleRef,
		}
	}

	if manifest.Role != nil {
		mergeRole(manifest.Role, app.Role, prefix, namespace, labels)
	} else if app.Role != nil {
		manifest.Role = &rbac_v1.Role{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Rules: app.Role.Rules,
		}
	}

	if manifest.RoleBinding != nil {
		mergeRoleBinding(manifest.RoleBinding, app.RoleBinding, prefix, namespace, labels)
	} else if app.RoleBinding != nil {
		manifest.RoleBinding = &rbac_v1.RoleBinding{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Subjects: app.RoleBinding.Subjects,
			RoleRef:  *app.RoleBinding.RoleRef,
		}
	}

	if manifest.Ingress != nil {
		mergeIngress(manifest.Ingress, app.Ingress, prefix, namespace, labels)
	} else if app.Ingress != nil {
		manifest.Ingress = &networking_v1.Ingress{
			ObjectMeta: metav1.ObjectMeta{
				Name:      fmt.Sprintf("%s-%s", prefix, name),
				Namespace: namespace,
				Labels:    labels,
			},
			Spec: networking_v1.IngressSpec{
				Rules:            app.Ingress.Rules,
				IngressClassName: app.Ingress.IngressClassName,
				DefaultBackend:   app.Ingress.DefaultBackend,
				TLS:              app.Ingress.TLS,
			},
		}
	}

	return nil
}
