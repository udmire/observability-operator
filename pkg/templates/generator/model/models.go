package model

import (
	"bufio"
	"fmt"
	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
	"github.com/udmire/observability-operator/pkg/utils"
)

type WorkloadType string

const (
	WT_Deployment  WorkloadType = "deployment"
	WT_DaemonSet   WorkloadType = "daemonset"
	WT_StatefulSet WorkloadType = "statefulset"
	WT_ReplicaSet  WorkloadType = "replicaset"
	WT_Job         WorkloadType = "job"
	WT_CronJob     WorkloadType = "cronjob"
)

var (
	WorkloadTypes = []WorkloadType{
		WT_Deployment,
		WT_DaemonSet,
		WT_StatefulSet,
		WT_ReplicaSet,
		WT_Job,
		WT_CronJob,
	}
)

const (
	C_ConfigMap          string = "configmap"
	C_Secret             string = "secret"
	C_Service            string = "service"
	C_ClusterRole        string = "clusterrole"
	C_ClusterRoleBinding string = "clusterrolebinding"
	C_Role               string = "role"
	C_RoleBinding        string = "rolebinding"
	C_ServiceAccount     string = "serviceaccount"
	C_Ingress            string = "ingress"

	C_APP       string = "app"
	C_Component string = "component"
	C_Workload  string = "workload"
	C_Hpa       string = "hpa"
)

const (
	A_Name      string = "name"
	A_Version   string = "version"
	A_Namespace string = "namespace"
	A_Labels    string = "labels"
)

type Base struct {
	Name   string            `yaml:"name"`
	Labels map[string]string `yaml:"labels"`
}

type App struct {
	GenericModel `json:",inline"`
	Common       `yaml:",inline"`
	Version      string       `yaml:"version"`
	Components   []*Component `yaml:"components,omitempty"`
}

// 列出可选类型，并选择
func (a *App) NextOperation(scanner *bufio.Scanner) []generator.Operation {
	return generator.Operations
}

// 列出可选类型，并选择
func (a *App) Create(scanner *bufio.Scanner) generator.OperationHolder {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Component}
	typ := generator.ChooseType[string](scanner, opts)
	if typ == C_Component {
		comp := &Component{GenericModel: inheritCreate(a.GenericModel)}
		a.Components = append(a.Components, comp)
		return comp
	} else {
		return createCommonOperations(a.GenericModel, &a.Common, typ)
	}
}

func (a *App) Modify(scanner *bufio.Scanner) generator.OperationHolder {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Component}
	typ := generator.ChooseType[string](scanner, opts)
	holders := a.List(typ)
	if len(holders) == 0 {
		fmt.Printf("没有可修改的%s", typ)
		return nil
	}

	parsedIdx := generator.ChooseIndex(scanner, typ, holders)

	return holders[parsedIdx]
}

func (a *App) Remove(scanner *bufio.Scanner) {
	opts := []string{C_ConfigMap, C_Secret, C_Service, C_ClusterRole, C_ClusterRoleBinding, C_Role, C_RoleBinding, C_ServiceAccount, C_Ingress, C_Component}
	typ := generator.ChooseType[string](scanner, opts)

	holders := a.List(typ)
	if len(holders) == 0 {
		fmt.Printf("没有可删除的%s", typ)
		return
	}

	switch typ {
	case C_Secret:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		a.Secrets = append(a.Secrets[:parsedIdx], a.Secrets[parsedIdx+1:]...)
	case C_ConfigMap:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		a.ConfigMaps = append(a.ConfigMaps[:parsedIdx], a.ConfigMaps[parsedIdx+1:]...)
	case C_Service:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		a.Services = append(a.Services[:parsedIdx], a.Services[parsedIdx+1:]...)
	case C_Component:
		parsedIdx := generator.ChooseIndex(scanner, typ, holders)
		a.Components = append(a.Components[:parsedIdx], a.Components[parsedIdx+1:]...)
	case C_Role:
		a.Role = nil
	case C_RoleBinding:
		a.RoleBinding = nil
	case C_ClusterRole:
		a.ClusterRole = nil
	case C_ClusterRoleBinding:
		a.ClusterRoleBinding = nil
	case C_Ingress:
		a.Ingress = nil
	case C_ServiceAccount:
		a.ServiceAccount = nil
	}

	return
}

func (a *App) General(scanner *bufio.Scanner) generator.GeneralCommand {
	opts := []generator.GeneralCommand{generator.OP_Finish, generator.OP_Generate}
	typ := generator.ChooseType[generator.GeneralCommand](scanner, opts)
	switch typ {
	case generator.OP_Finish:
		fmt.Printf("是否结束应用构建(Y/n): ")
		input := scanner.Text()

		for strings.ToUpper(input) != "Y" && strings.ToUpper(input) != "N" {
			fmt.Print("输入错误，请重新输入(Y/n): ")
			scanner.Scan()
			input = scanner.Text()
		}
		if strings.ToUpper(input) == "Y" {
			return generator.OP_Finish
		} else {
			return generator.OP_NOOP
		}
	case generator.OP_Generate:
		return typ
	}

	return generator.OP_NOOP
}

func (a *App) List(typ string) (result []generator.OperationHolder) {
	switch typ {
	case C_Secret:
		for _, secret := range a.Secrets {
			result = append(result, secret)
		}
		return
	case C_ConfigMap:
		for _, cm := range a.ConfigMaps {
			result = append(result, cm)
		}
		return
	case C_Service:
		for _, svc := range a.Services {
			result = append(result, svc)
		}
		return
	case C_Component:
		for _, com := range a.Components {
			result = append(result, com)
		}
		return
	case C_Role:
		result = append(result, a.Role)
		return
	case C_RoleBinding:
		result = append(result, a.RoleBinding)
		return
	case C_ClusterRole:
		result = append(result, a.ClusterRole)
		return
	case C_ClusterRoleBinding:
		result = append(result, a.ClusterRoleBinding)
		return
	case C_Ingress:
		result = append(result, a.Ingress)
		return
	case C_ServiceAccount:
		result = append(result, a.ServiceAccount)
		return
	}
	return
}

func (a *App) Type() string {
	return C_APP
}

func (a *App) ArgsExample() string {
	return "name version namespace label1:values;label2:value2"
}

func (a *App) ParseArgs(input string) error {
	sections := strings.Split(input, " ")
	if len(sections) != len(a.Args()) {
		return fmt.Errorf("invalid args, required pattern: %s", a.ArgsExample())
	}

	if !generator.IsValidName(sections[0]) {
		return fmt.Errorf("invalid names, must match pattern: %s", generator.NamePatternString)
	}
	a.Name = sections[0]
	a.Labels[utils.AppLabel] = a.Name

	if sections[1] != "_" {
		a.Version = sections[1]
		a.Labels[utils.VersionLabel] = a.Version
	}

	if sections[2] != "_" {
		a.DefaultNamespace = sections[2]
	}

	if sections[3] != "_" {
		return convertAndMergeLabels(a.Labels, sections[3])
	}

	return nil
}

func (a *App) String() string {
	return fmt.Sprintf("%s: %s/%s", a.Type(), a.DefaultNamespace, a.Name)
}

func (a *App) Args() []string {
	return []string{A_Name, A_Version, A_Namespace, A_Labels}
}
