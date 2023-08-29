package model

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type CompCronJob struct {
	GenericModel `yaml:",inline"`
	JobTemplate  *CompJob

	Schedule                   string  `json:"schedule"`
	TimeZone                   *string `json:"timeZone,omitempty"`
	StartingDeadlineSeconds    *int64  `json:"startingDeadlineSeconds,omitempty"`
	ConcurrencyPolicy          string  `json:"concurrencyPolicy,omitempty"`
	SuccessfulJobsHistoryLimit *int32  `json:"successfulJobsHistoryLimit,omitempty"`
	FailedJobsHistoryLimit     *int32  `json:"failedJobsHistoryLimit,omitempty"`
}

func (a *CompCronJob) AvailableOperations() map[generator.Operation][]string {
	return map[generator.Operation][]string{}
}

func (g *CompCronJob) Type() string {
	return ""
}

func (a *CompCronJob) Args() []string {
	return []string{A_Labels, A_Schedule, A_TimeZone, A_StartingDeadlineSeconds, A_ConcurrencyPolicy, A_SuccessfulJobsHistoryLimit, A_FailedJobsHistoryLimit}
}

func (a *CompCronJob) ArgsExample() string {
	return "label1:values;label2:value2 schedule timeZone startingDeadlineSeconds concurrencyPolicy successfulJobsHistoryLimit failedJobsHistoryLimit"
}

func (a *CompCronJob) ParseArgs(input string) (err error) {
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

	a.Schedule = sections[1]

	if sections[2] != "_" {
		a.TimeZone = &sections[2]
	}

	if sections[3] != "_" {
		val, err := strconv.Atoi(sections[3])
		if err != nil {
			return err
		}
		int64Val := int64(val)
		a.StartingDeadlineSeconds = &int64Val
	}

	if sections[4] != "_" {
		val, err := convertConcurrencyPolicy(sections[4])
		if err != nil {
			return err
		}
		a.ConcurrencyPolicy = val
	}

	if sections[5] != "_" {
		val, err := strconv.Atoi(sections[5])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.SuccessfulJobsHistoryLimit = &int32Val
	}

	if sections[6] != "_" {
		val, err := strconv.Atoi(sections[6])
		if err != nil {
			return err
		}
		int32Val := int32(val)
		a.FailedJobsHistoryLimit = &int32Val
	}
	return nil
}

func convertConcurrencyPolicy(value string) (string, error) {
	if value == "Allow" || value == "Forbid" || value == "Replace" {
		return value, nil
	}
	return "", fmt.Errorf("invalid values")
}
