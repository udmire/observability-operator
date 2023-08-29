package specs

import (
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"github.com/udmire/observability-operator/api/v1alpha1"
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	"github.com/udmire/observability-operator/pkg/utils"
)

func mergeServiceAccount(manifest *core_v1.ServiceAccount, serviceAccount *v1alpha1.ServiceAccountSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if serviceAccount == nil {
		return
	}

	if serviceAccount.AutomountServiceAccountToken != nil {
		manifest.AutomountServiceAccountToken = serviceAccount.AutomountServiceAccountToken
	}

	if len(serviceAccount.Secrets) > 0 {
		manifest.Secrets = serviceAccount.Secrets
	}

	if len(serviceAccount.ImagePullSecrets) > 0 {
		manifest.ImagePullSecrets = serviceAccount.ImagePullSecrets
	}
}

func mergeClusterRole(manifest *rbac_v1.ClusterRole, clusterRole *v1alpha1.ClusterRoleSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if clusterRole == nil {
		return
	}

	if clusterRole.AggregationRule != nil {
		manifest.AggregationRule = clusterRole.AggregationRule
	}

	if len(clusterRole.Rules) > 0 {
		manifest.Rules = clusterRole.Rules
	}
}

func mergeClusterRoleBinding(manifest *rbac_v1.ClusterRoleBinding, clusterRoleBinding *v1alpha1.ClusterRoleBindingSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	manifest.Subjects = updateSubjectsNamespace(manifest.Subjects, ns)

	if clusterRoleBinding == nil {
		return
	}

	if clusterRoleBinding.RoleRef != nil {
		manifest.RoleRef = *clusterRoleBinding.RoleRef
	}

	if len(clusterRoleBinding.Subjects) > 0 {
		manifest.Subjects = clusterRoleBinding.Subjects
	}
}

func mergeRole(manifest *rbac_v1.Role, role *v1alpha1.RoleSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if role == nil {
		return
	}

	if len(role.Rules) > 0 {
		manifest.Rules = role.Rules
	}
}

func mergeRoleBinding(manifest *rbac_v1.RoleBinding, roleBinding *v1alpha1.RoleBindingSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	manifest.Subjects = updateSubjectsNamespace(manifest.Subjects, ns)

	if roleBinding == nil {
		return
	}

	if roleBinding.RoleRef != nil {
		manifest.RoleRef = *roleBinding.RoleRef
	}

	if len(roleBinding.Subjects) > 0 {
		manifest.Subjects = roleBinding.Subjects
	}
}

func mergeIngress(manifest *networking_v1.Ingress, ingress *v1alpha1.IngressSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	if ingress == nil {
		return
	}

	if ingress.IngressClassName != nil {
		manifest.Spec.IngressClassName = ingress.IngressClassName
	}

	if ingress.DefaultBackend != nil {
		manifest.Spec.DefaultBackend = ingress.DefaultBackend
	}

	if len(ingress.TLS) > 0 {
		manifest.Spec.TLS = ingress.TLS
	}

	if len(ingress.Rules) > 0 {
		manifest.Spec.Rules = ingress.Rules
	}
}

func configMapsCustom(manifest *manifest.Manifests, configmaps map[string]*v1alpha1.ConfigMapSpec, ns string, labels map[string]string) {
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

func secretsCustom(manifest *manifest.Manifests, secrets map[string]*v1alpha1.SecretSpec, ns string, labels map[string]string) {
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

func servicesCustom(manifest *manifest.Manifests, services map[string]*v1alpha1.ServiceSpec, ns string, labels map[string]string) {
	for _, svc := range manifest.Services {
		mergeService(svc, nil, ns, labels)
	}

	if len(services) < 1 {
		return
	}

	merged := make(map[string]string)
	for _, srv := range manifest.Services {
		if service, ok := services[srv.Name]; ok {
			mergeService(srv, service, ns, labels)
			merged[srv.Name] = ""
		}
		continue
	}

	for name, srv := range services {
		if _, ok := merged[name]; ok {
			continue
		}
		manifest.Services = append(manifest.Services, &core_v1.Service{
			ObjectMeta: metav1.ObjectMeta{
				Name:      name,
				Namespace: ns,
				Labels:    labels,
			},
			Spec: core_v1.ServiceSpec{
				ClusterIP: srv.ClusterIP,
				Ports:     srv.Ports,
				Selector:  srv.Selector,
			},
		})
	}
}

func mergeService(manifest *core_v1.Service, service *v1alpha1.ServiceSpec, ns string, labels map[string]string) {
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	mergeSelectorLabels(manifest.Spec.Selector, labels)

	if service == nil {
		return
	}

	manifest.Spec.ClusterIP = service.ClusterIP

	if len(service.Ports) > 0 {
		manifest.Spec.Ports = service.Ports
	}

	if len(service.Selector) > 0 {
		manifest.Spec.Selector = service.Selector
	}
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

func mergePodTemplateObjectMeta(meta *metav1.ObjectMeta, labels map[string]string) {
	mergeSelectorLabels(meta.Labels, labels)
}

func mergeSelectorLabels(meta map[string]string, labels map[string]string) {
	if len(labels) == 0 {
		return
	}

	for key, value := range labels {
		if _, ok := meta[key]; ok {
			continue
		}
		if utils.ShouldIgnoredInSelector(key) {
			continue
		}
		meta[key] = value
	}
}

func updateSubjectsNamespace(subjects []rbac_v1.Subject, ns string) (result []rbac_v1.Subject) {
	if len(ns) > 0 {
		for _, sub := range subjects {
			result = append(result, rbac_v1.Subject{
				Kind:      sub.Kind,
				APIGroup:  sub.APIGroup,
				Name:      sub.Name,
				Namespace: ns,
			})
		}
		return result
	}
	return subjects
}
