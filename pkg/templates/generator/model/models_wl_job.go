package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type CompJob struct {
	GenericModel            `yaml:",inline"`
	Selector                map[string]string `yaml:"selector,omitempty"`
	Parallelism             *int32            `yaml:"parallelism,omitempty"`
	Completions             *int32            `yaml:"completions,omitempty"`
	ActiveDeadlineSeconds   *int64            `yaml:"activeDeadlineSeconds,omitempty"`
	BackoffLimit            *int32            `yaml:"backoffLimit,omitempty"`
	TTLSecondsAfterFinished *int32            `yaml:"ttlSecondsAfterFinished,omitempty"`
	CompletionMode          string            `yaml:"completionMode,omitempty"`
	Suspend                 *bool             `yaml:"suspend,omitempty"`
}

func (a *CompJob) AvailableOperations() map[generator.Operation][]string {
	return map[generator.Operation][]string{}
}
func (g *CompJob) Type() string {
	return ""
}
func (a *CompJob) Args() []string {
	return []string{A_Labels, A_Selector, A_Parallelism, A_Completions, A_ActiveDeadlineSeconds, A_BackoffLimit, A_TTLSecondsAfterFinished, A_CompletionMode, A_Suspend}
}

func (a *CompJob) ArgsExample() string {
	return "label1:values;label2:value2 label1:values parallelism completions activeDeadlineSeconds backoffLimit ttlSecondsAfterFinished completionMode suspend"
}
func (a *CompJob) ParseArgs(input string) (err error) {
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
		val, err := strconv.Atoi(sections[2])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.Parallelism = &int32Val
	}

	if sections[3] != "_" {
		val, err := strconv.Atoi(sections[3])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.Completions = &int32Val
	}

	if sections[4] != "_" {
		val, err := strconv.Atoi(sections[4])
		if err != nil {
			return err
		}
		int64Val := int64(val)
		a.ActiveDeadlineSeconds = &int64Val
	}

	if sections[5] != "_" {
		val, err := strconv.Atoi(sections[5])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.BackoffLimit = &int32Val
	}
	if sections[6] != "_" {
		val, err := strconv.Atoi(sections[6])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.TTLSecondsAfterFinished = &int32Val
	}
	if sections[7] != "_" {
		val, err := convertCompletionMode(sections[7])
		if err != nil {
			return err
		}
		a.CompletionMode = val
	}
	if sections[8] != "_" {
		val, err := strconv.ParseBool(sections[8])
		if err != nil {
			return err
		}
		a.Suspend = &val
	}

	return nil
}

func convertCompletionMode(value string) (string, error) {
	if value == "NonIndexed" || value == "Indexed" {
		return value, nil
	}
	return "", fmt.Errorf("invalid values")
}
