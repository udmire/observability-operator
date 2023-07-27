package templates

type TemplateProvider interface {
	GetTemplates(name string) []*AppTemplate
	GetTemplate(name, version string) *AppTemplate
}

type RemoteTemplateProvider struct{}

type LocalTemplateProvider struct{}
