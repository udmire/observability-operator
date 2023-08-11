package utils

import "strings"

func IsNewerThan(version, previous string) bool {
	version = strings.TrimLeft(strings.ToLower(version), "v")
	previous = strings.TrimLeft(strings.ToLower(previous), "v")

	versionParts := strings.Split(version, ".")
	previousParts := strings.Split(previous, ".")

	length := len(versionParts)
	if len(previousParts) < length {
		length = len(previousParts)
	}

	for i := 0; i < length; i++ {
		if versionParts[i] > previousParts[i] {
			return true
		} else if versionParts[i] < previousParts[i] {
			return false
		}
		continue
	}

	return len(versionParts) >= len(previousParts)
}
