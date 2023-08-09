package utils

func AppInstanceLabels(instance, template, version string) map[string]string {
	return map[string]string{
		AppLabel:       template,
		InstanceLabel:  instance,
		VersionLabel:   version,
		ManagedByLabel: DefaultManagedByValue,
	}
}

func ComponentLabels(instance, template, version, component string) map[string]string {
	ils := AppInstanceLabels(instance, template, version)
	ils[ComponentLabel] = component
	return ils
}
