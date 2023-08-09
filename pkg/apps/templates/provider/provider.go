package provider

import (
	"github.com/udmire/observability-operator/pkg/apps/templates/template"
)

type TemplateProvider interface {
	SearchTemplates(name string) []*template.AppTemplate
	GetTemplate(name, version string) *template.AppTemplate
	GetLatestTemplate(name string) *template.AppTemplate
}
