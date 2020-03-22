package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nealhardesty/kubectlplugins/internal/cmd/context"
	"github.com/nealhardesty/kubectlplugins/internal/cmd/namespace"
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/util"
)

func main() {
	subcommand := util.KubectlSubcommand()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: kubectl %v [options] [contextRegex] [namespaceRegex]\n\n", subcommand)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	options := global.NewGlobalOptionsWithFlag()
	flag.Parse()

	args := flag.Args()

	contextSelectorRegex := ""
	if len(args) > 0 {
		contextSelectorRegex = args[0]
	}

	// Atempt to set context
	newContextName, err := context.SwitchContextByRegex(contextSelectorRegex, options)
	if err != nil {
		util.Die(err.Error())
	}
	//printContexts(os.Stdout)
	fmt.Fprintf(os.Stdout, "switched context to '%v'\n", newContextName)

	if len(args) > 1 {
		namespaceRegex := args[1]
		namespace, err := namespace.SwitchNamespaceByRegex(namespaceRegex, options)
		if err != nil {
			util.Die(err.Error())
		}
		fmt.Fprintf(os.Stdout, "switched namespace to '%v'\n", namespace)
	}

}
