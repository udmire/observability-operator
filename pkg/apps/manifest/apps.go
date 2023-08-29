package manifest

import (
	app_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
)

type ManifestType int

const (
	ConfigMap ManifestType = iota
	Secret
	ServiceAccount
	ClusterRole
	ClusterRoleBinding
	Role
	RoleBinding
	Ingress
	Service
	Deployment
	DaemonSet
	StatefulSet
	ReplicaSet
	Job
	CronJob
	HPA
)

var ManifestTypes = []ManifestType{
	ConfigMap,
	Secret,
	ServiceAccount,
	ClusterRole,
	ClusterRoleBinding,
	Role,
	RoleBinding,
	Ingress,
	Service,
	Deployment,
	DaemonSet,
	StatefulSet,
	ReplicaSet,
	Job,
	CronJob,
	HPA,
}

type Manifests struct {
	ConfigMaps []*core_v1.ConfigMap
	Secrets    []*core_v1.Secret
	Services   []*core_v1.Service

	ServiceAccount     *core_v1.ServiceAccount
	ClusterRole        *rbac_v1.ClusterRole
	ClusterRoleBinding *rbac_v1.ClusterRoleBinding
	Role               *rbac_v1.Role
	RoleBinding        *rbac_v1.RoleBinding
	Ingress            *networking_v1.Ingress
}

type AppManifests struct {
	Manifests

	CompsMenifests []*CompManifests
}

type CompManifests struct {
	Manifests

	Name string

	Deployment  *app_v1.Deployment
	DaemonSet   *app_v1.DaemonSet
	StatefulSet *app_v1.StatefulSet
	ReplicaSet  *app_v1.ReplicaSet
	Job         *batch_v1.Job
	CronJob     *batch_v1.CronJob

	HPA *autoscaling_v1.HorizontalPodAutoscaler
}
