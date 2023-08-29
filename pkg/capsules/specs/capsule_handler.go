package specs

import (
	"fmt"

	core_v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/capsules/manifest"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	"github.com/udmire/observability-operator/pkg/templates/template"
	"github.com/udmire/observability-operator/pkg/utils"
)

type CapsuleHandler interface {
	Handle(app v1alpha1.CapsuleSpec) (*manifest.CapsuleManifests, error)
}

type capsuleHandler struct {
	logger log.Logger

	provider provider.TemplateProvider
}

func New(provider provider.TemplateProvider, logger log.Logger) CapsuleHandler {
	return &capsuleHandler{
		logger:   logger,
		provider: provider,
	}
}

func (h *capsuleHandler) Handle(capsule v1alpha1.CapsuleSpec) (*manifest.CapsuleManifests, error) {
	template := h.getTemplate(capsule.Template.Name, capsule.Template.Version)
	if template == nil {
		version := "latest"
		if len(capsule.Template.Version) > 0 {
			version = capsule.Template.Version
		}
		level.Warn(h.logger).Log("msg", "template not found", "name", capsule.Template.Name, "version", version)
		return nil, fmt.Errorf("template %s:%s not found", capsule.Template.Name, version)
	}

	manifest := manifest.New(template).Build()
	return h.customerizeApp(manifest, capsule)
}

func (h *capsuleHandler) getTemplate(name, version string) *template.AppTemplate {
	var template *template.AppTemplate
	if len(version) == 0 {
		template = h.provider.GetLatestTemplate(name)
	} else {
		template = h.provider.GetTemplate(name, version)
	}

	return template
}

func (h *capsuleHandler) customerizeApp(manifest *manifest.CapsuleManifests, app v1alpha1.CapsuleSpec) (*manifest.CapsuleManifests, error) {
	instanceLabels := utils.AppInstanceLabels(app.Name, app.Template.Name, app.Template.Version)
	namespace := app.Namespace

	h.customerize(&manifest.Manifest, app.CapsuleCommonSpec, app.Name, namespace, instanceLabels)

	if len(manifest.CompsManifests) > 0 {
		for _, component := range manifest.CompsManifests {
			componentName := component.Name
			var err error
			if compSpec, exists := app.Components[componentName]; exists {
				err = h.customerizeComponent(component, compSpec, app.Name, app.Template.Name, app.Template.Version, componentName, namespace)
			} else {
				err = h.customerizeComponent(component, v1alpha1.CapsuleCommonSpec{}, app.Name, app.Template.Name, app.Template.Version, componentName, namespace)
			}

			if err != nil {
				level.Error(h.logger).Log("msg", "failed to customerize component", "name", componentName, "err", err)
				return nil, err
			}
		}
	}

	return manifest, nil
}

func (h *capsuleHandler) customerizeComponent(manifest *manifest.CompManifests, spec v1alpha1.CapsuleCommonSpec, instance, template, version, component, namespace string) error {
	compLabels := utils.ComponentLabels(instance, template, version, component)

	return h.customerize(&manifest.Manifest, spec, component, namespace, compLabels)
}

func (h *capsuleHandler) customerize(manifest *manifest.Manifest, app v1alpha1.CapsuleCommonSpec, name, namespace string, labels map[string]string) error {
	configMapsCustom(manifest, app.ConfigMaps, namespace, labels)
	secretsCustom(manifest, app.Secrets, namespace, labels)

	return nil
}

func configMapsCustom(manifest *manifest.Manifest, configmaps map[string]*v1alpha1.ConfigMapSpec, ns string, labels map[string]string) {
	for _, cm := range manifest.ConfigMaps {
		mergeConfigMap(cm, nil, ns, labels)
	}

	if len(configmaps) < 1 {
		return
	}

	merged := make(map[string]string)
	for _, cm := range manifest.ConfigMaps {
		if configmap, ok := configmaps[cm.Name]; ok {
			mergeConfigMap(cm, configmap, ns, labels)
			merged[cm.Name] = ""
		}
		continue
	}

	for name, cm := range configmaps {
		if _, ok := merged[name]; ok {
			continue
		}
		manifest.ConfigMaps = append(manifest.ConfigMaps, &core_v1.ConfigMap{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
				Labels:    labels,
			},
			Data: cm.Data,
		})
	}
}

func mergeConfigMap(manifest *core_v1.ConfigMap, configmap *v1alpha1.ConfigMapSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if configmap == nil {
		return
	}

	if manifest.Immutable != nil && *manifest.Immutable {
		return
	}

	if len(configmap.Data) < 1 {
		return
	}

	manifest.Data = configmap.Data
}

func secretsCustom(manifest *manifest.Manifest, secrets map[string]*v1alpha1.SecretSpec, ns string, labels map[string]string) {
	for _, sec := range manifest.Secrets {
		mergeSecret(sec, nil, ns, labels)
	}

	if len(secrets) < 1 {
		return
	}

	merged := make(map[string]string)
	for _, sec := range manifest.Secrets {
		if secret, ok := secrets[sec.Name]; ok {
			mergeSecret(sec, secret, ns, labels)
			merged[sec.Name] = ""
		}
		continue
	}

	for name, sec := range secrets {
		if _, ok := merged[name]; ok {
			continue
		}
		manifest.Secrets = append(manifest.Secrets, &core_v1.Secret{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
				Labels:    labels,
			},
			StringData: sec.StringData,
		})
	}
}

func mergeSecret(manifest *core_v1.Secret, secret *v1alpha1.SecretSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if secret == nil {
		return
	}

	if manifest.Immutable != nil && *manifest.Immutable {
		return
	}

	if len(secret.StringData) < 1 {
		return
	}

	manifest.StringData = secret.StringData
}

func mergeObjectMeta(meta *metav1.ObjectMeta, ns string, labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	if len(ns) > 0 {
		meta.Namespace = ns
	}

	for key, value := range labels {
		if _, ok := meta.Labels[key]; ok {
			continue
		}
		meta.Labels[key] = value
	}
}
