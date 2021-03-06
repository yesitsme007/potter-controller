package controllersdi

import (
	"github.com/gardener/potter-controller/pkg/deployutil"
	"github.com/gardener/potter-controller/pkg/helm"
	"github.com/gardener/potter-controller/pkg/kapp"
	"github.com/gardener/potter-controller/pkg/synchronize"
	"github.com/gardener/potter-controller/pkg/util"

	"github.com/pkg/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type DeployerFactory interface {
	GetDeployer(configType string) (deployutil.DeployItemDeployer, error)
}

func NewDeploymentFactory(crAndSecretClient client.Client, uncachedClient synchronize.UncachedClient,
	appRepoClient client.Client, blockObject *synchronize.BlockObject,
	reconcileIntervalMinutes int64) DeployerFactory {
	return &deployerFactoryImpl{
		crAndSecretClient:        crAndSecretClient,
		uncachedClient:           uncachedClient,
		appRepoClient:            appRepoClient,
		blockObject:              blockObject,
		reconcileIntervalMinutes: reconcileIntervalMinutes,
	}
}

type deployerFactoryImpl struct {
	crAndSecretClient        client.Client
	uncachedClient           synchronize.UncachedClient
	appRepoClient            client.Client
	blockObject              *synchronize.BlockObject
	reconcileIntervalMinutes int64
}

func (r *deployerFactoryImpl) GetDeployer(configType string) (deployutil.DeployItemDeployer, error) {
	var deployer deployutil.DeployItemDeployer
	switch configType {
	case util.ConfigTypeHelm:
		deployer = helm.NewHelmDeployerDI(r.crAndSecretClient, r.uncachedClient, r.appRepoClient, r.blockObject)
	case util.ConfigTypeKapp:
		deployer = kapp.NewKappDeployerDI(r.crAndSecretClient, r.uncachedClient, r.blockObject, r.reconcileIntervalMinutes)
	default:
		return nil, errors.New("Unsupported configtype " + configType)
	}

	return deployer, nil
}
