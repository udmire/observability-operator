package template

type TemplateBase struct {
	Name          string
	Version       string
	TemplateFiles []*TemplateFile
}

type TemplateFile struct {
	FileName string
	Content  []byte
}

type AppTemplate struct {
	TemplateBase
	Workloads map[string]*WorkloadTemplate
}

type WorkloadTemplate struct {
	TemplateBase
}
