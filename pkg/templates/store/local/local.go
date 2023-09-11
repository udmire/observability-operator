package local

import (
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/udmire/observability-operator/pkg/templates/template"
	"github.com/udmire/observability-operator/pkg/utils"
)

const (
	Name = "local"
)

type Config struct {
	Directory string `yaml:"directory"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.Directory, "templates.store.directory", "/data/templates", "where the templates stored in local.")
}

type LocalStore struct {
	*services.BasicService

	watcher *fsnotify.Watcher

	cfg Config

	logger    log.Logger
	lock      sync.Mutex
	loader    template.TemplateLoader
	templates map[string]*template.AppTemplate
}

func New(cfg Config, logger log.Logger) *LocalStore {
	err := os.MkdirAll(cfg.Directory, os.ModeDir)
	if err != nil {
		level.Warn(logger).Log("msg", "failed to validate the template folder", "folder", cfg.Directory, "err", err)
		panic("invalid config for templates")
	}

	store := &LocalStore{
		cfg:       cfg,
		lock:      sync.Mutex{},
		logger:    logger,
		loader:    template.NewTemplateLoader(logger),
		templates: make(map[string]*template.AppTemplate),
	}
	store.BasicService = services.NewBasicService(store.startup, store.watching, store.shutdown)
	return store
}

func (l *LocalStore) Load() error {
	err := fs.WalkDir(os.DirFS(l.cfg.Directory), ".", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if entry.IsDir() {
			if path == "." {
				return nil
			}
			return fs.SkipDir // ignore all the content in sub folder.
		}

		return l.LoadTemplate(filepath.Join(l.cfg.Directory, path))
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

	if temp == nil {
		return nil
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

func (l *LocalStore) GetLatestTemplate(name string) *template.AppTemplate {
	temps := l.templates
	appPrefix := fmt.Sprintf("%s_", name)

	var result *template.AppTemplate
	var latest string

	for k, at := range temps {
		if !strings.HasPrefix(k, appPrefix) {
			continue
		}
		version := strings.TrimLeft(k, appPrefix)
		if len(latest) <= 0 {
			latest = version
			result = at
		} else if utils.IsNewerThan(version, latest) {
			latest = version
			result = at
		}
	}

	return result
}
