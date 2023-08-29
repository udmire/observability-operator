package provider

import (
	"github.com/grafana/dskit/services"
	"github.com/udmire/observability-operator/pkg/templates/template"
)

const (
	Apps     string = "apps"
	Capsules string = "capsules"
)

type TemplateProvider interface {
	services.Service

	SearchTemplates(name string) []*template.AppTemplate
	GetTemplate(name, version string) *template.AppTemplate
	GetLatestTemplate(name string) *template.AppTemplate
}

type CategryTemplateProvider interface {
	GetProvider(category string) TemplateProvider
}
