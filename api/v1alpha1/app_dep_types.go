package v1alpha1

type AppDepsSpec struct {
	Capsules map[string]CapsuleSpec `json:"capsules,omitempty"`
}
