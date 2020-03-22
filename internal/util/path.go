package util

import (
	"os"
	"path"
	"path/filepath"
	"strings"
)

// KubectlSubcommand attempt to determine the subcommand
func KubectlSubcommand() string {
	fullCommand := os.Args[0]
	base := path.Base(fullCommand)
	subcommand := strings.Replace(base, "kubectl-", "", 1)
	subcommand = strings.Replace(subcommand, ".exe", "", 1)
	return subcommand
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}
	return os.Getenv("USERPROFILE") // windows
}

func DefaultKubeconfigPath() string {
	return filepath.Join(homeDir(), ".kube", "config")
}
