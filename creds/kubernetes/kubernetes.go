package kubernetes

import (
	"code.cloudfoundry.org/lager"
	"github.com/cloudfoundry/bosh-cli/director/template"
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	k8s_errors "k8s.io/apimachinery/pkg/api/errors"
)

type Kubernetes struct {
	Clientset    *kubernetes.Clientset
	TeamName     string
	PipelineName string
	NamespacePrefix string
	Logger       lager.Logger
}

func (k Kubernetes) Get(varDef template.VariableDefinition) (interface{}, bool, error) {
	var namespace = k.NamespacePrefix + k.TeamName

	//first try
	secret, found, err := k.findSecret(namespace, k.PipelineName + "." + varDef.Name)

	if err != nil {
		return nil, false, err
	}

	if found {
		return k.getValueFromSecret(secret)
	}

	//fallback name
	secret, found, err = k.findSecret(namespace, varDef.Name)

	if err != nil {
		return nil, false, err
	}

	if found {
		return k.getValueFromSecret(secret)
	}

	//didn't find it.
	return nil, false, nil
}

func (k Kubernetes) getValueFromSecret(secret *v1.Secret) (interface{}, bool, error) {
	val, found := secret.Data["value"]
	if found {
		return val, true, nil
	}

	evenLessTyped := map[interface{}]interface{}{}
	for k, v := range secret.Data {
		evenLessTyped[k] = v
	}

	return evenLessTyped, true, nil
}

func (k Kubernetes) defautSecretName() string {
	return k.PipelineName + "-concourse-secrets"
}

func (k Kubernetes) findSecret(namespace, name string) (*v1.Secret, bool, error) {
	var secret *v1.Secret
	var err error

	secret, err = k.Clientset.Core().Secrets(namespace).Get(name, meta_v1.GetOptions{})

	if err != nil && k8s_errors.IsNotFound(err) {
		return nil, false, nil
	} else if err != nil {
		return nil, false, err
	} else {
		return secret, true, err
	}
}

func (k Kubernetes) List() ([]template.VariableDefinition, error) {
	// Don't think this works with vault.. if we need it to we'll figure it out
	// var defs []template.VariableDefinition

	// secret, err := v.vaultClient.List(v.PathPrefix)
	// if err != nil {
	// 	return defs, err
	// }

	// var def template.VariableDefinition
	// for name, _ := range secret.Data {
	// 	defs := append(defs, template.VariableDefinition{
	// 		Name: name,
	// 	})
	// }

	return []template.VariableDefinition{}, nil
}