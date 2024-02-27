package utils

import (
	"regexp"
	"strings"
)

const (
	DOMAIN_IMAGE_PATTERN = `^[a-zA-Z0-9-]+\.[a-zA-Z]{2,}(:[0-9]+)?$`
	IPV6_IMAGE_PATTERN   = `^\[[0-9a-fA-F:]+\](:[0-9]+)?$`
	IPV4_IMAGE_PATTERN   = `^(\d{1,3}\.){3}\d{1,3}(:[0-9]+)?$`
)

var (
	DOMAIN_IMAGE_REGEX = regexp.MustCompile(DOMAIN_IMAGE_PATTERN)
	IPV4_IMAGE_REGEX   = regexp.MustCompile(IPV4_IMAGE_PATTERN)
	IPV6_IMAGE_REGEX   = regexp.MustCompile(IPV6_IMAGE_PATTERN)
	PATHS              = "/"
)

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

func UpdateImageRegistry(registry, image string) string {
	splits := strings.Split(image, PATHS)
	if len(splits) < 2 {
		return strings.Join([]string{registry, image}, PATHS)
	}

	if DOMAIN_IMAGE_REGEX.Match([]byte(splits[0])) {
		splits[0] = registry
		return strings.Join(splits, PATHS)
	} else if IPV4_IMAGE_REGEX.Match([]byte(splits[0])) {
		splits[0] = registry
		return strings.Join(splits, PATHS)
	} else if IPV6_IMAGE_REGEX.Match([]byte(splits[0])) {
		splits[0] = registry
		return strings.Join(splits, PATHS)
	}

	return strings.Join([]string{registry, image}, PATHS)
}
