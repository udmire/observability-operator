package v1alpha1

import (
	apps_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type AppSpec struct {
	Name      string   `json:"name,omitempty"`
	Namespace string   `json:"namespace,omitempty"`
	Template  Template `json:"template"`
	Singleton bool     `json:"singleton,omitempty"`

	Registry string `json:"registry,omitempty"`

	CommonSpec `json:",inline"`
	Components map[string]ComponentSpec `json:"components,omitempty"`

	// App Dependencies, must be ready before AppSpec Applied.
	Dependencies AppDepsSpec `json:"deps,omitempty"`
}

type Template struct {
	Name    string `json:"name"`
	Version string `json:"version,omitempty"`
}

type ComponentSpec struct {
	CommonSpec   `json:",inline"`
	WorkloadSpec `json:",inline"`
	HPA          *HpaSpec `json:"hpa,omitempty"`
}

type CommonSpec struct {
	ConfigMaps map[string]*ConfigMapSpec `json:"configmaps,omitempty"`
	Secrets    map[string]*SecretSpec    `json:"secrets,omitempty"`
	Services   map[string]*ServiceSpec   `json:"services,omitempty"`

	ServiceAccount     *ServiceAccountSpec     `json:"serviceAccount,omitempty"`
	ClusterRole        *ClusterRoleSpec        `json:"clusterRole,omitempty"`
	ClusterRoleBinding *ClusterRoleBindingSpec `json:"clusterRoleBinding,omitempty"`
	Role               *RoleSpec               `json:"role,omitempty"`
	RoleBinding        *RoleBindingSpec        `json:"roleBinding,omitempty"`
	Ingress            *IngressSpec            `json:"ingress,omitempty"`
}

type WorkloadSpec struct {
	Deployment  *DeploymentSpec  `json:"deployment,omitempty"`
	DaemonSet   *DaemonSetSpec   `json:"daemonset,omitempty"`
	StatefulSet *StatefulSetSpec `json:"statefulset,omitempty"`
	ReplicaSet  *ReplicaSetSpec  `json:"replicaset,omitempty"`
	Job         *JobSpec         `json:"job,omitempty"`
	CronJob     *CronJobSpec     `json:"cronjob,omitempty"`
}

type ServiceSpec struct {
	Ports     []core_v1.ServicePort `json:"ports,omitempty"`
	Selector  map[string]string     `json:"selector,omitempty"`
	ClusterIP string                `json:"clusterIP,omitempty"`
}
type ConfigMapSpec struct {
	Data map[string]string `json:"data,omitempty"`
}
type SecretSpec struct {
	StringData map[string]string `json:"stringData,omitempty"`
}
type ServiceAccountSpec struct {
	Secrets                      []core_v1.ObjectReference      `json:"secrets,omitempty"`
	ImagePullSecrets             []core_v1.LocalObjectReference `json:"imagePullSecrets,omitempty"`
	AutomountServiceAccountToken *bool                          `json:"automountServiceAccountToken,omitempty"`
}
type ClusterRoleSpec struct {
	Rules           []rbac_v1.PolicyRule     `json:"rules,omitempty"`
	AggregationRule *rbac_v1.AggregationRule `json:"aggregationRule,omitempty"`
}
type ClusterRoleBindingSpec struct {
	Subjects []rbac_v1.Subject `json:"subjects,omitempty"`
	RoleRef  *rbac_v1.RoleRef  `json:"roleRef,omitempty"`
}
type RoleSpec struct {
	Rules []rbac_v1.PolicyRule `json:"rules,omitempty"`
}
type RoleBindingSpec struct {
	Subjects []rbac_v1.Subject `json:"subjects,omitempty"`
	RoleRef  *rbac_v1.RoleRef  `json:"roleRef,omitempty"`
}
type IngressSpec struct {
	IngressClassName *string                       `json:"ingressClassName,omitempty"`
	DefaultBackend   *networking_v1.IngressBackend `json:"defaultBackend,omitempty"`
	TLS              []networking_v1.IngressTLS    `json:"tls,omitempty"`
	Rules            []networking_v1.IngressRule   `json:"rules,omitempty"`
}
type DeploymentSpec struct {
	Replicas             *int32                      `json:"replicas,omitempty"`
	Selector             *metav1.LabelSelector       `json:"selector,omitempty"`
	Template             *PodTemplateSpec            `json:"template,omitempty"`
	Strategy             *apps_v1.DeploymentStrategy `json:"strategy,omitempty" patchStrategy:"retainKeys"`
	MinReadySeconds      *int32                      `json:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int32                      `json:"revisionHistoryLimit,omitempty"`
}
type DaemonSetSpec struct {
	Selector             *metav1.LabelSelector            `json:"selector,omitempty"`
	Template             *PodTemplateSpec                 `json:"template,omitempty"`
	UpdateStrategy       *apps_v1.DaemonSetUpdateStrategy `json:"updateStrategy,omitempty"`
	MinReadySeconds      *int32                           `json:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int32                           `json:"revisionHistoryLimit,omitempty"`
}
type StatefulSetSpec struct {
	Replicas             *int32                             `json:"replicas,omitempty"`
	Selector             *metav1.LabelSelector              `json:"selector,omitempty"`
	Template             *PodTemplateSpec                   `json:"template,omitempty"`
	VolumeClaimTemplates []core_v1.PersistentVolumeClaim    `json:"volumeClaimTemplates,omitempty"`
	ServiceName          *string                            `json:"serviceName,omitempty"`
	PodManagementPolicy  *apps_v1.PodManagementPolicyType   `json:"podManagementPolicy,omitempty"`
	UpdateStrategy       *apps_v1.StatefulSetUpdateStrategy `json:"updateStrategy,omitempty"`
	RevisionHistoryLimit *int32                             `json:"revisionHistoryLimit,omitempty"`
	MinReadySeconds      *int32                             `json:"minReadySeconds,omitempty"`

	PersistentVolumeClaimRetentionPolicy *apps_v1.StatefulSetPersistentVolumeClaimRetentionPolicy `json:"persistentVolumeClaimRetentionPolicy,omitempty"`

	Ordinals *apps_v1.StatefulSetOrdinals `json:"ordinals,omitempty"`
}
type ReplicaSetSpec struct {
	Replicas        *int32                `json:"replicas,omitempty"`
	MinReadySeconds *int32                `json:"minReadySeconds,omitempty"`
	Selector        *metav1.LabelSelector `json:"selector,omitempty"`
	Template        *PodTemplateSpec      `json:"template,omitempty"`
}
type JobSpec struct {
	Parallelism             *int32                     `json:"parallelism,omitempty"`
	Completions             *int32                     `json:"completions,omitempty"`
	ActiveDeadlineSeconds   *int64                     `json:"activeDeadlineSeconds,omitempty"`
	PodFailurePolicy        *batch_v1.PodFailurePolicy `json:"podFailurePolicy,omitempty"`
	BackoffLimit            *int32                     `json:"backoffLimit,omitempty"`
	Selector                *metav1.LabelSelector      `json:"selector,omitempty"`
	ManualSelector          *bool                      `json:"manualSelector,omitempty"`
	Template                *PodTemplateSpec           `json:"template,omitempty"`
	TTLSecondsAfterFinished *int32                     `json:"ttlSecondsAfterFinished,omitempty"`
	CompletionMode          *batch_v1.CompletionMode   `json:"completionMode,omitempty"`
	Suspend                 *bool                      `json:"suspend,omitempty"`
}
type CronJobSpec struct {
	Schedule                   *string                     `json:"schedule,omitempty"`
	TimeZone                   *string                     `json:"timeZone,omitempty"`
	StartingDeadlineSeconds    *int64                      `json:"startingDeadlineSeconds,omitempty"`
	ConcurrencyPolicy          *batch_v1.ConcurrencyPolicy `json:"concurrencyPolicy,omitempty"`
	Suspend                    *bool                       `json:"suspend,omitempty"`
	JobTemplate                *JobTemplateSpec            `json:"jobTemplate,omitempty"`
	SuccessfulJobsHistoryLimit *int32                      `json:"successfulJobsHistoryLimit,omitempty"`
	FailedJobsHistoryLimit     *int32                      `json:"failedJobsHistoryLimit,omitempty"`
}
type HpaSpec struct {
	ScaleTargetRef                 autoscaling_v1.CrossVersionObjectReference `json:"scaleTargetRef"`
	MinReplicas                    *int32                                     `json:"minReplicas,omitempty"`
	MaxReplicas                    int32                                      `json:"maxReplicas"`
	TargetCPUUtilizationPercentage *int32                                     `json:"targetCPUUtilizationPercentage,omitempty"`
}
type JobTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              JobSpec `json:"spec,omitempty"`
}
type PodTemplateSpec struct {
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              PodSpec `json:"spec,omitempty"`
}
type PodSpec struct {
	Containers                []core_v1.Container                `json:"containers,omitempty"`
	InitContainers            []core_v1.Container                `json:"initContainers,omitempty"`
	Volumes                   []core_v1.Volume                   `json:"volumes,omitempty"`
	DNSPolicy                 *core_v1.DNSPolicy                 `json:"dnsPolicy,omitempty"`
	NodeSelector              *map[string]string                 `json:"nodeSelector,omitempty"`
	ServiceAccountName        *string                            `json:"serviceAccountName,omitempty"`
	SecurityContext           *core_v1.PodSecurityContext        `json:"securityContext,omitempty"`
	ImagePullSecrets          []core_v1.LocalObjectReference     `json:"imagePullSecrets,omitempty"`
	Affinity                  *core_v1.Affinity                  `json:"affinity,omitempty"`
	Tolerations               []core_v1.Toleration               `json:"tolerations,omitempty"`
	HostAliases               []core_v1.HostAlias                `json:"hostAliases,omitempty"`
	PriorityClassName         *string                            `json:"priorityClassName,omitempty"`
	TopologySpreadConstraints []core_v1.TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty"`
}
