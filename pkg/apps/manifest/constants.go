package manifest

const (
	fileConfigMap          = "([^/]+)[-_](configmap|cm|config).ya?ml"
	fileSecret             = "([^/]+)[-_]secret.ya?ml"
	fileServiceAccount     = "([^/]+)[-_](sa|serviceaccount).ya?ml"
	fileClusterRole        = "([^/]+)[-_](cr|clusterrole).ya?ml"
	fileClusterRoleBinding = "([^/]+)[-_](crb|clusterrolebinding).ya?ml"
	fileRole               = "([^/]+)[-_]role.ya?ml"
	fileRoleBinding        = "([^/]+)[-_](rb|rolebinding).ya?ml"
	fileIngress            = "([^/]+)[-_]ingress.ya?ml"
	fileService            = "([^/]+)[-_](svc|service).ya?ml"
	fileDeployment         = "([^/]+)[-_](dep|deploy|deployment).ya?ml"
	fileDaemonSet          = "([^/]+)[-_](ds|daemonset).ya?ml"
	fileStatefulSet        = "([^/]+)[-_](sts|statefulset).ya?ml"
	fileReplicaSet         = "([^/]+)[-_](rs|replicaset).ya?ml"
	fileJob                = "([^/]+)[-_]job.ya?ml"
	fileCronJob            = "([^/]+)[-_]cronjob.ya?ml"
	fileHPA                = "([^/]+)[-_](hpa|horizontalpodautoscaler).ya?ml"
)

var filePatterns = []string{
	fileConfigMap,
	fileSecret,
	fileServiceAccount,
	fileClusterRole,
	fileClusterRoleBinding,
	fileRole,
	fileRoleBinding,
	fileIngress,
	fileService,
	fileDeployment,
	fileDaemonSet,
	fileStatefulSet,
	fileReplicaSet,
	fileJob,
	fileCronJob,
	fileHPA,
}
