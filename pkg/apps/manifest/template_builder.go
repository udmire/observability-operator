package manifest

import (
	"regexp"

	app_v1 "k8s.io/api/apps/v1"
	autoscaling_v1 "k8s.io/api/autoscaling/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/util/yaml"

	"github.com/udmire/observability-operator/pkg/templates/template"
)

type Builder interface {
	Build() *AppManifests
}

func NewTemplateBuilder(template *template.AppTemplate) Builder {
	if template == nil {
		return nil
	}
	return &templateBuilder{
		template: template,
	}
}

type templateBuilder struct {
	template *template.AppTemplate
}

func (b *templateBuilder) Build() *AppManifests {
	manifests := &AppManifests{}
	for _, tempFile := range b.template.TemplateFiles {
		resType, _ := recognize(tempFile)
		var err error

		switch resType {
		case ConfigMap:
			cm := core_v1.ConfigMap{}
			err = yaml.Unmarshal(tempFile.Content, &cm)
			if err != nil {
				panic(err)
			}
			manifests.ConfigMaps = append(manifests.ConfigMaps, &cm)
		case Secret:
			sec := &core_v1.Secret{}
			err = yaml.Unmarshal(tempFile.Content, sec)
			if err != nil {
				panic(err)
			}
			manifests.Secrets = append(manifests.Secrets, sec)
		case ServiceAccount:
			sa := &core_v1.ServiceAccount{}
			err = yaml.Unmarshal(tempFile.Content, sa)
			if err != nil {
				panic(err)
			}
			manifests.ServiceAccount = sa
		case ClusterRole:
			role := &rbac_v1.ClusterRole{}
			err = yaml.Unmarshal(tempFile.Content, role)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRole = role
		case ClusterRoleBinding:
			rb := &rbac_v1.ClusterRoleBinding{}
			err = yaml.Unmarshal(tempFile.Content, rb)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRoleBinding = rb
		case Role:
			role := &rbac_v1.Role{}
			err = yaml.Unmarshal(tempFile.Content, role)
			if err != nil {
				panic(err)
			}
			manifests.Role = role
		case RoleBinding:
			rb := &rbac_v1.RoleBinding{}
			err = yaml.Unmarshal(tempFile.Content, rb)
			if err != nil {
				panic(err)
			}
			manifests.RoleBinding = rb
		case Ingress:
			ing := &networking_v1.Ingress{}
			err = yaml.Unmarshal(tempFile.Content, ing)
			if err != nil {
				panic(err)
			}
			manifests.Ingress = ing

		default:
		}
	}

	for _, comp := range b.template.Workloads {
		manifests.CompsMenifests = append(manifests.CompsMenifests, b.BuildComp(comp))
	}

	return manifests
}

func (b *templateBuilder) BuildComp(template *template.WorkloadTemplate) *CompManifests {
	manifests := &CompManifests{Name: template.Name}
	for _, tempFile := range template.TemplateFiles {
		resType, _ := recognize(tempFile)
		var err error

		switch resType {
		case ConfigMap:
			cm := &core_v1.ConfigMap{}
			err = yaml.Unmarshal(tempFile.Content, cm)
			if err != nil {
				panic(err)
			}
			manifests.ConfigMaps = append(manifests.ConfigMaps, cm)
		case Secret:
			sec := &core_v1.Secret{}
			err = yaml.Unmarshal(tempFile.Content, sec)
			if err != nil {
				panic(err)
			}
			manifests.Secrets = append(manifests.Secrets, sec)
		case ServiceAccount:
			sa := &core_v1.ServiceAccount{}
			err = yaml.Unmarshal(tempFile.Content, sa)
			if err != nil {
				panic(err)
			}
			manifests.ServiceAccount = sa
		case ClusterRole:
			role := &rbac_v1.ClusterRole{}
			err = yaml.Unmarshal(tempFile.Content, role)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRole = role
		case ClusterRoleBinding:
			rb := &rbac_v1.ClusterRoleBinding{}
			err = yaml.Unmarshal(tempFile.Content, rb)
			if err != nil {
				panic(err)
			}
			manifests.ClusterRoleBinding = rb
		case Role:
			role := &rbac_v1.Role{}
			err = yaml.Unmarshal(tempFile.Content, role)
			if err != nil {
				panic(err)
			}
			manifests.Role = role
		case RoleBinding:
			rb := &rbac_v1.RoleBinding{}
			err = yaml.Unmarshal(tempFile.Content, rb)
			if err != nil {
				panic(err)
			}
			manifests.RoleBinding = rb
		case Ingress:
			ing := &networking_v1.Ingress{}
			err = yaml.Unmarshal(tempFile.Content, ing)
			if err != nil {
				panic(err)
			}
			manifests.Ingress = ing
		case Service:
			svc := &core_v1.Service{}
			err = yaml.Unmarshal(tempFile.Content, svc)
			if err != nil {
				panic(err)
			}
			manifests.Services = append(manifests.Services, svc)
		case Deployment:
			wl := &app_v1.Deployment{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.Deployment = wl
		case DaemonSet:
			wl := &app_v1.DaemonSet{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.DaemonSet = wl
		case StatefulSet:
			wl := &app_v1.StatefulSet{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.StatefulSet = wl
		case ReplicaSet:
			wl := &app_v1.ReplicaSet{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.ReplicaSet = wl
		case Job:
			wl := &batch_v1.Job{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.Job = wl
		case CronJob:
			wl := &batch_v1.CronJob{}
			err = yaml.Unmarshal(tempFile.Content, wl)
			if err != nil {
				panic(err)
			}
			manifests.CronJob = wl
		case HPA:
			hpa := &autoscaling_v1.HorizontalPodAutoscaler{}
			err = yaml.Unmarshal(tempFile.Content, hpa)
			if err != nil {
				panic(err)
			}
			manifests.HPA = hpa
		default:
		}
	}

	return manifests
}

func recognize(file *template.TemplateFile) (ManifestType, string) {
	for i := 0; i < len(filePatterns); i++ {
		match := regexp.MustCompile(filePatterns[i]).FindStringSubmatch(file.FileName)
		if len(match) > 1 {
			firstGroupContent := match[1]
			return ManifestTypes[i], firstGroupContent
		}
	}
	return -1, ""
}
