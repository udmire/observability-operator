package manager

import (
	"context"
	"flag"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"

	udmirecnv1alpha1 "github.com/udmire/observability-operator/api/v1alpha1"
	//+kubebuilder:scaffold:imports
)

type CtrlManagerWraper interface {
	Manager() ctrl.Manager
}

type Config struct {
	Port                  int    `yaml:"port"`
	MetricsAddress        string `yaml:"metric_address"`
	ProbeAddress          string `yaml:"probe_address"`
	EnabledLeaderElection bool   `yaml:"enable_leader_election"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Port, "manager.port", 9443, "The operator endpoint binds to.")
	f.StringVar(&c.MetricsAddress, "manager.metrics-address", ":8080", "The address the metric endpoint binds to.")
	f.StringVar(&c.ProbeAddress, "manager.probe-address", ":8081", "The address the probe endpoint binds to.")
	f.BoolVar(&c.EnabledLeaderElection, "manager.leader-elect", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
}

type managerWraper struct {
	*services.BasicService

	cfg    Config
	logger log.Logger

	Mgr manager.Manager
}

func NewManagerWraper(cfg Config, logger log.Logger) *managerWraper {
	wraper := &managerWraper{
		cfg:    cfg,
		logger: logger,
	}

	_ = wraper.createManager()

	wraper.BasicService = services.NewBasicService(nil, wraper.running, nil)

	return wraper
}

func (w *managerWraper) createManager() error {

	scheme := runtime.NewScheme()

	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(udmirecnv1alpha1.AddToScheme(scheme))

	ctrl.SetLogger(zap.New(zap.UseFlagOptions(&zap.Options{
		Development: true,
	})))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     w.cfg.MetricsAddress,
		Port:                   w.cfg.Port,
		HealthProbeBindAddress: w.cfg.ProbeAddress,
		LeaderElection:         w.cfg.EnabledLeaderElection,
		LeaderElectionID:       "6618dc8d.udmire.cn",
	})
	if err != nil {
		level.Error(w.logger).Log("msg", "unable to start controller manager", "err", err)
		return err
	}

	w.Mgr = mgr
	return nil
}

func (w *managerWraper) running(serviceContext context.Context) error {
	if err := w.Mgr.AddHealthzCheck("healthz", healthz.Ping); err != nil {
		level.Error(w.logger).Log("msg", "unable to set up health check", "err", err)
		return err
	}
	if err := w.Mgr.AddReadyzCheck("readyz", healthz.Ping); err != nil {
		level.Error(w.logger).Log("msg", "unable to set up ready check", "err", err)
		return err
	}

	level.Info(w.logger).Log("msg", "starting manager")
	if err := w.Mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		level.Error(w.logger).Log("msg", "problem running manager", "err", err)
		return err
	}
	return nil
}

func (w *managerWraper) Manager() ctrl.Manager {
	return w.Mgr
}
