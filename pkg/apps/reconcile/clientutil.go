package reconcile

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	apps_v1 "k8s.io/api/apps/v1"
	batch_v1 "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
	networking_v1 "k8s.io/api/networking/v1"
	rbac_v1 "k8s.io/api/rbac/v1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

var invalidDNS1123Characters = regexp.MustCompile("[^-a-z0-9]+")

// SanitizeVolumeName ensures that the given volume name is a valid DNS-1123 label
// accepted by Kubernetes.
//
// Copied from github.com/prometheus-operator/prometheus-operator/pkg/k8sutil.
func SanitizeVolumeName(name string) string {
	name = strings.ToLower(name)
	name = invalidDNS1123Characters.ReplaceAllString(name, "-")
	if len(name) > validation.DNS1123LabelMaxLength {
		name = name[0:validation.DNS1123LabelMaxLength]
	}
	return strings.Trim(name, "-")
}

// CreateOrUpdateSecret applies the given secret against the client.
func CreateOrUpdateSecret(ctx context.Context, c client.Client, s *v1.Secret) error {
	var exist v1.Secret
	err := c.Get(ctx, client.ObjectKeyFromObject(s), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing service: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, s)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
	} else {
		s.ResourceVersion = exist.ResourceVersion
		s.SetOwnerReferences(mergeOwnerReferences(s.GetOwnerReferences(), exist.GetOwnerReferences()))
		s.SetLabels(mergeMaps(s.Labels, exist.Labels))
		s.SetAnnotations(mergeMaps(s.Annotations, exist.Annotations))

		err := c.Update(ctx, s)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update service: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateConfigMap applies the given secret against the client.
func CreateOrUpdateConfigMap(ctx context.Context, c client.Client, s *v1.ConfigMap) error {
	var exist v1.ConfigMap
	err := c.Get(ctx, client.ObjectKeyFromObject(s), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing ConfigMap: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, s)
		if err != nil {
			return fmt.Errorf("failed to create ConfigMap: %w", err)
		}
	} else {
		s.ResourceVersion = exist.ResourceVersion
		s.SetOwnerReferences(mergeOwnerReferences(s.GetOwnerReferences(), exist.GetOwnerReferences()))
		s.SetLabels(mergeMaps(s.Labels, exist.Labels))
		s.SetAnnotations(mergeMaps(s.Annotations, exist.Annotations))

		err := c.Update(ctx, s)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update ConfigMap: %w", err)
		}
	}

	return nil
}

func CreateOrUpdateIngress(ctx context.Context, c client.Client, i *networking_v1.Ingress) error {
	var exist networking_v1.Ingress
	err := c.Get(ctx, client.ObjectKeyFromObject(i), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing Ingress: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, i)
		if err != nil {
			return fmt.Errorf("failed to create Ingress: %w", err)
		}
	} else {
		i.ResourceVersion = exist.ResourceVersion
		i.SetOwnerReferences(mergeOwnerReferences(i.GetOwnerReferences(), exist.GetOwnerReferences()))
		i.SetLabels(mergeMaps(i.Labels, exist.Labels))
		i.SetAnnotations(mergeMaps(i.Annotations, exist.Annotations))

		err := c.Update(ctx, i)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update Ingress: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateClusterRole applies the given clusterRole against the client.
func CreateOrUpdateClusterRole(ctx context.Context, c client.Client, clusterRole *rbac_v1.ClusterRole) error {
	var exist rbac_v1.ClusterRole
	err := c.Get(ctx, client.ObjectKeyFromObject(clusterRole), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing clusterRole: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, clusterRole)
		if err != nil {
			return fmt.Errorf("failed to create clusterRole: %w", err)
		}
	} else {
		clusterRole.ResourceVersion = exist.ResourceVersion
		clusterRole.SetOwnerReferences(mergeOwnerReferences(clusterRole.GetOwnerReferences(), exist.GetOwnerReferences()))
		clusterRole.SetLabels(mergeMaps(clusterRole.Labels, exist.Labels))
		clusterRole.SetAnnotations(mergeMaps(clusterRole.Annotations, exist.Annotations))

		err := c.Update(ctx, clusterRole)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update clusterRole: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateClusterRoleBinding applies the given crb against the client.
func CreateOrUpdateClusterRoleBinding(ctx context.Context, c client.Client, crb *rbac_v1.ClusterRoleBinding) error {
	var exist v1.Service
	err := c.Get(ctx, client.ObjectKeyFromObject(crb), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing clusterRoleBinding: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, crb)
		if err != nil {
			return fmt.Errorf("failed to create clusterRoleBinding: %w", err)
		}
	} else {
		crb.ResourceVersion = exist.ResourceVersion
		crb.SetOwnerReferences(mergeOwnerReferences(crb.GetOwnerReferences(), exist.GetOwnerReferences()))
		crb.SetLabels(mergeMaps(crb.Labels, exist.Labels))
		crb.SetAnnotations(mergeMaps(crb.Annotations, exist.Annotations))

		err := c.Update(ctx, crb)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update clusterRoleBinding: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateRole applies the given role against the client.
func CreateOrUpdateRole(ctx context.Context, c client.Client, role *rbac_v1.Role) error {
	var exist rbac_v1.Role
	err := c.Get(ctx, client.ObjectKeyFromObject(role), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing role: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, role)
		if err != nil {
			return fmt.Errorf("failed to create role: %w", err)
		}
	} else {
		role.ResourceVersion = exist.ResourceVersion
		role.SetOwnerReferences(mergeOwnerReferences(role.GetOwnerReferences(), exist.GetOwnerReferences()))
		role.SetLabels(mergeMaps(role.Labels, exist.Labels))
		role.SetAnnotations(mergeMaps(role.Annotations, exist.Annotations))

		err := c.Update(ctx, role)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update service: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateRoleBinding applies the given rolebinding against the client.
func CreateOrUpdateRoleBinding(ctx context.Context, c client.Client, roleBinding *rbac_v1.RoleBinding) error {
	var exist rbac_v1.RoleBinding
	err := c.Get(ctx, client.ObjectKeyFromObject(roleBinding), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing roleBinding: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, roleBinding)
		if err != nil {
			return fmt.Errorf("failed to create roleBinding: %w", err)
		}
	} else {
		roleBinding.ResourceVersion = exist.ResourceVersion
		roleBinding.SetOwnerReferences(mergeOwnerReferences(roleBinding.GetOwnerReferences(), exist.GetOwnerReferences()))
		roleBinding.SetLabels(mergeMaps(roleBinding.Labels, exist.Labels))
		roleBinding.SetAnnotations(mergeMaps(roleBinding.Annotations, exist.Annotations))

		err := c.Update(ctx, roleBinding)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update roleBinding: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateService applies the given svc against the client.
func CreateOrUpdateService(ctx context.Context, c client.Client, svc *v1.Service) error {
	var exist v1.Service
	err := c.Get(ctx, client.ObjectKeyFromObject(svc), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing service: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, svc)
		if err != nil {
			return fmt.Errorf("failed to create service: %w", err)
		}
	} else {
		svc.ResourceVersion = exist.ResourceVersion
		svc.Spec.IPFamilies = exist.Spec.IPFamilies
		svc.SetOwnerReferences(mergeOwnerReferences(svc.GetOwnerReferences(), exist.GetOwnerReferences()))
		svc.SetLabels(mergeMaps(svc.Labels, exist.Labels))
		svc.SetAnnotations(mergeMaps(svc.Annotations, exist.Annotations))

		err := c.Update(ctx, svc)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update service: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateServiceAccount applies the given sa against the client.
func CreateOrUpdateServiceAccount(ctx context.Context, c client.Client, sa *v1.ServiceAccount) error {
	var exist v1.ServiceAccount
	err := c.Get(ctx, client.ObjectKeyFromObject(sa), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing serviceAccount: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, sa)
		if err != nil {
			return fmt.Errorf("failed to create serviceAccount: %w", err)
		}
	} else {
		sa.ResourceVersion = exist.ResourceVersion
		sa.SetOwnerReferences(mergeOwnerReferences(sa.GetOwnerReferences(), exist.GetOwnerReferences()))
		sa.SetLabels(mergeMaps(sa.Labels, exist.Labels))
		sa.SetAnnotations(mergeMaps(sa.Annotations, exist.Annotations))

		err := c.Update(ctx, sa)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update serviceAccount: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateEndpoints applies the given eps against the client.
func CreateOrUpdateEndpoints(ctx context.Context, c client.Client, eps *v1.Endpoints) error {
	var exist v1.Endpoints
	err := c.Get(ctx, client.ObjectKeyFromObject(eps), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing endpoints: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, eps)
		if err != nil {
			return fmt.Errorf("failed to create endpoints: %w", err)
		}
	} else {
		eps.ResourceVersion = exist.ResourceVersion
		eps.SetOwnerReferences(mergeOwnerReferences(eps.GetOwnerReferences(), exist.GetOwnerReferences()))
		eps.SetLabels(mergeMaps(eps.Labels, exist.Labels))
		eps.SetAnnotations(mergeMaps(eps.Annotations, exist.Annotations))

		err := c.Update(ctx, eps)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update endpoints: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateStatefulSet applies the given StatefulSet against the client.
func CreateOrUpdateStatefulSet(ctx context.Context, c client.Client, ss *apps_v1.StatefulSet) error {
	var exist apps_v1.StatefulSet
	err := c.Get(ctx, client.ObjectKeyFromObject(ss), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing statefulset: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, ss)
		if err != nil {
			return fmt.Errorf("failed to create statefulset: %w", err)
		}
	} else {
		ss.ResourceVersion = exist.ResourceVersion
		ss.SetOwnerReferences(mergeOwnerReferences(ss.GetOwnerReferences(), exist.GetOwnerReferences()))
		ss.SetLabels(mergeMaps(ss.Labels, exist.Labels))
		ss.SetAnnotations(mergeMaps(ss.Annotations, exist.Annotations))

		err := c.Update(ctx, ss)
		// Statefulsets have a large number of fields that are immutable after creation,
		// so we sometimes need to delete and recreate.
		// We should be mindful when making changes to try and avoid this when possible.
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			ss.ResourceVersion = ""

			// do a quicker deletion of the old statefulset to minimize downtime before we spin up new pods
			err = c.Delete(ctx, ss, client.GracePeriodSeconds(5))
			if err != nil {
				return fmt.Errorf("failed to update statefulset when deleting old statefulset: %w", err)
			}
			err = c.Create(ctx, ss)
			if err != nil {
				return fmt.Errorf("failed to update statefulset when creating replacement statefulset: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update statefulset: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateDaemonSet applies the given DaemonSet against the client.
func CreateOrUpdateDaemonSet(ctx context.Context, c client.Client, ss *apps_v1.DaemonSet) error {
	var exist apps_v1.DaemonSet
	err := c.Get(ctx, client.ObjectKeyFromObject(ss), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing daemonset: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, ss)
		if err != nil {
			return fmt.Errorf("failed to create daemonset: %w", err)
		}
	} else {
		ss.ResourceVersion = exist.ResourceVersion
		ss.SetOwnerReferences(mergeOwnerReferences(ss.GetOwnerReferences(), exist.GetOwnerReferences()))
		ss.SetLabels(mergeMaps(ss.Labels, exist.Labels))
		ss.SetAnnotations(mergeMaps(ss.Annotations, exist.Annotations))

		err := c.Update(ctx, ss)
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			ss.ResourceVersion = ""

			err = c.Delete(ctx, ss)
			if err != nil {
				return fmt.Errorf("failed to update daemonset: deleting old daemonset: %w", err)
			}
			err = c.Create(ctx, ss)
			if err != nil {
				return fmt.Errorf("failed to update daemonset: creating new deamonset: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update daemonset: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateDeployment applies the given DaemonSet against the client.
func CreateOrUpdateDeployment(ctx context.Context, c client.Client, d *apps_v1.Deployment) error {
	var exist apps_v1.Deployment
	err := c.Get(ctx, client.ObjectKeyFromObject(d), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing Deployment: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, d)
		if err != nil {
			return fmt.Errorf("failed to create Deployment: %w", err)
		}
	} else {
		d.ResourceVersion = exist.ResourceVersion
		d.SetOwnerReferences(mergeOwnerReferences(d.GetOwnerReferences(), exist.GetOwnerReferences()))
		d.SetLabels(mergeMaps(d.Labels, exist.Labels))
		d.SetAnnotations(mergeMaps(d.Annotations, exist.Annotations))

		err := c.Update(ctx, d)
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			d.ResourceVersion = ""

			err = c.Delete(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update Deployment: deleting old Deployment: %w", err)
			}
			err = c.Create(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update Deployment: creating new Deployment: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update Deployment: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateReplicaSet applies the given ReplicaSet against the client.
func CreateOrUpdateReplicaSet(ctx context.Context, c client.Client, d *apps_v1.ReplicaSet) error {
	var exist apps_v1.ReplicaSet
	err := c.Get(ctx, client.ObjectKeyFromObject(d), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing ReplicaSet: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, d)
		if err != nil {
			return fmt.Errorf("failed to create ReplicaSet: %w", err)
		}
	} else {
		d.ResourceVersion = exist.ResourceVersion
		d.SetOwnerReferences(mergeOwnerReferences(d.GetOwnerReferences(), exist.GetOwnerReferences()))
		d.SetLabels(mergeMaps(d.Labels, exist.Labels))
		d.SetAnnotations(mergeMaps(d.Annotations, exist.Annotations))

		err := c.Update(ctx, d)
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			d.ResourceVersion = ""

			err = c.Delete(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update ReplicaSet: deleting old ReplicaSet: %w", err)
			}
			err = c.Create(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update ReplicaSet: creating new ReplicaSet: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update ReplicaSet: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateJob applies the given Job against the client.
func CreateOrUpdateJob(ctx context.Context, c client.Client, d *batch_v1.Job) error {
	var exist batch_v1.Job
	err := c.Get(ctx, client.ObjectKeyFromObject(d), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing Job: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, d)
		if err != nil {
			return fmt.Errorf("failed to create Job: %w", err)
		}
	} else {
		d.ResourceVersion = exist.ResourceVersion
		d.SetOwnerReferences(mergeOwnerReferences(d.GetOwnerReferences(), exist.GetOwnerReferences()))
		d.SetLabels(mergeMaps(d.Labels, exist.Labels))
		d.SetAnnotations(mergeMaps(d.Annotations, exist.Annotations))

		err := c.Update(ctx, d)
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			d.ResourceVersion = ""

			err = c.Delete(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update Job: deleting old Job: %w", err)
			}
			err = c.Create(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update Job: creating new Job: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update Job: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateCronJob applies the given CronJob against the client.
func CreateOrUpdateCronJob(ctx context.Context, c client.Client, d *batch_v1.CronJob) error {
	var exist batch_v1.CronJob
	err := c.Get(ctx, client.ObjectKeyFromObject(d), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing CronJob: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, d)
		if err != nil {
			return fmt.Errorf("failed to create CronJob: %w", err)
		}
	} else {
		d.ResourceVersion = exist.ResourceVersion
		d.SetOwnerReferences(mergeOwnerReferences(d.GetOwnerReferences(), exist.GetOwnerReferences()))
		d.SetLabels(mergeMaps(d.Labels, exist.Labels))
		d.SetAnnotations(mergeMaps(d.Annotations, exist.Annotations))

		err := c.Update(ctx, d)
		if k8s_errors.IsNotAcceptable(err) || k8s_errors.IsInvalid(err) {
			// Resource version should only be set when updating
			d.ResourceVersion = ""

			err = c.Delete(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update CronJob: deleting old CronJob: %w", err)
			}
			err = c.Create(ctx, d)
			if err != nil {
				return fmt.Errorf("failed to update CronJob: creating new CronJob: %w", err)
			}
		} else if err != nil {
			return fmt.Errorf("failed to update Deployment: %w", err)
		}
	}

	return nil
}

func mergeOwnerReferences(new, old []meta_v1.OwnerReference) []meta_v1.OwnerReference {
	existing := make(map[types.UID]bool)
	for _, ref := range old {
		existing[ref.UID] = true
	}
	for _, ref := range new {
		if _, ok := existing[ref.UID]; !ok {
			old = append(old, ref)
		}
	}
	return old
}

func mergeMaps(new, old map[string]string) map[string]string {
	if old == nil {
		old = make(map[string]string, len(new))
	}
	for k, v := range new {
		old[k] = v
	}
	return old
}
