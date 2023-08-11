package operator

import (
	"github.com/go-kit/log/level"
	"github.com/grafana/dskit/modules"
	"github.com/grafana/dskit/services"
	"github.com/udmire/observability-operator/pkg/apps/templates/store/local"
	"github.com/udmire/observability-operator/pkg/operator/agents"
	"github.com/udmire/observability-operator/pkg/operator/apps"
	"github.com/udmire/observability-operator/pkg/operator/exporters"
	"github.com/udmire/observability-operator/pkg/operator/manager"
	util_log "github.com/udmire/observability-operator/pkg/utils/log"
)

// The various modules that make up Mimir.
const (
	TemplateStorage string = "template-storage"
	CtrlManager     string = "ctrl-manager"
	Apps            string = "apps"
	Agents          string = "agents"
	Exporters       string = "exporters"
	All             string = "all"
)

func (op *Operator) initTemplateStore() (serv services.Service, err error) {
	if !op.Cfg.isAnyModuleEnabled(Apps, Agents, Exporters, All) {
		level.Info(util_log.Logger).Log("msg", "The templatestore is not being started because you need to configure the template storage.")
		return
	}

	store := local.New(op.Cfg.TemplateStore.Local, util_log.Logger)
	op.TemplateStore = store
	return store, nil
}

func (op *Operator) initCtrlManager() (serv services.Service, err error) {
	wrapper := manager.NewManagerWraper(op.Cfg.Manager, util_log.Logger)
	op.ControllerManager = wrapper
	return wrapper, nil
}

func (op *Operator) initAgentsController() (serv services.Service, err error) {
	ctrl := agents.New(op.ControllerManager.Manager().GetClient(), op.ControllerManager.Manager().GetScheme(), op.TemplateStore, util_log.Logger)

	ctrl.SetManager(op.ControllerManager.Manager())

	return ctrl, nil
}

func (op *Operator) initAppsController() (serv services.Service, err error) {
	ctrl := apps.New(op.ControllerManager.Manager().GetClient(), op.ControllerManager.Manager().GetScheme(), op.Cfg.Apps, op.TemplateStore, util_log.Logger)
	ctrl.SetManager(op.ControllerManager.Manager())

	return ctrl, nil
}

func (op *Operator) initExportersController() (serv services.Service, err error) {
	ctrl := exporters.New(op.ControllerManager.Manager().GetClient(), op.ControllerManager.Manager().GetScheme(), op.Cfg.Exporters, op.TemplateStore, util_log.Logger)
	ctrl.SetManager(op.ControllerManager.Manager())

	return ctrl, nil
}

func (op *Operator) SetupModuleManager() error {
	mm := modules.NewManager(util_log.Logger)

	mm.RegisterModule(TemplateStorage, op.initTemplateStore, modules.UserInvisibleModule)
	mm.RegisterModule(CtrlManager, op.initCtrlManager, modules.UserInvisibleModule)
	mm.RegisterModule(Agents, op.initAgentsController)
	mm.RegisterModule(Apps, op.initAppsController)
	mm.RegisterModule(Exporters, op.initExportersController)
	mm.RegisterModule(All, nil)

	// Add dependencies
	deps := map[string][]string{
		Apps:      {TemplateStorage, CtrlManager},
		Agents:    {TemplateStorage, CtrlManager},
		Exporters: {TemplateStorage, CtrlManager},
		All:       {Apps, Agents, Exporters},
	}

	for mod, targets := range deps {
		if err := mm.AddDependency(mod, targets...); err != nil {
			return err
		}
	}

	op.ModuleManager = mm

	return nil
}
