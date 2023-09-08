package client

import (
	"context"
	"fmt"

	"github.com/udmire/observability-operator/api/v1alpha1"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// CreateOrUpdateCapsule applies the given capsule against the client.
func CreateOrUpdateCapsule(ctx context.Context, c client.Client, s *v1alpha1.Capsule) error {
	var exist v1alpha1.Capsule
	err := c.Get(ctx, client.ObjectKeyFromObject(s), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing capsule: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, s)
		if err != nil {
			return fmt.Errorf("failed to create capsule: %w", err)
		}
	} else {
		s.ResourceVersion = exist.ResourceVersion
		s.SetOwnerReferences(mergeOwnerReferences(s.GetOwnerReferences(), exist.GetOwnerReferences()))
		s.SetLabels(mergeMaps(s.Labels, exist.Labels))
		s.SetAnnotations(mergeMaps(s.Annotations, exist.Annotations))

		err := c.Update(ctx, s)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update capsule: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateApps applies the given apps against the client.
func CreateOrUpdateApps(ctx context.Context, c client.Client, s *v1alpha1.Apps) error {
	var exist v1alpha1.Apps
	err := c.Get(ctx, client.ObjectKeyFromObject(s), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing apps: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, s)
		if err != nil {
			return fmt.Errorf("failed to create apps: %w", err)
		}
	} else {
		s.ResourceVersion = exist.ResourceVersion
		s.SetOwnerReferences(mergeOwnerReferences(s.GetOwnerReferences(), exist.GetOwnerReferences()))
		s.SetLabels(mergeMaps(s.Labels, exist.Labels))
		s.SetAnnotations(mergeMaps(s.Annotations, exist.Annotations))

		err := c.Update(ctx, s)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update apps: %w", err)
		}
	}

	return nil
}

// CreateOrUpdateExporters applies the given exporters against the client.
func CreateOrUpdateExporters(ctx context.Context, c client.Client, s *v1alpha1.Exporters) error {
	var exist v1alpha1.Exporters
	err := c.Get(ctx, client.ObjectKeyFromObject(s), &exist)
	if err != nil && !k8s_errors.IsNotFound(err) {
		return fmt.Errorf("failed to retrieve existing exporters: %w", err)
	}

	if k8s_errors.IsNotFound(err) {
		err := c.Create(ctx, s)
		if err != nil {
			return fmt.Errorf("failed to create exporters: %w", err)
		}
	} else {
		s.ResourceVersion = exist.ResourceVersion
		s.SetOwnerReferences(mergeOwnerReferences(s.GetOwnerReferences(), exist.GetOwnerReferences()))
		s.SetLabels(mergeMaps(s.Labels, exist.Labels))
		s.SetAnnotations(mergeMaps(s.Annotations, exist.Annotations))

		err := c.Update(ctx, s)
		if err != nil && !k8s_errors.IsNotFound(err) {
			return fmt.Errorf("failed to update exporters: %w", err)
		}
	}

	return nil
}
