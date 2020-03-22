package context

import (
	"fmt"
	"os"

	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/ui"
)

func PrintContextsAndExit() {
	contexts := k8s.GetContextNamesSorted()
	currentContext := k8s.CurrentContextName()
	for _, name := range contexts {
		if name == currentContext {
			fmt.Println(" * " + name)
		} else {
			fmt.Println("   " + name)
		}
	}
	os.Exit(0)
}

func SwitchContextByRegex(contextNameRegex string, options *global.GlobalOptions) (string, error) {
	if options.PrintOnly {
		PrintContextsAndExit()
	}
	newContextName, err := ui.SelectContextWithUI(contextNameRegex)

	if err != nil {
		//fmt.Fprintf(os.Stderr, "%v\n\n", err)
		return "", err
	}

	err = k8s.SwitchContext(newContextName)
	if err != nil {
		return newContextName, err
	}

	return newContextName, nil
}
