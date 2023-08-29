package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type CompDaemonSet struct {
	GenericModel         `yaml:",inline"`
	LabelSelector        map[string]string `yaml:"selector,omitempty"`
	UpdateStrategy       string            `yaml:"updateStrategy,omitempty"`
	MinReadySeconds      int32             `yaml:"minReadySeconds,omitempty"`
	RevisionHistoryLimit *int32            `yaml:"revisionHistoryLimit,omitempty"`
}

func (a *CompDaemonSet) AvailableOperations() map[generator.Operation][]string {
	return map[generator.Operation][]string{}
}

func (g *CompDaemonSet) Type() string {
	return C_Workload
}

func (a *CompDaemonSet) Args() []string {
	return []string{A_Labels, A_Selector, A_UpdateStrategy, A_MinReadySeconds, A_RevisionHistoryLimit}
}
func (a *CompDaemonSet) ArgsExample() string {
	return "label1:values;label2:value2 label1:values updateStrategy minReadySeconds RevisionHistoryLimit"
}

func (a *CompDaemonSet) ParseArgs(input string) (err error) {
	sections := strings.Split(input, " ")
	if len(sections) != len(a.Args()) {
		return fmt.Errorf("invalid args, required pattern: %s", a.ArgsExample())
	}

	if sections[0] != "_" {
		err = convertAndMergeLabels(a.Labels, sections[0])
		if err != nil {
			return err
		}
	}

	if sections[1] != "_" {
		err = convertAndMergeLabels(a.LabelSelector, sections[1])
		if err != nil {
			return err
		}
	}

	if sections[2] != "_" {
		a.UpdateStrategy, err = convertUpdateStrategy(sections[2])
		if err != nil {
			return err
		}
	}

	if sections[3] != "_" {
		val, err := strconv.Atoi(sections[3])
		if err != nil {
			return err
		}
		a.MinReadySeconds = int32(val)
	}

	if sections[4] != "_" {
		val, err := strconv.Atoi(sections[4])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.RevisionHistoryLimit = &int32Val
	}

	return nil
}

func convertUpdateStrategy(value string) (string, error) {
	if value == "RollingUpdate" || value == "OnDelete" {
		return value, nil
	}
	return "", fmt.Errorf("invalid values")
}
