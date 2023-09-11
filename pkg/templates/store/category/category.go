package category

import (
	"context"
	"flag"
	"path/filepath"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	"github.com/udmire/observability-operator/pkg/templates/store/local"
	"github.com/udmire/observability-operator/pkg/templates/store/sync"
)

type Config struct {
	BaseDirectory string                 `yaml:"base_directory"`
	Categories    flagext.StringSliceCSV `yaml:"categories"`
	Synchronize   sync.Config            `yaml:"sync"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.StringVar(&c.BaseDirectory, "templates.store.category.base-directory", "/data/templates", "where the templates stored in local.")
	c.Categories = []string{provider.Apps, provider.Capsules}
	f.Var(&c.Categories, "templates.store.category.categories", "Comma-separated list of template categories.")

	c.Synchronize.RegisterFlags(f)
}

type CategoryStore struct {
	*services.BasicService

	cfg    Config
	logger log.Logger

	subservices        *services.Manager
	subservicesWatcher *services.FailureWatcher

	providers   map[string]provider.TemplateProvider
	synchorizer provider.TemplatesSynchronizer
}

func New(cfg Config, logger log.Logger) *CategoryStore {
	store := &CategoryStore{
		cfg:    cfg,
		logger: logger,
	}
	store.providers = buildProviders(cfg, logger)
	if store.cfg.Synchronize.Enabled {
		store.synchorizer = sync.NewHttpSynchronizer(cfg.Synchronize, cfg.BaseDirectory, logger)
	}
	store.BasicService = services.NewBasicService(store.starting, store.run, store.stopping)
	return store
}

func buildProviders(cfg Config, logger log.Logger) map[string]provider.TemplateProvider {
	providers := make(map[string]provider.TemplateProvider)

	for _, typ := range cfg.Categories {
		lc := local.Config{
			Directory: filepath.Join(cfg.BaseDirectory, typ),
		}
		provider := local.New(lc, logger)
		providers[typ] = provider
	}
	return providers
}

func (r *CategoryStore) starting(ctx context.Context) error {
	var err error

	var svcs []services.Service
	for _, provider := range r.providers {
		svcs = append(svcs, provider)
	}

	if r.synchorizer != nil {
		svcs = append(svcs, r.synchorizer)
	}

	if r.subservices, err = services.NewManager(svcs...); err != nil {
		return errors.Wrap(err, "unable to start category stores")
	}

	r.subservicesWatcher = services.NewFailureWatcher()
	r.subservicesWatcher.WatchManager(r.subservices)

	if err = services.StartManagerAndAwaitHealthy(ctx, r.subservices); err != nil {
		return errors.Wrap(err, "unable to start category stores")
	}

	return nil
}

func (r *CategoryStore) run(ctx context.Context) error {
	level.Info(r.logger).Log("msg", "category store up and running")

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-r.subservicesWatcher.Chan():
			return errors.Wrap(err, "category store subservice failed")
		}
	}
}

func (r *CategoryStore) stopping(_ error) error {
	if r.subservices != nil {
		_ = services.StopManagerAndAwaitStopped(context.Background(), r.subservices)
	}
	return nil
}

func (r *CategoryStore) GetProvider(category string) provider.TemplateProvider {
	return r.providers[category]
}
