package specs

import (
	"encoding/json"
	"fmt"

	core_v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/util/strategicpatch"
)

func MergePatchContainers(base, patches []core_v1.Container) ([]core_v1.Container, error) {
	return MergePatch[core_v1.Container, func(core_v1.Container) string, func() core_v1.Container](
		base, patches, func(c core_v1.Container) string { return c.Name }, func() core_v1.Container { return core_v1.Container{} },
	)
}

func MergePatchVolumes(base, patches []core_v1.Volume) ([]core_v1.Volume, error) {
	return MergePatch[core_v1.Volume, func(core_v1.Volume) string, func() core_v1.Volume](
		base, patches, func(v core_v1.Volume) string { return v.Name }, func() core_v1.Volume { return core_v1.Volume{} },
	)
}

func MergePatch[T any, I func(T) string, N func() T](base, patches []T, identifier I, create N) ([]T, error) {
	var out []T

	// map of containers that still need to be patched by name
	toPatch := make(map[string]T)
	for _, patch := range patches {
		id := identifier(patch)
		toPatch[id] = patch
	}

	for _, b := range base {
		id := identifier(b)
		// If we have a patch result, iterate over each container and try and calculate the patch
		if patch, ok := toPatch[id]; ok {
			// Get the json for the container and the patch
			baseBytes, err := json.Marshal(b)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal json for base %s: %w", id, err)
			}
			patchBytes, err := json.Marshal(patch)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal json for patch container %s: %w", id, err)
			}

			// Calculate the patch result
			jsonResult, err := strategicpatch.StrategicMergePatch(baseBytes, patchBytes, create())
			if err != nil {
				return nil, fmt.Errorf("failed to generate merge patch for %s: %w", id, err)
			}
			var patchResult T
			if err := json.Unmarshal(jsonResult, &patchResult); err != nil {
				return nil, fmt.Errorf("failed to unmarshal merged container %s: %w", id, err)
			}

			// Add the patch result and remove the corresponding key from the to do list
			out = append(out, patchResult)
			delete(toPatch, id)
		} else {
			// This container didn't need to be patched
			out = append(out, b)
		}
	}

	// Iterate over the patches and add all the containers that were not previously part of a patch result
	for _, patch := range patches {
		id := identifier(patch)
		if _, ok := toPatch[id]; ok {
			out = append(out, patch)
		}
	}

	return out, nil
}
