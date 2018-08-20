package radar

import (
	"context"
	"errors"
	"reflect"
	"time"

	"code.cloudfoundry.org/clock"
	"code.cloudfoundry.org/lager"
	"github.com/concourse/atc"
	"github.com/concourse/atc/creds"
	"github.com/concourse/atc/db"
	"github.com/concourse/atc/metric"
	"github.com/concourse/atc/resource"
	"github.com/concourse/atc/worker"
)

//go:generate counterfeiter . Scannable

type Scannable interface {
	Name() string
	CheckEvery() string
	Type() string
	Source() atc.Source
	SetResourceConfig(int) error
	Paused() bool
	Tags() atc.Tags
	SetCheckError(error) error
}

type resourceConfigScanner struct {
	clock                             clock.Clock
	resourceFactory                   resource.ResourceFactory
	resourceConfigCheckSessionFactory db.ResourceConfigCheckSessionFactory
	resourceConfigFactory             db.ResourceConfigFactory
	defaultInterval                   time.Duration
	dbPipeline                        db.Pipeline
	externalURL                       string
	variables                         creds.Variables
}

func NewResourceConfigScanner(
	clock clock.Clock,
	resourceFactory resource.ResourceFactory,
	resourceConfigCheckSessionFactory db.ResourceConfigCheckSessionFactory,
	resourceConfigFactory db.ResourceConfigFactory,
	defaultInterval time.Duration,
	dbPipeline db.Pipeline,
	externalURL string,
	variables creds.Variables,
) Scanner {
	return &resourceConfigScanner{
		clock:                             clock,
		resourceFactory:                   resourceFactory,
		resourceConfigCheckSessionFactory: resourceConfigCheckSessionFactory,
		resourceConfigFactory:             resourceConfigFactory,
		defaultInterval:                   defaultInterval,
		dbPipeline:                        dbPipeline,
		externalURL:                       externalURL,
		variables:                         variables,
	}
}

var ErrFailedToAcquireLock = errors.New("failed-to-acquire-lock")

func (scanner *resourceConfigScanner) Run(logger lager.Logger, scannable Scannable) (time.Duration, error) {
	interval, err := scanner.scan(logger.Session("tick"), scannable, nil, false)

	err = swallowErrResourceScriptFailed(err)

	return interval, err
}

func (scanner *resourceConfigScanner) ScanFromVersion(logger lager.Logger, scannable Scannable, fromVersion atc.Version) error {
	_, err := scanner.scan(logger, scannable, fromVersion, true)

	return err
}

func (scanner *resourceConfigScanner) Scan(logger lager.Logger, scannable Scannable) error {
	_, err := scanner.scan(logger, scannable, nil, true)

	err = swallowErrResourceScriptFailed(err)

	return err
}

func (scanner *resourceConfigScanner) scan(logger lager.Logger, scannable Scannable, fromVersion atc.Version, mustComplete bool) (time.Duration, error) {
	lockLogger := logger.Session("lock", lager.Data{
		"resource": scannable.Name(),
	})

	interval, err := scanner.checkInterval(scannable.CheckEvery())
	if err != nil {
		scanner.setResourceCheckError(logger, scannable, err)

		return 0, err
	}

	resourceTypes, err := scanner.dbPipeline.ResourceTypes()
	if err != nil {
		logger.Error("failed-to-get-resource-types", err)
		return 0, err
	}

	for _, parentType := range resourceTypes {
		if parentType.Name() == scannable.Name() {
			continue
		}
		if parentType.Name() != scannable.Type() {
			continue
		}

		logger.Info("parent-resource-type-has-no-versions", lager.Data{"resource-type": parentType.Name()})
		return interval, nil
		// err = scanner.customScanner.Scan(logger.Session("resource-type-scanner"), parentType)
		// if err != nil {
		// 	logger.Error("failed-to-scan-parent-resource-type-version", err)
		// 	scanner.setResourceCheckError(logger, scannable, err)
		// 	return 0, err
		// }
	}

	// resourceTypes, err = scanner.dbPipeline.ResourceTypes()
	// if err != nil {
	// 	logger.Error("failed-to-get-resource-types", err)
	// 	return 0, err
	// }

	versionedResourceTypes := creds.NewVersionedResourceTypes(
		scanner.variables,
		resourceTypes.Deserialize(),
	)

	source, err := creds.NewSource(scanner.variables, scannable.Source()).Evaluate()
	if err != nil {
		logger.Error("failed-to-evaluate-resource-source", err)
		scanner.setResourceCheckError(logger, scannable, err)
		return 0, err
	}

	resourceConfigCheckSession, err := scanner.resourceConfigCheckSessionFactory.FindOrCreateResourceConfigCheckSession(
		logger,
		scannable.Type(),
		source,
		versionedResourceTypes.Without(scannable.Name()),
		ContainerExpiries,
	)
	if err != nil {
		logger.Error("failed-to-find-or-create-resource-config-check-session", err)
		scanner.setResourceCheckError(logger, scannable, err)
		return 0, err
	}

	resourceConfig := resourceConfigCheckSession.ResourceConfig()
	err = scannable.SetResourceConfig(resourceConfig.ID())
	if err != nil {
		logger.Error("failed-to-set-resource-config-id-on-resource", err)
		scanner.setResourceConfigCheckError(logger, resourceConfig, err)
		return 0, err
	}

	for breaker := true; breaker == true; breaker = mustComplete {
		lock, acquired, err := resourceConfig.AcquireResourceConfigCheckingLockWithIntervalCheck(
			logger,
			interval,
			mustComplete,
		)
		if err != nil {
			lockLogger.Error("failed-to-get-lock", err, lager.Data{
				"resource_name":   scannable.Name(),
				"resource_config": resourceConfig.ID(),
			})
			return interval, ErrFailedToAcquireLock
		}

		if !acquired {
			lockLogger.Debug("did-not-get-lock")
			if mustComplete {
				scanner.clock.Sleep(time.Second)
				continue
			} else {
				return interval, ErrFailedToAcquireLock
			}
		}

		defer lock.Release()

		break
	}

	if fromVersion == nil {
		rcv, found, err := resourceConfig.GetLatestVersion()
		if err != nil {
			logger.Error("failed-to-get-current-version", err)
			return interval, err
		}

		if found {
			fromVersion = atc.Version(rcv.Version())
		}
	}

	return interval, scanner.check(
		logger,
		scannable,
		resourceConfigCheckSession,
		fromVersion,
		versionedResourceTypes,
		source,
	)
}

func (scanner *resourceConfigScanner) check(
	logger lager.Logger,
	scannable Scannable,
	resourceConfigCheckSession db.ResourceConfigCheckSession,
	fromVersion atc.Version,
	resourceTypes creds.VersionedResourceTypes,
	source atc.Source,
) error {
	pipelinePaused, err := scanner.dbPipeline.CheckPaused()
	if err != nil {
		logger.Error("failed-to-check-if-pipeline-paused", err)
		return err
	}

	if pipelinePaused {
		logger.Debug("pipeline-paused")
		return nil
	}

	if scannable.Paused() {
		logger.Debug("resource-paused")
		return nil
	}

	found, err := scanner.dbPipeline.Reload()
	if err != nil {
		logger.Error("failed-to-reload-scannerdb", err)
		return err
	}
	if !found {
		logger.Info("pipeline-removed")
		return errPipelineRemoved
	}

	containerSpec := worker.ContainerSpec{
		ImageSpec: worker.ImageSpec{
			ResourceType: scannable.Type(),
		},
		Tags:   scannable.Tags(),
		TeamID: scanner.dbPipeline.TeamID(),
	}

	res, err := scanner.resourceFactory.NewResource(
		context.Background(),
		logger,
		db.NewResourceConfigCheckSessionContainerOwner(resourceConfigCheckSession, scanner.dbPipeline.TeamID()),
		db.ContainerMetadata{
			Type: db.ContainerTypeCheck,
		},
		containerSpec,
		resourceTypes.Without(scannable.Name()),
		worker.NoopImageFetchingDelegate{},
	)
	if err != nil {
		logger.Error("failed-to-initialize-new-container", err)
		scanner.setResourceConfigCheckError(logger, resourceConfigCheckSession.ResourceConfig(), err)
		return err
	}

	logger.Debug("checking", lager.Data{
		"from": fromVersion,
	})

	newVersions, err := res.Check(source, fromVersion)

	// XXX: change the metrics for resource and resource type check
	scanner.setResourceConfigCheckError(logger, resourceConfigCheckSession.ResourceConfig(), err)
	metric.ResourceCheck{
		PipelineName: scanner.dbPipeline.Name(),
		ResourceName: scannable.Name(),
		TeamName:     scanner.dbPipeline.TeamName(),
		Success:      err == nil,
	}.Emit(logger)

	if err != nil {
		if rErr, ok := err.(resource.ErrResourceScriptFailed); ok {
			logger.Info("check-failed", lager.Data{"exit-status": rErr.ExitStatus})
			return rErr
		}

		logger.Error("failed-to-check", err)
		return err
	}

	if len(newVersions) == 0 || reflect.DeepEqual(newVersions, []atc.Version{fromVersion}) {
		logger.Debug("no-new-versions")
		return nil
	}

	logger.Info("versions-found", lager.Data{
		"versions": newVersions,
		"total":    len(newVersions),
	})

	err = resourceConfigCheckSession.ResourceConfig().SaveVersions(newVersions)
	if err != nil {
		logger.Error("failed-to-save-resource-config-versions", err, lager.Data{
			"versions": newVersions,
		})
	}

	return nil
}

func swallowErrResourceScriptFailed(err error) error {
	if _, ok := err.(resource.ErrResourceScriptFailed); ok {
		return nil
	}
	return err
}

func (scanner *resourceConfigScanner) checkInterval(checkEvery string) (time.Duration, error) {
	interval := scanner.defaultInterval
	if checkEvery != "" {
		configuredInterval, err := time.ParseDuration(checkEvery)
		if err != nil {
			return 0, err
		}

		interval = configuredInterval
	}

	return interval, nil
}

func (scanner *resourceConfigScanner) setResourceCheckError(logger lager.Logger, scannable Scannable, err error) {
	setErr := scannable.SetCheckError(err)
	if setErr != nil {
		logger.Error("failed-to-set-check-error", err)
	}
}

func (scanner *resourceConfigScanner) setResourceConfigCheckError(logger lager.Logger, resourceConfig db.ResourceConfig, err error) {
	setErr := resourceConfig.SetCheckError(err)
	if setErr != nil {
		logger.Error("failed-to-set-check-error", err)
	}
}

var errPipelineRemoved = errors.New("pipeline removed")
