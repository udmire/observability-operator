package operator

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/flagext"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/services"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
	"github.com/udmire/observability-operator/pkg/configs/logging"
	"github.com/udmire/observability-operator/pkg/operator/agents"
	"github.com/udmire/observability-operator/pkg/operator/apps"
	"github.com/udmire/observability-operator/pkg/operator/exporters"
	"github.com/udmire/observability-operator/pkg/operator/manager"
	info "github.com/udmire/observability-operator/pkg/operator/providers"
	"github.com/udmire/observability-operator/pkg/templates/provider"
	"github.com/udmire/observability-operator/pkg/templates/store/category"
	"github.com/udmire/observability-operator/pkg/utils"
	util_log "github.com/udmire/observability-operator/pkg/utils/log"
	"github.com/udmire/observability-operator/pkg/utils/process"
	"github.com/udmire/observability-operator/pkg/utils/signals"
	"go.uber.org/atomic"
	"gopkg.in/yaml.v2"
)

type Config struct {
	Target                 flagext.StringSliceCSV `yaml:"target"`
	EnableGoRuntimeMetrics bool                   `yaml:"enable_go_runtime_metrics" category:"advanced"`
	ShutdownDelay          time.Duration          `yaml:"shutdown_delay" category:"experimental"`
	PrintConfig            bool                   `yaml:"-"`
	ApplicationName        string                 `yaml:"-"`

	Logging logging.Config `yaml:"logging"`
	Manager manager.Config `yaml:"manager"`

	Apps          apps.Config      `yaml:"apps"`
	Exporters     exporters.Config `yaml:"exporters"`
	TemplateStore category.Config  `yaml:"template_store"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	c.ApplicationName = "Observperator"
	c.Target = []string{All}

	f.Var(&c.Target, "target", "Comma-separated list of components to include in the instantiated process. "+
		"The default value 'all' includes all components that are required to form a functional Observability Operator instance in single-binary mode. "+
		"Use the '-modules' command line flag to get a list of available components, and to see which components are included with 'all'.")
	f.BoolVar(&c.PrintConfig, "print.config", false, "Print the config and exit.")
	f.DurationVar(&c.ShutdownDelay, "shutdown-delay", 0, "How long to wait between SIGTERM and shutdown. After receiving SIGTERM, Operator will report not-ready status via /ready endpoint.")

	c.Logging.RegisterFlags(f)
	c.Manager.RegisterFlags(f)

	c.Apps.RegisterFlags(f)
	c.Exporters.RegisterFlags(f)
	c.TemplateStore.RegisterFlags(f)
}

func (c *Config) isAnyModuleEnabled(modules ...string) bool {
	for _, m := range modules {
		if utils.StringsContain(c.Target, m) {
			return true
		}
	}

	return false
}

type Operator struct {
	Cfg        Config
	Registerer prometheus.Registerer

	// set during initialization
	ServiceMap    map[string]services.Service
	ModuleManager *modules.Manager

	TemplateStore     provider.CategryTemplateProvider
	ControllerManager manager.CtrlManagerWraper
	InfoProviders     info.Providers

	AppsController      *apps.AppsReconciler
	AgentsController    *agents.AgentsReconciler
	ExportersController *exporters.ExportersReconciler
}

// New makes a new Mimir.
func New(cfg Config, reg prometheus.Registerer) (*Operator, error) {
	if cfg.PrintConfig {
		if err := yaml.NewEncoder(os.Stdout).Encode(&cfg); err != nil {
			fmt.Println("Error encoding config:", err)
		}
		os.Exit(0)
	}

	if cfg.EnableGoRuntimeMetrics {
		// unregister default Go collector
		reg.Unregister(collectors.NewGoCollector())
		// register Go collector with all available runtime metrics
		reg.MustRegister(collectors.NewGoCollector(
			collectors.WithGoCollectorRuntimeMetrics(collectors.MetricsAll),
		))
	}

	op := &Operator{
		Cfg:        cfg,
		Registerer: reg,
	}

	if err := op.SetupModuleManager(); err != nil {
		return nil, err
	}

	return op, nil
}

func (op *Operator) Run() error {
	// Register custom process metrics.
	if c, err := process.NewProcessCollector(); err == nil {
		if op.Registerer != nil {
			op.Registerer.MustRegister(c)
		}
	} else {
		level.Warn(util_log.Logger).Log("msg", "skipped registration of custom process metrics collector", "err", err)
	}

	var err error
	op.ServiceMap, err = op.ModuleManager.InitModuleServices(op.Cfg.Target...)
	if err != nil {
		return err
	}

	// get all services, create service manager and tell it to start
	servs := []services.Service(nil)
	for _, s := range op.ServiceMap {
		servs = append(servs, s)
	}

	sm, err := services.NewManager(servs...)
	if err != nil {
		return err
	}

	// Used to delay shutdown but return "not ready" during this delay.
	shutdownRequested := atomic.NewBool(false)

	// Let's listen for events from this manager, and log them.
	healthy := func() { level.Info(util_log.Logger).Log("msg", "Application started") }
	stopped := func() { level.Info(util_log.Logger).Log("msg", "Application stopped") }
	serviceFailed := func(service services.Service) {
		// if any service fails, stop entire Mimir
		sm.StopAsync()

		// let's find out which module failed
		for m, s := range op.ServiceMap {
			if s == service {
				if errors.Is(service.FailureCase(), modules.ErrStopProcess) {
					level.Info(util_log.Logger).Log("msg", "received stop signal via return error", "module", m, "err", service.FailureCase())
				} else {
					level.Error(util_log.Logger).Log("msg", "module failed", "module", m, "err", service.FailureCase())
				}
				return
			}
		}

		level.Error(util_log.Logger).Log("msg", "module failed", "module", "unknown", "err", service.FailureCase())
	}

	sm.AddListener(services.NewManagerListener(healthy, stopped, serviceFailed))

	// // Setup signal handler to gracefully shutdown in response to SIGTERM or SIGINT
	handler := signals.NewHandler(util_log.Logger)
	go func() {
		handler.Loop()

		shutdownRequested.Store(true)

		if op.Cfg.ShutdownDelay > 0 {
			time.Sleep(op.Cfg.ShutdownDelay)
		}

		sm.StopAsync()
	}()

	// Start all services. This can really only fail if some service is already
	// in other state than New, which should not be the case.
	err = sm.StartAsync(context.Background())
	if err == nil {
		// Wait until service manager stops. It can stop in two ways:
		// 1) Signal is received and manager is stopped.
		// 2) Any service fails.
		err = sm.AwaitStopped(context.Background())
	}

	// If there is no error yet (= service manager started and then stopped without problems),
	// but any service failed, report that failure as an error to caller.
	if err == nil {
		if failed := sm.ServicesByState()[services.Failed]; len(failed) > 0 {
			for _, f := range failed {
				if !errors.Is(f.FailureCase(), modules.ErrStopProcess) {
					// Details were reported via failure listener before
					err = errors.New("failed services")
					break
				}
			}
		}
	}
	return err
}
