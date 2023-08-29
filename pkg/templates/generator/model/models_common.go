package model

import "fmt"

type Common struct {
	ConfigMaps []*ConfigMap `yaml:"configmaps,omitempty"`
	Secrets    []*Secret    `yaml:"secrets,omitempty"`
	Services   []*Service   `yaml:"services,omitempty"`

	ServiceAccount     *ServiceAccount     `yaml:"serviceaccount,omitempty"`
	ClusterRole        *ClusterRole        `yaml:"clusterrole,omitempty"`
	ClusterRoleBinding *ClusterRoleBinding `yaml:"clusterrolebinding,omitempty"`
	Role               *Role               `yaml:"role,omitempty"`
	RoleBinding        *RoleBinding        `yaml:"rolebinding,omitempty"`
	Ingress            *Ingress            `yaml:"ingress,omitempty"`
}

type ConfigMap struct {
	GenericModel `yaml:",inline"`
}

func (c *ConfigMap) Type() string {
	return C_ConfigMap
}

func (c *ConfigMap) String() string {
	return fmt.Sprintf("%s: %s/%s", c.Type(), c.DefaultNamespace, c.Name)
}

type Secret struct {
	GenericModel `yaml:",inline"`
}

func (s *Secret) Type() string {
	return C_Secret
}
func (s *Secret) String() string {
	return fmt.Sprintf("%s: %s/%s", s.Type(), s.DefaultNamespace, s.Name)
}

type Service struct {
	GenericModel `yaml:",inline"`
}

func (s *Service) Type() string {
	return C_Service
}
func (s *Service) String() string {
	return fmt.Sprintf("%s: %s/%s", s.Type(), s.DefaultNamespace, s.Name)
}

type ServiceAccount struct {
	GenericModel `yaml:",inline"`
}

func (a *ServiceAccount) Type() string {
	return C_ServiceAccount
}
func (a *ServiceAccount) String() string {
	return fmt.Sprintf("%s: %s/%s", a.Type(), a.DefaultNamespace, a.Name)
}

type ClusterRole struct {
	GenericModel `yaml:",inline"`
}

func (r *ClusterRole) Type() string {
	return C_ClusterRole
}
func (r *ClusterRole) String() string {
	return fmt.Sprintf("%s: %s/%s", r.Type(), r.DefaultNamespace, r.Name)
}

type ClusterRoleBinding struct {
	GenericModel `yaml:",inline"`
}

func (b *ClusterRoleBinding) Type() string {
	return C_ClusterRoleBinding
}
func (b *ClusterRoleBinding) String() string {
	return fmt.Sprintf("%s: %s/%s", b.Type(), b.DefaultNamespace, b.Name)
}

type Role struct {
	GenericModel `yaml:",inline"`
}

func (r *Role) Type() string {
	return C_Role
}
func (r *Role) String() string {
	return fmt.Sprintf("%s: %s/%s", r.Type(), r.DefaultNamespace, r.Name)
}

type RoleBinding struct {
	GenericModel `yaml:",inline"`
}

func (b *RoleBinding) Type() string {
	return C_RoleBinding
}

func (b *RoleBinding) String() string {
	return fmt.Sprintf("%s: %s/%s", b.Type(), b.DefaultNamespace, b.Name)
}

type Ingress struct {
	GenericModel `yaml:",inline"`
}

func (i *Ingress) Type() string {
	return C_Ingress
}

func (i *Ingress) String() string {
	return fmt.Sprintf("%s: %s/%s", i.Type(), i.DefaultNamespace, i.Name)
}
