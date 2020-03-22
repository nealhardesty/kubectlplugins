package execute

import (
	"fmt"
	"os"

	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/ui"
	"github.com/nealhardesty/kubectlplugins/internal/util"
)

// ExecByPodRegex TODO
func ExecByPodRegex(podSelectorRegex string, contextSelectorRegex string, namespaceSelectorRegex string, args []string, options *global.GlobalOptions) error {
	var err error

	context := k8s.CurrentContextName()
	if contextSelectorRegex != "" {
		context, err = ui.SelectContextWithUI(contextSelectorRegex)
		if err != nil {
			util.Die(err.Error())
		}
	}

	fmt.Fprintf(os.Stderr, "Using context %v\n", context)

	namespace, _ := k8s.CurrentNamespaceNameWithContext(context)
	if namespaceSelectorRegex != "" {
		namespace, err = ui.SelectNamespaceWithUI(namespaceSelectorRegex, context, options)
		if err != nil {
			util.Die(err.Error())
		}
	}

	fmt.Fprintf(os.Stderr, "Using namespace %v\n", namespace)

	pod, container, err := ui.SelectPodContainerWithUI(podSelectorRegex, context, namespace, options)
	if err != nil {
		util.Die(err.Error())
	}

	exitCode, err := k8s.Execute(context, namespace, pod, container, args)
	if err != nil {
		util.Die(err)
	}

	os.Exit(exitCode)

	return nil
}
