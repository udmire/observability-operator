package manifest

import core_v1 "k8s.io/api/core/v1"

type CapsuleType string

const (
	ConfigmapType CapsuleType = "configmap"
	SecretType    CapsuleType = "secret"
)

const (
	CapsuleFile string = ".capsule"
)

type Capsule struct {
	Name         string      `json:"name"`
	Type         CapsuleType `json:"type,omitempty"`
	CapsuleItems `json:",inline"`
}

type CapsuleItems struct {
	Items        map[string]string `json:"items,omitempty"`
	DynamicItems *string           `json:"dynamics,omitempty"`
}

type Manifest struct {
	Name       string
	ConfigMaps []*core_v1.ConfigMap
	Secrets    []*core_v1.Secret
}

type CapsuleManifests struct {
	Manifest

	CompsManifests []*CompManifests
}

type CompManifests struct {
	Manifest
}
