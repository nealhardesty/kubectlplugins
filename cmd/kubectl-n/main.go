package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nealhardesty/kubectlplugins/internal/cmd/namespace"
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/util"
)

func PrintNamespacesAndExit() {
	namespaces := k8s.GetNamespaceNamesOrDie()
	currentNamespace := k8s.CurrentNamespaceName()
	for _, name := range namespaces {
		if name == currentNamespace {
			fmt.Println(" * " + name)
		} else {
			fmt.Println("   " + name)
		}
	}
	os.Exit(0)
}

func main() {
	subcommand := util.KubectlSubcommand()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: kubectl %v [options] [namespaceRegex]\n\n", subcommand)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	options := global.NewGlobalOptionsWithFlag()
	flag.Parse()

	args := flag.Args()

	if options.PrintOnly {
		PrintNamespacesAndExit()
	}

	namespaceSelectorRegex := "."
	if len(args) > 0 {
		namespaceSelectorRegex = args[0]
	}

	namespace, err := namespace.SwitchNamespaceByRegex(namespaceSelectorRegex, options)
	if err != nil {
		util.Die(err.Error())
	}
	fmt.Fprintf(os.Stdout, "switched namespace to '%v'\n", namespace)

}
