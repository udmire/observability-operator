package manager

import (
	"context"

	"github.com/go-kit/log/level"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/healthz"
	//+kubebuilder:scaffold:imports
)

func (w *managerWraper) startManager(_ context.Context) error {
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
