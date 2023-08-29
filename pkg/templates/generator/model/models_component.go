package model

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
	"github.com/udmire/observability-operator/pkg/utils"
)

const (
	A_Selector                string = "selector"
	A_Replica                 string = "replicas"
	A_MinReadySeconds         string = "minReadySeconds"
	A_RevisionHistoryLimit    string = "revisionHistoryLimit"
	A_ProgressDeadlineSeconds string = "progressDeadlineSeconds"
	A_PodManagementPolicy     string = "podManagementPolicy"
	A_UpdateStrategy          string = "updateStrategy"

	A_Parallelism             string = "parallelism"
	A_Completions             string = "completions"
	A_ActiveDeadlineSeconds   string = "activeDeadlineSeconds"
	A_BackoffLimit            string = "backoffLimit"
	A_TTLSecondsAfterFinished string = "ttlSecondsAfterFinished"
	A_CompletionMode          string = "completionMode"
	A_Suspend                 string = "suspend"

	A_Schedule                   string = "schedule"
	A_TimeZone                   string = "timeZone"
	A_StartingDeadlineSeconds    string = "startingDeadlineSeconds"
	A_ConcurrencyPolicy          string = "concurrencyPolicy"
	A_SuccessfulJobsHistoryLimit string = "successfulJobsHistoryLimit"
	A_FailedJobsHistoryLimit     string = "failedJobsHistoryLimit"
)

type Component struct {
	GenericModel `json:",inline"`
	Common       `yaml:",inline"`
	WorkloadType WorkloadType `yaml:"type"`

	Deployment  *CompDeployment  `yaml:"deployment,omitempty"`
	DaemonSet   *CompDaemonSet   `yaml:"daemonSet,omitempty"`
	StatefulSet *CompStatefulSet `yaml:"statefulSet,omitempty"`
	ReplicaSet  *CompReplicaSet  `yaml:"replicaSet,omitempty"`
	Job         *CompJob         `yaml:"job,omitempty"`
	CronJob     *CompCronJob     `yaml:"cronJob,omitempty"`

	HPA *HPA `yaml:"hpa,omitempty"`
}

func (c *Component) NextOperation(scanner *bufio.Scanner) []generator.Operation {
	return generator.Operations
}

// 列出可选类型，并选择
func (c *Component) Create(scanner *bufio.Scanner) generator.OperationHolder {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Hpa, C_Workload}
	typ := generator.ChooseType[string](scanner, opts)

	switch typ {
	case C_Hpa:
		comp := &HPA{GenericModel: inheritCreate(c.GenericModel)}
		c.HPA = comp
		return comp
	case C_Workload:
		fmt.Printf("可选Workload类型:\n%s\n请选择: ", generator.BuildOptions[WorkloadType](WorkloadTypes))
		scanner.Scan()
		input := scanner.Text()
		parsedIdx, err := generator.ParseOptions[WorkloadType](input, WorkloadTypes)
		for err != nil {
			fmt.Print("输入错误，请重新选择: ")
			scanner.Scan()
			input = scanner.Text()
			parsedIdx, err = generator.ParseOptions[WorkloadType](input, WorkloadTypes)
		}
		wt := WorkloadTypes[parsedIdx]
		switch wt {
		case WT_CronJob:
			comp := &CompCronJob{GenericModel: inheritCreate(c.GenericModel)}
			c.CronJob = comp
			return comp
		case WT_Job:
			comp := &CompJob{GenericModel: inheritCreate(c.GenericModel)}
			c.Job = comp
			return comp
		case WT_Deployment:
			comp := &CompDeployment{GenericModel: inheritCreate(c.GenericModel)}
			c.Deployment = comp
			c.WorkloadType = WT_Deployment
			return comp
		case WT_DaemonSet:
			comp := &CompDaemonSet{GenericModel: inheritCreate(c.GenericModel)}
			c.DaemonSet = comp
			c.WorkloadType = WT_DaemonSet
			return comp
		case WT_StatefulSet:
			comp := &CompStatefulSet{GenericModel: inheritCreate(c.GenericModel)}
			c.StatefulSet = comp
			c.WorkloadType = WT_StatefulSet
			return comp
		case WT_ReplicaSet:
			comp := &CompReplicaSet{GenericModel: inheritCreate(c.GenericModel)}
			c.ReplicaSet = comp
			c.WorkloadType = WT_ReplicaSet
			return comp
		}
		c.WorkloadType = wt
	}
	return createCommonOperations(c.GenericModel, &c.Common, typ)
}

func (c *Component) Modify(scanner *bufio.Scanner) generator.OperationHolder {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Hpa, C_Workload}
	typ := generator.ChooseType[string](scanner, opts)
	holders := c.List(typ)
	if len(holders) == 0 {
		fmt.Printf("没有可修改的%s", typ)
		return nil
	}

	parsedIdx := generator.ChooseIndex(scanner, typ, holders)
	return holders[parsedIdx]
}

func (c *Component) Remove(scanner *bufio.Scanner) {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Hpa, C_Workload}
	typ := generator.ChooseType[string](scanner, opts)

	holders := c.List(typ)
	if len(holders) == 0 {
		fmt.Printf("没有可删除的%s", typ)
		return
	}

	switch typ {
	case C_Secret:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		c.Secrets = append(c.Secrets[:parsedIdx], c.Secrets[parsedIdx+1:]...)
	case C_ConfigMap:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		c.ConfigMaps = append(c.ConfigMaps[:parsedIdx], c.ConfigMaps[parsedIdx+1:]...)
	case C_Service:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		c.Services = append(c.Services[:parsedIdx], c.Services[parsedIdx+1:]...)
	case C_Role:
		c.Role = nil
	case C_RoleBinding:
		c.RoleBinding = nil
	case C_ClusterRole:
		c.ClusterRole = nil
	case C_ClusterRoleBinding:
		c.ClusterRoleBinding = nil
	case C_Ingress:
		c.Ingress = nil
	case C_ServiceAccount:
		c.ServiceAccount = nil
	case C_Hpa:
		c.HPA = nil
	case C_Workload:
		switch c.WorkloadType {
		case WT_Deployment:
			c.Deployment = nil
		case WT_DaemonSet:
			c.DaemonSet = nil
		case WT_StatefulSet:
			c.StatefulSet = nil
		case WT_ReplicaSet:
			c.ReplicaSet = nil
		case WT_Job:
			c.Job = nil
		case WT_CronJob:
			c.CronJob = nil
		}
		c.WorkloadType = ""
	}
}

func (c *Component) General(scanner *bufio.Scanner) generator.GeneralCommand {
	opts := []generator.GeneralCommand{generator.OP_Finish, generator.OP_Cancel}
	return generator.ChooseType(scanner, opts)
}

func (c *Component) List(typ string) (result []generator.OperationHolder) {
	switch typ {
	case C_Secret:
		for _, secret := range c.Secrets {
			result = append(result, secret)
		}
		return
	case C_ConfigMap:
		for _, cm := range c.ConfigMaps {
			result = append(result, cm)
		}
		return
	case C_Service:
		for _, svc := range c.Services {
			result = append(result, svc)
		}
		return
	case C_Role:
		result = append(result, c.Role)
		return
	case C_RoleBinding:
		result = append(result, c.RoleBinding)
		return
	case C_ClusterRole:
		result = append(result, c.ClusterRole)
		return
	case C_ClusterRoleBinding:
		result = append(result, c.ClusterRoleBinding)
		return
	case C_Ingress:
		result = append(result, c.Ingress)
		return
	case C_ServiceAccount:
		result = append(result, c.ServiceAccount)
		return
	case C_Hpa:
		result = append(result, c.HPA)
		return
	case C_Workload:
		switch c.WorkloadType {
		case WT_Deployment:
			result = append(result, c.Deployment)
		case WT_DaemonSet:
			result = append(result, c.DaemonSet)
		case WT_StatefulSet:
			result = append(result, c.StatefulSet)
		case WT_ReplicaSet:
			result = append(result, c.ReplicaSet)
		case WT_Job:
			result = append(result, c.Job)
		case WT_CronJob:
			result = append(result, c.CronJob)
		}
		c.WorkloadType = ""
	}
	return
}

func (c *Component) Type() string {
	return C_Component
}

func (c *Component) Args() []string {
	return []string{A_Name, A_Namespace, A_Labels}
}

func (c *Component) ParseArgs(input string) error {
	sections := strings.Split(input, " ")
	if len(sections) != len(c.Args()) {
		return fmt.Errorf("invalid args, required pattern: %s", c.ArgsExample())
	}

	if !generator.IsValidName(sections[0]) {
		return fmt.Errorf("invalid names, must match pattern: %s", generator.NamePatternString)
	}
	c.Name = sections[0]
	c.Labels[utils.ComponentLabel] = c.Name

	if sections[1] != "_" {
		c.DefaultNamespace = sections[1]
	}

	if sections[2] != "_" {
		return convertAndMergeLabels(c.Labels, sections[2])
	}

	return nil
}
func (c *Component) String() string {
	return fmt.Sprintf("%s: %s/%s", c.Type(), c.DefaultNamespace, c.Name)
}

type HPA struct {
	GenericModel `yaml:",inline"`
}

func (h *HPA) Type() string {
	return C_Hpa
}

func (h *HPA) Args() []string {
	return []string{A_Name, A_Namespace, A_Labels}
}
