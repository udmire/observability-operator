package model

import (
	"bufio"
	"fmt"

	"strings"

	"github.com/udmire/observability-operator/pkg/templates/generator"
)

type GenericModel struct {
	Base             `yaml:",inline"`
	DefaultNamespace string `yaml:"namespace,omitempty"`
}

func (g *GenericModel) NextOperation(scanner *bufio.Scanner) []generator.Operation {
	return []generator.Operation{}
}

// 列出可选类型，并选择
func (g *GenericModel) Create(scanner *bufio.Scanner) generator.OperationHolder {
	panic("not implemented") // TODO: Implement
}

func (g *GenericModel) Modify(scanner *bufio.Scanner) generator.OperationHolder {
	panic("not implemented") // TODO: Implement
}

func (g *GenericModel) Remove(scanner *bufio.Scanner) {
	panic("not implemented") // TODO: Implement
}

func (g *GenericModel) General(scanner *bufio.Scanner) generator.GeneralCommand {
	panic("not implemented") // TODO: Implement
}

func (g *GenericModel) Type() string {
	panic("not implemented")
}

func (g *GenericModel) Args() []string {
	return []string{A_Name, A_Namespace, A_Labels}
}

func (g *GenericModel) ArgsExample() string {
	return "name namespace label1:values;label2:value2"
}

func (g *GenericModel) ParseArgs(input string) error {
	sections := strings.Split(input, " ")
	if len(sections) != len(g.Args()) {
		return fmt.Errorf("invalid args, required pattern: %s", g.ArgsExample())
	}

	if sections[0] != "_" {
		if !generator.IsValidName(sections[0]) {
			return fmt.Errorf("invalid names, must match pattern: %s", generator.NamePatternString)
		}
		g.Name = sections[0]
	}

	if sections[1] != "_" {
		g.DefaultNamespace = sections[1]
	}

	if sections[2] != "_" {
		return convertAndMergeLabels(g.Labels, sections[2])
	}

	return nil
}

func (g *GenericModel) String() string {
	return fmt.Sprintf("%s: %s/%s", g.Type(), g.DefaultNamespace, g.Name)
}

func createCommonOperations(parent GenericModel, comm *Common, typ string) generator.OperationHolder {
	switch typ {
	case C_ConfigMap:
		cm := &ConfigMap{GenericModel: inheritCreate(parent)}
		comm.ConfigMaps = append(comm.ConfigMaps, cm)
		return cm
	case C_Secret:
		sec := &Secret{GenericModel: inheritCreate(parent)}
		comm.Secrets = append(comm.Secrets, sec)
		return sec
	case C_Service:
		svc := &Service{GenericModel: inheritCreate(parent)}
		comm.Services = append(comm.Services, svc)
		return svc
	case C_ServiceAccount:
		sa := &ServiceAccount{GenericModel: inheritCreate(parent)}
		comm.ServiceAccount = sa
		return sa
	case C_ClusterRole:
		role := &ClusterRole{GenericModel: inheritCreateWithoutNamespace(parent)}
		comm.ClusterRole = role
		return role
	case C_ClusterRoleBinding:
		bind := &ClusterRoleBinding{GenericModel: inheritCreateWithoutNamespace(parent)}
		comm.ClusterRoleBinding = bind
		return bind
	case C_Role:
		role := &Role{GenericModel: inheritCreate(parent)}
		comm.Role = role
		return role
	case C_RoleBinding:
		bind := &RoleBinding{GenericModel: inheritCreate(parent)}
		comm.RoleBinding = bind
		return bind
	case C_Ingress:
		ing := &Ingress{GenericModel: inheritCreate(parent)}
		comm.Ingress = ing
		return ing
	}
	return nil
}

func inheritCreate(parent GenericModel) GenericModel {
	result := inheritCreateWithoutNamespace(parent)
	result.DefaultNamespace = parent.DefaultNamespace
	return result
}

func inheritCreateWithoutNamespace(parent GenericModel) GenericModel {
	result := GenericModel{Base: Base{Labels: make(map[string]string)}}
	result.Name = parent.Name
	mergeLabels(result.Labels, parent.Labels)
	return result
}

func convertAndMergeLabels(ori map[string]string, content string) (err error) {
	if len(content) < 1 {
		return
	}

	labels := strings.Split(content, ";")
	for _, label := range labels {
		if len(label) == 0 {
			continue
		}

		sections := strings.Split(label, ":")
		if len(sections) != 2 || len(sections[0]) == 0 {
			return fmt.Errorf("invalid labels format, should be \"a:b;c:d\"")
		}

		if len(sections[1]) == 0 {
			delete(ori, sections[0])
		} else {
			ori[sections[0]] = sections[1]
		}
	}
	return nil
}
