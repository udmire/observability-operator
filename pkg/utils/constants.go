package utils

const (
	AppLabel       = "app.kubernetes.io/name"
	InstanceLabel  = "app.kubernetes.io/instance"
	ManagedByLabel = "app.kubernetes.io/managed-by"
	ComponentLabel = "app.kubernetes.io/component"
	VersionLabel   = "app.kubernetes.io/version"
	PartOfLabel    = "app.kubernetes.io/part-of"

	DefaultManagedByValue = "observability-operator"
)

var SelectorIgnoredLabels = []string{ManagedByLabel, VersionLabel, PartOfLabel}

func ShouldIgnoredInSelector(key string) bool {
	for _, v := range SelectorIgnoredLabels {
		if v == key {
			return true
		}
	}
	return false
}
