package model

import (
	"github.com/udmire/observability-operator/pkg/apps/manifest"
	apps_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (a *App) Build() *manifest.AppManifests {
	mani := &manifest.AppManifests{Manifests: manifest.Manifests{}}
	convertCommon(a.Common, &mani.Manifests)
	for _, comp := range a.Components {
		component := &manifest.CompManifests{Manifests: manifest.Manifests{}, Name: comp.Name}
		convertComponent(comp, component)
		mani.CompsMenifests = append(mani.CompsMenifests, component)
	}
	return mani
}

func convertCommon(c Common, mani *manifest.Manifests) {
	if c.ClusterRole != nil {
		mani.ClusterRole = &v1.ClusterRole{
			ObjectMeta: buildObjectMetaWithoutNamespace(c.ClusterRole.GenericModel),
		}
	}
	if c.ClusterRoleBinding != nil {
		mani.ClusterRoleBinding = &v1.ClusterRoleBinding{
			ObjectMeta: buildObjectMetaWithoutNamespace(c.ClusterRoleBinding.GenericModel),
		}
	}

	if c.Role != nil {
		mani.Role = &v1.Role{
			ObjectMeta: buildObjectMeta(c.Role.GenericModel),
		}
	}
	if c.RoleBinding != nil {
		mani.RoleBinding = &v1.RoleBinding{
			ObjectMeta: buildObjectMeta(c.RoleBinding.GenericModel),
		}
	}
	if c.ServiceAccount != nil {
		mani.ServiceAccount = &core_v1.ServiceAccount{
			ObjectMeta: buildObjectMeta(c.ServiceAccount.GenericModel),
		}
	}
	if c.Ingress != nil {
		mani.Ingress = &networking_v1.Ingress{
			ObjectMeta: buildObjectMeta(c.Ingress.GenericModel),
		}
	}
	for _, sec := range c.Secrets {
		mani.Secrets = append(mani.Secrets, &core_v1.Secret{
			ObjectMeta: buildObjectMeta(sec.GenericModel),
		})
	}
	for _, svc := range c.Services {
		mani.Services = append(mani.Services, &core_v1.Service{
			ObjectMeta: buildObjectMeta(svc.GenericModel),
		})
	}
	for _, cm := range c.ConfigMaps {
		mani.ConfigMaps = append(mani.ConfigMaps, &core_v1.ConfigMap{
			ObjectMeta: buildObjectMeta(cm.GenericModel),
		})
	}
}

func convertComponent(comp *Component, mani *manifest.CompManifests) {

	convertCommon(comp.Common, &mani.Manifests)

	if comp.HPA != nil {
		mani.HPA = &autoscaling_v1.HorizontalPodAutoscaler{
			ObjectMeta: buildObjectMeta(comp.HPA.GenericModel),
		}
	}

	if comp.Deployment != nil {
		mani.Deployment = &apps_v1.Deployment{
			ObjectMeta: buildObjectMeta(comp.Deployment.GenericModel),
		}
	}
	if comp.DaemonSet != nil {
		mani.DaemonSet = &apps_v1.DaemonSet{
			ObjectMeta: buildObjectMeta(comp.DaemonSet.GenericModel),
		}
	}
	if comp.StatefulSet != nil {
		mani.StatefulSet = &apps_v1.StatefulSet{
			ObjectMeta: buildObjectMeta(comp.StatefulSet.GenericModel),
		}
	}
	if comp.ReplicaSet != nil {
		mani.ReplicaSet = &apps_v1.ReplicaSet{
			ObjectMeta: buildObjectMeta(comp.ReplicaSet.GenericModel),
		}
	}
	if comp.Job != nil {
		mani.Job = &batch_v1.Job{
			ObjectMeta: buildObjectMeta(comp.Job.GenericModel),
		}
	}
	if comp.CronJob != nil {
		mani.CronJob = &batch_v1.CronJob{
			ObjectMeta: buildObjectMeta(comp.CronJob.GenericModel),
		}
	}
}

func buildObjectMeta(model GenericModel) metav1.ObjectMeta {
	meta := buildObjectMetaWithoutNamespace(model)

	if model.DefaultNamespace != "" {
		meta.Namespace = model.DefaultNamespace
	}

	return meta
}

func buildObjectMetaWithoutNamespace(model GenericModel) metav1.ObjectMeta {
	meta := metav1.ObjectMeta{
		Name:   model.Name,
		Labels: make(map[string]string),
	}

	mergeLabels(meta.Labels, model.Labels)
	return meta
}

// Merge not overwrite
func mergeLabels(ori, merge map[string]string) {
	for key, value := range merge {
		if _, exists := ori[key]; !exists {
			ori[key] = value
		}
	}
}
