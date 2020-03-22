package namespace

import (
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/ui"
)

// SwitchNamespaceByRegex
func SwitchNamespaceByRegex(namespaceRegex string, options *global.GlobalOptions) (string, error) {
	namespace, err := ui.SelectNamespaceWithUI(namespaceRegex, k8s.CurrentContextName(), options)
	if err != nil {
		return "", err
	}

	err = k8s.SwitchNamespace(namespace)
	if err != nil {
		return "", err
	}

	return namespace, nil

}
