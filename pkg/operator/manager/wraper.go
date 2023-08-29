package manager

import (
	"context"
	"flag"

	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/services"
	"github.com/pkg/errors"
	udmirecnv1alpha1 "github.com/udmire/observability-operator/api/v1alpha1"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type ClusterNameProvider func() string

type CtrlManagerWraper interface {
	Manager() ctrl.Manager
	ClusterNameProvider() ClusterNameProvider
}

type Config struct {
	Port                  int    `yaml:"port"`
	MetricsAddress        string `yaml:"metric_address"`
	ProbeAddress          string `yaml:"probe_address"`
	EnabledLeaderElection bool   `yaml:"enable_leader_election"`
	ClusterName           string `yaml:"cluster_name"`
}

func (c *Config) RegisterFlags(f *flag.FlagSet) {
	f.IntVar(&c.Port, "manager.port", 9443, "The operator endpoint binds to.")
	f.StringVar(&c.MetricsAddress, "manager.metrics-address", ":8080", "The address the metric endpoint binds to.")
	f.StringVar(&c.ProbeAddress, "manager.probe-address", ":8081", "The address the probe endpoint binds to.")
	f.BoolVar(&c.EnabledLeaderElection, "manager.leader-elect", false, "Enable leader election for controller manager. Enabling this will ensure there is only one active controller manager.")
	f.StringVar(&c.ClusterName, "manager.cluster-name", "", "The k8s cluster name to specified.")
}

type managerWraper struct {
	*services.BasicService

	cfg    Config
	logger log.Logger

	clusterInfo *clusterInfoProvider

	subservices        *services.Manager
	subservicesWatcher *services.FailureWatcher

	Mgr manager.Manager
}

func NewManagerWraper(cfg Config, logger log.Logger) *managerWraper {
	wraper := &managerWraper{
		cfg:    cfg,
		logger: logger,
	}

	_ = wraper.createManager()
	wraper.clusterInfo = newClusterInfoProvider(&wraper.cfg, wraper.Mgr.GetClient(), logger)

	wraper.BasicService = services.NewBasicService(wraper.starting, wraper.run, wraper.stopping)

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

func (w *managerWraper) Manager() ctrl.Manager {
	return w.Mgr
}

func (w *managerWraper) ClusterNameProvider() ClusterNameProvider {
	return func() string {
		return w.clusterInfo.cfg.ClusterName
	}
}

func (r *managerWraper) starting(ctx context.Context) error {
	var err error

	if r.subservices, err = services.NewManager(services.NewBasicService(nil, r.startManager, nil), r.clusterInfo); err != nil {
		return errors.Wrap(err, "unable to start category stores")
	}

	r.subservicesWatcher = services.NewFailureWatcher()
	r.subservicesWatcher.WatchManager(r.subservices)

	if err = services.StartManagerAndAwaitHealthy(ctx, r.subservices); err != nil {
		return errors.Wrap(err, "unable to start controller manager wraper")
	}

	return nil
}

func (r *managerWraper) run(ctx context.Context) error {
	level.Info(r.logger).Log("msg", "controller manager up and running")

	for {
		select {
		case <-ctx.Done():
			return nil
		case err := <-r.subservicesWatcher.Chan():
			return errors.Wrap(err, "controller manager subservice failed")
		}
	}
}

func (r *managerWraper) stopping(_ error) error {
	if r.subservices != nil {
		_ = services.StopManagerAndAwaitStopped(context.Background(), r.subservices)
	}
	return nil
}
