package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/nealhardesty/kubectlplugins/internal/cmd/execute"
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/util"
)

func main() {
	subcommand := util.KubectlSubcommand()
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: kubectl %v [options] [podRegex] [command(default=/bin/sh)]*\n\n", subcommand)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\n")
	}
	options := global.NewGlobalOptionsWithFlag()

	var contextRegex, namespaceRegex string
	flag.StringVar(&contextRegex, "context", "", "kubecontext regex to search for (default to current context)")
	flag.StringVar(&contextRegex, "c", "", "kubecontext regex to search for (default to current context) (shorthand)")
	flag.StringVar(&namespaceRegex, "namespace", "", "namespace regex to search for (default to current namespace)")
	flag.StringVar(&namespaceRegex, "n", "", "namespace regex to search for (default to current namespace) (shorthand)")
	flag.Parse()

	args := flag.Args()

	podSelectorRegex := "."
	if len(args) > 0 {
		args = args[1:]
	}

	if len(args) == 0 {
		args = []string{"/bin/sh"}
	}

	err := execute.ExecByPodRegex(podSelectorRegex, contextRegex, namespaceRegex, args, options)
	if err != nil {
		util.Die(err.Error())
	}
}
