package templates

// type Templates struct {
// 	services.Service

// 	cfg Config

// 	lifecycler  *ring.BasicLifecycler
// 	ring        *ring.Ring
// 	directStore templatestore.TemplateStore
// 	cachedStore templatestore.TemplateStore

// 	// metrics *templatesMetrics

// 	registry prometheus.Registerer
// 	logger   log.Logger

// 	subservices        *services.Manager
// 	subservicesWatcher *services.FailureWatcher

// 	// // Pool of clients used to connect to other ruler replicas.
// 	// clientsPool ClientsPool

// 	// // Queue where we push rules syncing notifications to send to other ruler instances.
// 	// // This queue is also used to de-amplify the outbound notifications.
// 	// outboundSyncQueue          *templateSyncQueue
// 	// outboundSyncQueueProcessor *templateSyncQueueProcessor

// 	// // Queue where we pull rules syncing notifications received from other ruler instances.
// 	// // This queue is also used to de-amplify the inbound notifications.
// 	// inboundSyncQueue *templatesSyncQueue
// }

// func New(cfg Config, reg prometheus.Registerer, logger log.Logger, directStore, cachedStore templatestore.TemplateStore) (*Templates, error) {
// 	// If the cached store is not configured, just fallback to the direct one.
// 	if cachedStore == nil {
// 		cachedStore = directStore
// 	}

// 	templates := &Templates{
// 		cfg:         cfg,
// 		directStore: directStore,
// 		cachedStore: cachedStore,
// 		registry:    reg,
// 		logger:      logger,
// 		// clientsPool:       clientPool,
// 		// outboundSyncQueue: newRulerSyncQueue(cfg.syncQueuePollFrequency()),
// 		// inboundSyncQueue:  newRulerSyncQueue(cfg.syncQueuePollFrequency()),
// 		// metrics: newTemplateMetrics(reg),
// 	}

// 	ringStore, err := kv.NewClient(
// 		cfg.Ring.Common.KVStore,
// 		ring.GetCodec(),
// 		kv.RegistererWithKVName(prometheus.WrapRegistererWithPrefix("observability_operator_", reg), "template"),
// 		logger,
// 	)
// 	if err != nil {
// 		return nil, errors.Wrap(err, "create KV store client")
// 	}

// 	if err := enableSharding(templates, ringStore); err != nil {
// 		return nil, errors.Wrap(err, "setup ruler sharding ring")
// 	}

// 	templates.Service = services.NewBasicService(templates.starting, templates.run, templates.stopping)
// 	return templates, nil
// }

// func (t *Templates) SearchTemplates(name string) []*template.AppTemplate {
// 	return t.cachedStore.SearchTemplates(name)
// }

// func (t *Templates) GetTemplate(name string, version string) *template.AppTemplate {
// 	return t.cachedStore.GetTemplate(name, version)
// }

// type templatesMetrics struct {
// 	listRules       prometheus.Histogram
// 	loadRuleGroups  prometheus.Histogram
// 	ringCheckErrors prometheus.Counter
// 	rulerSync       *prometheus.CounterVec
// }

// func enableSharding(r *Templates, ringStore kv.Client) error {
// 	lifecyclerCfg, err := r.cfg.Ring.ToLifecyclerConfig(r.logger)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to initialize template's lifecycler config")
// 	}

// 	// Define lifecycler delegates in reverse order (last to be called defined first because they're
// 	// chained via "next delegate").
// 	delegate := ring.BasicLifecyclerDelegate(ring.NewInstanceRegisterDelegate(ring.JOINING, r.cfg.Ring.NumTokens))
// 	delegate = ring.NewLeaveOnStoppingDelegate(delegate, r.logger)
// 	delegate = ring.NewAutoForgetDelegate(r.cfg.Ring.Common.HeartbeatTimeout*ringAutoForgetUnhealthyPeriods, delegate, r.logger)

// 	ringName := "templates"
// 	r.lifecycler, err = ring.NewBasicLifecycler(lifecyclerCfg, ringName, RulerRingKey, ringStore, delegate, r.logger, prometheus.WrapRegistererWithPrefix("cortex_", r.registry))
// 	if err != nil {
// 		return errors.Wrap(err, "failed to initialize templates's lifecycler")
// 	}

// 	r.ring, err = ring.NewWithStoreClientAndStrategy(r.cfg.Ring.toRingConfig(), ringName, RulerRingKey, ringStore, ring.NewIgnoreUnhealthyInstancesReplicationStrategy(), prometheus.WrapRegistererWithPrefix("cortex_", r.registry), r.logger)
// 	if err != nil {
// 		return errors.Wrap(err, "failed to initialize templates's ring")
// 	}

// 	return nil
// }

// func (r *Templates) starting(ctx context.Context) error {
// 	var err error

// 	if r.subservices, err = services.NewManager(r.lifecycler, r.ring); err != nil {
// 		return errors.Wrap(err, "unable to start templates subservices")
// 	}

// 	r.subservicesWatcher = services.NewFailureWatcher()
// 	r.subservicesWatcher.WatchManager(r.subservices)

// 	if err = services.StartManagerAndAwaitHealthy(ctx, r.subservices); err != nil {
// 		return errors.Wrap(err, "unable to start templates subservices")
// 	}

// 	// Sync the rule when the ruler is JOINING the ring.
// 	// Activate the rule evaluation after the ruler is ACTIVE in the ring.
// 	// This is to make sure that the ruler is ready to evaluate rules immediately after it is ACTIVE in the ring.
// 	level.Info(r.logger).Log("msg", "waiting until templates is JOINING in the ring")
// 	if err := ring.WaitInstanceState(ctx, r.ring, r.lifecycler.GetInstanceID(), ring.JOINING); err != nil {
// 		return err
// 	}
// 	level.Info(r.logger).Log("msg", "ruler is JOINING in the ring")

// 	// Here during joining, we can download rules from object storage and sync them to the local rule manager
// 	r.syncRules(ctx, nil, rulerSyncReasonInitial, true)

// 	if err = r.lifecycler.ChangeState(ctx, ring.ACTIVE); err != nil {
// 		return errors.Wrapf(err, "switch instance to %s in the ring", ring.ACTIVE)
// 	}

// 	level.Info(r.logger).Log("msg", "waiting until ruler is ACTIVE in the ring")
// 	if err := ring.WaitInstanceState(ctx, r.ring, r.lifecycler.GetInstanceID(), ring.ACTIVE); err != nil {
// 		return err
// 	}
// 	level.Info(r.logger).Log("msg", "ruler is ACTIVE in the ring")

// 	r.manager.Start()
// 	level.Info(r.logger).Log("msg", "ruler is only now starting to evaluate rules")

// 	// TODO: ideally, ruler would wait until its queryable is finished starting.
// 	return nil
// }

// // Stop stops the Ruler.
// // Each function of the ruler is terminated before leaving the ring
// func (r *Templates) stopping(_ error) error {
// 	r.manager.Stop()

// 	if r.subservices != nil {
// 		_ = services.StopManagerAndAwaitStopped(context.Background(), r.subservices)
// 	}
// 	return nil
// }
