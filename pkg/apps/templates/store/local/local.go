package local

import (
	"fmt"
	"io/fs"
	"os"
	"strings"
	"sync"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/udmire/observability-operator/pkg/apps/templates/template"
)

const (
	Name = "local"
)

type Config struct {
	Directory string `yaml:"directory"`
}

type LocalStore struct {
	cfg Config

	logger    log.Logger
	lock      sync.Mutex
	loader    template.TemplateLoader
	templates map[string]*template.AppTemplate
}

func New(cfg Config, logger log.Logger) *LocalStore {
	return &LocalStore{
		cfg:       cfg,
		lock:      sync.Mutex{},
		logger:    logger,
		loader:    template.NewTemplateLoader(logger),
		templates: make(map[string]*template.AppTemplate),
	}
}

func (l *LocalStore) Load() error {
	err := fs.WalkDir(os.DirFS(l.cfg.Directory), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			return nil // ignore all the content in sub folder.
		}

		return l.LoadTemplate(path)
	})

	if err != nil {
		level.Warn(l.logger).Log("msg", "failed to load templates", "err", err)
		return err
	}
	return nil
}

func (l *LocalStore) LoadTemplate(path string) error {
	appVer, _ := l.loader.TemplateName(path)
	temp, err := l.loader.LoadTemplate(path)
	if err != nil {
		return err
	}
	l.lock.Lock()
	defer l.lock.Unlock()
	l.templates[appVer] = temp
	return nil
}

func (l *LocalStore) UnloadTemplate(path string) {
	appVer, _ := l.loader.TemplateName(path)

	l.lock.Lock()
	defer l.lock.Unlock()
	delete(l.templates, appVer)
}

func (l *LocalStore) SyncTemplates() {
	l.lock.Lock()
	defer l.lock.Unlock()
	l.Load()
}

func (l *LocalStore) ListAppTemplates() {

}

func (l *LocalStore) SearchTemplates(name string) []*template.AppTemplate {
	temps := l.templates
	var result []*template.AppTemplate
	for k, at := range temps {
		if strings.Contains(k, name) {
			result = append(result, at)
		}
	}
	return result
}

func (l *LocalStore) GetTemplate(name, version string) *template.AppTemplate {
	temps := l.templates
	appVer := fmt.Sprintf("%s_%s", name, version)
	if app, ok := temps[appVer]; ok {
		return app
	}

	return nil
}
