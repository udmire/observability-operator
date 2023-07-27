package templates

import (
	"io/ioutil"
	"os"
	"regexp"

	"gopkg.in/yaml.v2"
	app_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
)

const (
	fileConfigMap          = "(.+)[-_](configmap|cm|config).ya?ml"
	fileSecret             = "(.+)[-_]secret.ya?ml"
	fileServiceAccount     = "(.+)[-_](sa|serviceaccount).ya?ml"
	fileClusterRole        = "(.+)[-_](cr|clusterrole).ya?ml"
	fileClusterRoleBinding = "(.+)[-_](crb|clusterrolebinding).ya?ml"
	fileRole               = "(.+)[-_]role.ya?ml"
	fileRoleBinding        = "(.+)[-_](rb|rolebinding).ya?ml"
	fileIngress            = "(.+)[-_]ingress.ya?ml"
	fileService            = "(.+)[-_](svc|service).ya?ml"
	fileDeployment         = "(.+)[-_](dep|deploy|deployment).ya?ml"
	fileDaemonSet          = "(.+)[-_](ds|daemonset).ya?ml"
	fileStatefulSet        = "(.+)[-_](sts|statefulset).ya?ml"
	fileReplicaSet         = "(.+)[-_](rs|replicaset).ya?ml"
	fileJob                = "(.+)[-_]job.ya?ml"
	fileCronJob            = "(.+)[-_]cronjob.ya?ml"
	fileHPA                = "(.+)[-_](hpa|horizontalpodautoscaler).ya?ml"
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

	Services []*core_v1.Service

	Deployment  *app_v1.Deployment
	DaemonSet   *app_v1.DaemonSet
	StatefulSet *app_v1.StatefulSet
	ReplicaSet  *app_v1.ReplicaSet
	Job         *batch_v1.Job
	CronJob     *batch_v1.CronJob

	HPA *autoscaling_v1.HorizontalPodAutoscaler
}

type TemplateBase struct {
	Name          string
	Version       string
	TemplateFiles []*os.File
}

type AppTemplate struct {
	TemplateBase
	Workloads []*WorkloadTemplate
}

type WorkloadTemplate struct {
	TemplateBase
}

type TemplateBuilder interface {
	Build(template *AppTemplate) *AppManifests
}

func NewTemplateBuilder(template *AppTemplate) TemplateBuilder {
	if template == nil {
		return nil
	}
	return &templateBuilder{
		template: template,
	}
}

type templateBuilder struct {
	template *AppTemplate
}

func (b *templateBuilder) Build(template *AppTemplate) *AppManifests {
	var manifests *AppManifests
	for _, tempFile := range template.TemplateFiles {
		resType, _, content := b.recognize(tempFile)
		var err error

		switch resType {
		case ConfigMap:
			cm := &core_v1.ConfigMap{}
			err = yaml.Unmarshal(content, cm)
			if err != nil {
				panic(err)
			}
			manifests.ConfigMaps = append(manifests.ConfigMaps, cm)
		case Secret:
			sec := &core_v1.Secret{}
			err = yaml.Unmarshal(content, sec)
			if err != nil {
				panic(err)
			}
			manifests.Secrets = append(manifests.Secrets, sec)
		case ServiceAccount:
			sa := &core_v1.ServiceAccount{}
			err = yaml.Unmarshal(content, sa)
			if err != nil {
				panic(err)
			}
			manifests.ServiceAccount = sa
		case ClusterRole:
			role := &rbac_v1.ClusterRole{}
			err = yaml.Unmarshal(content, role)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRole = role
		case ClusterRoleBinding:
			rb := &rbac_v1.ClusterRoleBinding{}
			err = yaml.Unmarshal(content, rb)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRoleBinding = rb
		case Role:
			role := &rbac_v1.Role{}
			err = yaml.Unmarshal(content, role)
			if err != nil {
				panic(err)
			}
			manifests.Role = role
		case RoleBinding:
			rb := &rbac_v1.RoleBinding{}
			err = yaml.Unmarshal(content, rb)
			if err != nil {
				panic(err)
			}
			manifests.RoleBinding = rb
		case Ingress:
			ing := &networking_v1.Ingress{}
			err = yaml.Unmarshal(content, ing)
			if err != nil {
				panic(err)
			}
			manifests.Ingress = ing

		default:
		}
	}

	for _, comp := range template.Workloads {
		manifests.CompsMenifests = append(manifests.CompsMenifests, b.BuildComp(comp))
	}

	return manifests
}

func (b *templateBuilder) BuildComp(template *WorkloadTemplate) *CompManifests {
	var manifests *CompManifests
	for _, tempFile := range template.TemplateFiles {
		resType, _, content := b.recognize(tempFile)
		var err error

		switch resType {
		case ConfigMap:
			cm := &core_v1.ConfigMap{}
			err = yaml.Unmarshal(content, cm)
			if err != nil {
				panic(err)
			}
			manifests.ConfigMaps = append(manifests.ConfigMaps, cm)
		case Secret:
			sec := &core_v1.Secret{}
			err = yaml.Unmarshal(content, sec)
			if err != nil {
				panic(err)
			}
			manifests.Secrets = append(manifests.Secrets, sec)
		case ServiceAccount:
			sa := &core_v1.ServiceAccount{}
			err = yaml.Unmarshal(content, sa)
			if err != nil {
				panic(err)
			}
			manifests.ServiceAccount = sa
		case ClusterRole:
			role := &rbac_v1.ClusterRole{}
			err = yaml.Unmarshal(content, role)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRole = role
		case ClusterRoleBinding:
			rb := &rbac_v1.ClusterRoleBinding{}
			err = yaml.Unmarshal(content, rb)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRoleBinding = rb
		case Role:
			role := &rbac_v1.Role{}
			err = yaml.Unmarshal(content, role)
			if err != nil {
				panic(err)
			}
			manifests.Role = role
		case RoleBinding:
			rb := &rbac_v1.RoleBinding{}
			err = yaml.Unmarshal(content, rb)
			if err != nil {
				panic(err)
			}
			manifests.RoleBinding = rb
		case Ingress:
			ing := &networking_v1.Ingress{}
			err = yaml.Unmarshal(content, ing)
			if err != nil {
				panic(err)
			}
			manifests.Ingress = ing
		case Service:
			svc := &core_v1.Service{}
			err = yaml.Unmarshal(content, svc)
			if err != nil {
				panic(err)
			}
			manifests.Services = append(manifests.Services, svc)
		case Deployment:
			wl := &app_v1.Deployment{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.Deployment = wl
		case DaemonSet:
			wl := &app_v1.DaemonSet{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.DaemonSet = wl
		case StatefulSet:
			wl := &app_v1.StatefulSet{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.StatefulSet = wl
		case ReplicaSet:
			wl := &app_v1.ReplicaSet{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.ReplicaSet = wl
		case Job:
			wl := &batch_v1.Job{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.Job = wl
		case CronJob:
			wl := &batch_v1.CronJob{}
			err = yaml.Unmarshal(content, wl)
			if err != nil {
				panic(err)
			}
			manifests.CronJob = wl
		case HPA:
			hpa := &autoscaling_v1.HorizontalPodAutoscaler{}
			err = yaml.Unmarshal(content, hpa)
			if err != nil {
				panic(err)
			}
			manifests.HPA = hpa
		default:
		}
	}

	return manifests
}

func (b *templateBuilder) recognize(file *os.File) (ManifestType, string, []byte) {
	content, err := ioutil.ReadFile(file.Name())
	if err != nil {
		return -1, "", nil
	}
	for i := 0; i < len(filePatterns); i++ {
		match := regexp.MustCompile(filePatterns[i]).FindStringSubmatch(file.Name())
		if len(match) > 1 {
			firstGroupContent := match[1]
			return ManifestTypes[i], firstGroupContent, content
		}
	}
	return -1, "", nil
}
