package specs

import (
	"fmt"

	app_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	core_v1 "k8s.io/api/core/v1"

	"github.com/udmire/observability-operator/api/v1alpha1"
)

func mergeDeployment(manifest *app_v1.Deployment, spec *v1alpha1.DeploymentSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}

	if spec.Replicas != nil {
		manifest.Spec.Replicas = spec.Replicas
	}
	if spec.Selector != nil {
		manifest.Spec.Selector = spec.Selector
	}
	if spec.Strategy != nil {
		manifest.Spec.Strategy = *spec.Strategy
	}
	if spec.MinReadySeconds != nil {
		manifest.Spec.MinReadySeconds = *spec.MinReadySeconds
	}
	if spec.RevisionHistoryLimit != nil {
		manifest.Spec.RevisionHistoryLimit = spec.RevisionHistoryLimit
	}
	return nil
}

func mergeDaemonset(manifest *app_v1.DaemonSet, spec *v1alpha1.DaemonSetSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}

	if spec.Selector != nil {
		manifest.Spec.Selector = spec.Selector
	}
	if spec.UpdateStrategy != nil {
		manifest.Spec.UpdateStrategy = *spec.UpdateStrategy
	}
	if spec.MinReadySeconds != nil {
		manifest.Spec.MinReadySeconds = *spec.MinReadySeconds
	}
	if spec.RevisionHistoryLimit != nil {
		manifest.Spec.RevisionHistoryLimit = spec.RevisionHistoryLimit
	}
	return nil
}

func mergeStatefulSet(manifest *app_v1.StatefulSet, spec *v1alpha1.StatefulSetSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)
	updateServiceName(manifest, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}
	if spec.Replicas != nil {
		manifest.Spec.Replicas = spec.Replicas
	}
	if spec.Selector != nil {
		manifest.Spec.Selector = spec.Selector
	}

	if spec.VolumeClaimTemplates != nil {
		manifest.Spec.VolumeClaimTemplates = spec.VolumeClaimTemplates
	}
	if spec.ServiceName != nil {
		manifest.Spec.ServiceName = *spec.ServiceName
		updateServiceName(manifest, prefix)
	}
	if spec.PodManagementPolicy != nil {
		manifest.Spec.PodManagementPolicy = *spec.PodManagementPolicy
	}
	if spec.UpdateStrategy != nil {
		manifest.Spec.UpdateStrategy = *spec.UpdateStrategy
	}
	if spec.RevisionHistoryLimit != nil {
		manifest.Spec.RevisionHistoryLimit = spec.RevisionHistoryLimit
	}
	if spec.MinReadySeconds != nil {
		manifest.Spec.MinReadySeconds = *spec.MinReadySeconds
	}
	if spec.PersistentVolumeClaimRetentionPolicy != nil {
		manifest.Spec.PersistentVolumeClaimRetentionPolicy = spec.PersistentVolumeClaimRetentionPolicy
	}
	if spec.Ordinals != nil {
		manifest.Spec.Ordinals = spec.Ordinals
	}
	return nil
}

func mergeJob(manifest *batch_v1.Job, spec *v1alpha1.JobSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}

	if spec.Parallelism != nil {
		manifest.Spec.Parallelism = spec.Parallelism
	}
	if spec.Completions != nil {
		manifest.Spec.Completions = spec.Completions
	}
	if spec.ActiveDeadlineSeconds != nil {
		manifest.Spec.ActiveDeadlineSeconds = spec.ActiveDeadlineSeconds
	}
	if spec.PodFailurePolicy != nil {
		manifest.Spec.PodFailurePolicy = spec.PodFailurePolicy
	}
	if spec.BackoffLimit != nil {
		manifest.Spec.BackoffLimit = spec.BackoffLimit
	}
	if spec.Selector != nil {
		manifest.Spec.Selector = spec.Selector
	}
	if spec.ManualSelector != nil {
		manifest.Spec.ManualSelector = spec.ManualSelector
	}
	if spec.TTLSecondsAfterFinished != nil {
		manifest.Spec.TTLSecondsAfterFinished = spec.TTLSecondsAfterFinished
	}
	if spec.CompletionMode != nil {
		manifest.Spec.CompletionMode = spec.CompletionMode
	}
	if spec.Suspend != nil {
		manifest.Spec.Suspend = spec.Suspend
	}
	return nil
}

func mergeCronJob(manifest *batch_v1.CronJob, spec *v1alpha1.CronJobSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)

	var merge *v1alpha1.JobTemplateSpec
	if spec != nil {
		merge = spec.JobTemplate
	}
	err := mergeJobTemplate(manifest.Spec.JobTemplate, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}

	if spec.Schedule != nil {
		manifest.Spec.Schedule = *spec.Schedule
	}
	if spec.TimeZone != nil {
		manifest.Spec.TimeZone = spec.TimeZone
	}
	if spec.StartingDeadlineSeconds != nil {
		manifest.Spec.StartingDeadlineSeconds = spec.StartingDeadlineSeconds
	}
	if spec.ConcurrencyPolicy != nil {
		manifest.Spec.ConcurrencyPolicy = *spec.ConcurrencyPolicy
	}
	if spec.Suspend != nil {
		manifest.Spec.Suspend = spec.Suspend
	}
	if spec.SuccessfulJobsHistoryLimit != nil {
		manifest.Spec.SuccessfulJobsHistoryLimit = spec.SuccessfulJobsHistoryLimit
	}
	if spec.FailedJobsHistoryLimit != nil {
		manifest.Spec.FailedJobsHistoryLimit = spec.FailedJobsHistoryLimit
	}
	return nil
}

func mergeReplicaSet(manifest *app_v1.ReplicaSet, spec *v1alpha1.ReplicaSetSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}
	if spec.Replicas != nil {
		manifest.Spec.Replicas = spec.Replicas
	}
	if spec.MinReadySeconds != nil {
		manifest.Spec.MinReadySeconds = *spec.MinReadySeconds
	}
	if spec.Selector != nil {
		manifest.Spec.Selector = spec.Selector
	}

	return nil
}

func mergePodTemplate(manifest *core_v1.PodTemplateSpec, spec *v1alpha1.PodTemplateSpec, prefix, ns string, labels map[string]string) error {
	mergePodTemplateObjectMeta(&manifest.ObjectMeta, labels)
	manifest.Spec.Volumes = updateVolumes(manifest.Spec.Volumes, prefix)

	if spec == nil {
		return nil
	}

	if spec.Spec.Containers != nil {
		merged, err := MergePatchContainers(manifest.Spec.Containers, spec.Spec.Containers)
		if err != nil {
			return err
		}
		manifest.Spec.Containers = merged
	}
	if spec.Spec.InitContainers != nil {
		merged, err := MergePatchContainers(manifest.Spec.InitContainers, spec.Spec.InitContainers)
		if err != nil {
			return err
		}
		manifest.Spec.InitContainers = merged
	}
	if spec.Spec.Volumes != nil {
		toMerge := updateVolumes(spec.Spec.Volumes, prefix)
		merged, err := MergePatchVolumes(manifest.Spec.Volumes, toMerge)
		if err != nil {
			return err
		}
		manifest.Spec.Volumes = merged
	}
	if spec.Spec.DNSPolicy != nil {
		manifest.Spec.DNSPolicy = *spec.Spec.DNSPolicy
	}
	if spec.Spec.NodeSelector != nil {
		manifest.Spec.NodeSelector = *spec.Spec.NodeSelector
	}
	if spec.Spec.ServiceAccountName != nil {
		manifest.Spec.ServiceAccountName = *spec.Spec.ServiceAccountName
		updateServiceAccount(manifest, prefix)
	}
	if spec.Spec.SecurityContext != nil {
		manifest.Spec.SecurityContext = spec.Spec.SecurityContext
	}
	if spec.Spec.ImagePullSecrets != nil {
		manifest.Spec.ImagePullSecrets = spec.Spec.ImagePullSecrets
	}
	if spec.Spec.Affinity != nil {
		manifest.Spec.Affinity = spec.Spec.Affinity
	}
	if spec.Spec.Tolerations != nil {
		manifest.Spec.Tolerations = spec.Spec.Tolerations
	}
	if spec.Spec.HostAliases != nil {
		manifest.Spec.HostAliases = spec.Spec.HostAliases
	}
	if spec.Spec.PriorityClassName != nil {
		manifest.Spec.PriorityClassName = *spec.Spec.PriorityClassName
	}
	if spec.Spec.TopologySpreadConstraints != nil {
		manifest.Spec.TopologySpreadConstraints = spec.Spec.TopologySpreadConstraints
	}
	return nil
}

func mergeJobTemplate(manifest batch_v1.JobTemplateSpec, spec *v1alpha1.JobTemplateSpec, prefix, ns string, labels map[string]string) error {
	updateNameWithPrefix(prefix, &manifest.ObjectMeta)
	mergeObjectMeta(&manifest.ObjectMeta, ns, labels)
	updateServiceAccount(&manifest.Spec.Template, prefix)

	var merge *v1alpha1.PodTemplateSpec
	if spec != nil {
		merge = spec.Spec.Template
	}
	err := mergePodTemplate(&manifest.Spec.Template, merge, prefix, ns, labels)
	if err != nil {
		return err
	}

	if spec == nil {
		return nil
	}

	if spec.Spec.Parallelism != nil {
		manifest.Spec.Parallelism = spec.Spec.Parallelism
	}
	if spec.Spec.Completions != nil {
		manifest.Spec.Completions = spec.Spec.Completions
	}
	if spec.Spec.ActiveDeadlineSeconds != nil {
		manifest.Spec.ActiveDeadlineSeconds = spec.Spec.ActiveDeadlineSeconds
	}
	if spec.Spec.PodFailurePolicy != nil {
		manifest.Spec.PodFailurePolicy = spec.Spec.PodFailurePolicy
	}
	if spec.Spec.BackoffLimit != nil {
		manifest.Spec.BackoffLimit = spec.Spec.BackoffLimit
	}
	if spec.Spec.Selector != nil {
		manifest.Spec.Selector = spec.Spec.Selector
	}
	if spec.Spec.ManualSelector != nil {
		manifest.Spec.ManualSelector = spec.Spec.ManualSelector
	}
	if spec.Spec.TTLSecondsAfterFinished != nil {
		manifest.Spec.TTLSecondsAfterFinished = spec.Spec.TTLSecondsAfterFinished
	}
	if spec.Spec.CompletionMode != nil {
		manifest.Spec.CompletionMode = spec.Spec.CompletionMode
	}
	if spec.Spec.Suspend != nil {
		manifest.Spec.Suspend = spec.Spec.Suspend
	}
	return nil
}

func updateServiceAccount(template *core_v1.PodTemplateSpec, prefix string) {
	if len(template.Spec.ServiceAccountName) > 0 {
		template.Spec.ServiceAccountName = fmt.Sprintf("%s%s", prefix, template.Spec.ServiceAccountName)
		return
	}

	if len(template.Spec.DeprecatedServiceAccount) > 0 {
		template.Spec.ServiceAccountName = fmt.Sprintf("%s%s", prefix, template.Spec.DeprecatedServiceAccount)
		template.Spec.DeprecatedServiceAccount = ""
	}
}

func updateServiceName(sts *app_v1.StatefulSet, prefix string) {
	if len(sts.Spec.ServiceName) > 0 {
		sts.Spec.ServiceName = fmt.Sprintf("%s%s", prefix, sts.Spec.ServiceName)
	}
}

func updateVolumes(volumes []core_v1.Volume, prefix string) (result []core_v1.Volume) {
	for _, vol := range volumes {
		if vol.ConfigMap != nil {
			vol.ConfigMap.Name = fmt.Sprintf("%s%s", prefix, vol.ConfigMap.Name)
		}
		if vol.Secret != nil {
			vol.Secret.SecretName = fmt.Sprintf("%s%s", prefix, vol.Secret.SecretName)
		}
		result = append(result, vol)
	}
	return result
}
