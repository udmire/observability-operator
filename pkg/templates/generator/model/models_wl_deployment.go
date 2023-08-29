package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type CompDeployment struct {
	GenericModel            `yaml:",inline"`
	Selector                map[string]string `yaml:"selector,omitempty"`
	Replicas                int               `yaml:"replicas,omitempty"`
	MinReadySeconds         int32             `yaml:"minReadySeconds,omitempty"`
	RevisionHistoryLimit    *int32            `yaml:"revisionHistoryLimit,omitempty"`
	ProgressDeadlineSeconds *int32            `yaml:"progressDeadlineSeconds,omitempty"`
}

func (a *CompDeployment) AvailableOperations() map[generator.Operation][]string {
	return map[generator.Operation][]string{}
}

func (a *CompDeployment) Type() string {
	return C_Workload
}
func (a *CompDeployment) Args() []string {
	return []string{A_Labels, A_Selector, A_Replica, A_MinReadySeconds, A_RevisionHistoryLimit, A_ProgressDeadlineSeconds}
}

func (a *CompDeployment) ArgsExample() string {
	return "label1:values;label2:value2 label1:values replicas minReadySeconds RevisionHistoryLimit ProgressDeadlineSeconds"
}
func (a *CompDeployment) ParseArgs(input string) (err error) {
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
		err = convertAndMergeLabels(a.Selector, sections[1])
		if err != nil {
			return err
		}
	}

	if sections[2] != "_" {
		a.Replicas, err = strconv.Atoi(sections[2])
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

	if sections[5] != "_" {
		val, err := strconv.Atoi(sections[5])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.ProgressDeadlineSeconds = &int32Val
	}

	return nil
}
