package kubernetes

import (
	"code.cloudfoundry.org/lager"

	"github.com/concourse/atc/creds"
	"k8s.io/client-go/kubernetes"
)

type kubernetesFactory struct {
	clientset        *kubernetes.Clientset
	logger           lager.Logger
	namespacePrefix  string
	defaultNamespace string
	secretsName      string
}

func NewKubernetesFactory(logger lager.Logger, clientset *kubernetes.Clientset, namespacePrefix string, defaultNameSpace string, secretsName string) *kubernetesFactory {
	factory := &kubernetesFactory{
		clientset:        clientset,
		logger:           logger,
		namespacePrefix:  namespacePrefix,
		defaultNamespace: defaultNameSpace,
		secretsName:      secretsName,
	}

	return factory
}

func (factory *kubernetesFactory) NewVariables(teamName string, pipelineName string) creds.Variables {
	return &Kubernetes{
		Clientset:        factory.clientset,
		TeamName:         teamName,
		PipelineName:     pipelineName,
		NamespacePrefix:  factory.namespacePrefix,
		DefaultNamespace: factory.defaultNamespace,
		SecretsName:      factory.secretsName,
		logger:           factory.logger,
	}
}
