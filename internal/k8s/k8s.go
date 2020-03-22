package k8s

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"sort"

	//kcexec "github.com/kubernetes/kubectl/pkg/cmd/exec"
	"github.com/nealhardesty/kubectlplugins/internal/util"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/tools/clientcmd/api"
)

// GetConfigAccess Load clientcmd's ConfigAccess (PathOptions)
func GetConfigAccess() clientcmd.ConfigAccess {
	return clientcmd.NewDefaultPathOptions()
}

// GetStartingConfigOrDie Attempt to load clientconfig's default configuration
func GetStartingConfigOrDie() *api.Config {
	pathOptions := GetConfigAccess()

	config, err := pathOptions.GetStartingConfig()
	if err != nil {
		util.Die(err)
	}

	return config
}

func GetContexts() map[string]*api.Context {
	config := GetStartingConfigOrDie()
	return config.Contexts
}

func GetContextNamesSorted() []string {
	contexts := GetContexts()
	contextNames := []string{}
	for name := range contexts {
		contextNames = append(contextNames, name)
	}
	sort.Strings(contextNames)
	return contextNames

}

func GetConfigWithContextOrDie(contextName string) *rest.Config {
	kubeconfig := util.DefaultKubeconfigPath()
	config, err := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
		&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig},
		&clientcmd.ConfigOverrides{CurrentContext: contextName},
	).ClientConfig()

	if err != nil {
		util.Die(err)
	}

	return config
}

func GetConfigOrDie() *rest.Config {
	kubeconfig := util.DefaultKubeconfigPath()

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		util.Die(err)
	}

	return config
}

func GetClientsetOrDie() *kubernetes.Clientset {
	return kubernetes.NewForConfigOrDie(GetConfigOrDie())
}

func GetClientsetWithContextOrDie(contextName string) *kubernetes.Clientset {
	return kubernetes.NewForConfigOrDie(GetConfigWithContextOrDie(contextName))

}

func CurrentContextName() string {
	return GetStartingConfigOrDie().CurrentContext
}

func CurrentNamespaceName() string {
	config := GetStartingConfigOrDie()
	currentContext := config.Contexts[config.CurrentContext]

	return currentContext.Namespace

}

func CurrentNamespaceNameWithContext(contextName string) (string, error) {
	config := GetStartingConfigOrDie()
	context, exists := config.Contexts[contextName]
	if !exists {
		return "", fmt.Errorf("Context '%v' does not exist", contextName)
	}

	return context.Namespace, nil

}

func GetNamespacesWithContext(contextName string) ([]v1.Namespace, error) {
	clientset := GetClientsetWithContextOrDie(contextName)
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return namespaces.Items, err
}

func GetNamespaces() ([]v1.Namespace, error) {
	clientset := GetClientsetOrDie()
	namespaces, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	return namespaces.Items, err
}

func GetNamespaceNamesOrDie() []string {
	namespaces, err := GetNamespaces()
	if err != nil {
		util.Die(err)
	}
	namespaceNames := []string{}
	for _, namespace := range namespaces {
		namespaceNames = append(namespaceNames, namespace.Name)
	}
	sort.Strings(namespaceNames)
	return namespaceNames

}

func GetPodsWithContext(contextName string, namespace string) ([]v1.Pod, error) {
	clientset := GetClientsetWithContextOrDie(contextName)
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	return pods.Items, err
}

func GetPods(namespace string) ([]v1.Pod, error) {
	clientset := GetClientsetOrDie()
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	return pods.Items, err
}

// SwitchContext ...
func SwitchContext(contextName string) error {
	config := GetStartingConfigOrDie()
	config.CurrentContext = contextName

	return clientcmd.ModifyConfig(GetConfigAccess(), *config, true)
}

// SwitchNamespace ...
func SwitchNamespace(namespace string) error {
	currentContextName := CurrentContextName()

	config := GetStartingConfigOrDie()
	config.Contexts[currentContextName].Namespace = namespace

	return clientcmd.ModifyConfig(GetConfigAccess(), *config, true)
}

func Execute(contextName string, namespaceName string, podName string, containerName string, args []string) (exitCode int, err error) {
	fmt.Fprintf(os.Stderr, "Using pod %v container %v...\n", podName, containerName)

	kubectlPath, err := exec.LookPath("kubectl")
	if err != nil {
		return 0, err
	}

	kubectlArgs := []string{
		kubectlPath,
		"exec",
		"--tty=true",
		"--stdin=true",
		"--namespace",
		namespaceName,
		"--context",
		contextName,
		"--container",
		containerName,
		"pod/" + podName,
		"--",
	}

	for _, arg := range args {
		kubectlArgs = append(kubectlArgs, arg)
	}

	cmd := exec.Cmd{
		Path:   kubectlPath,
		Args:   kubectlArgs,
		Stdin:  os.Stdin,
		Stdout: os.Stdout,
		Stderr: os.Stderr,
	}

	err = cmd.Run()
	if err != nil {
		return -1, err
	}

	return cmd.ProcessState.ExitCode(), nil

}

/*
func Execute(contextName string, namespaceName string, podName string, containerName string, args []string, options *GlobalOptions) error {
	fmt.Fprintf(os.Stderr, "Execute pod=%v container=%v namespace=%v context=%v args=%v\n",
		podName, containerName, namespaceName, contextName, args)

	var err error

	execOptions := &kcexec.ExecOptions{
		PodClient:               GetClientsetWithContextOrDie(contextName).CoreV1(),
		Config:                  GetConfigWithContextOrDie(contextName),
		ResourceName:            "pod",
		EnforceNamespace:        true,
		EnableSuggestedCmdUsage: true,
		Command:                 args,
		StreamOptions: kcexec.StreamOptions{
			Namespace:     namespaceName,
			PodName:       podName,
			ContainerName: containerName,
			Stdin:         true,
			TTY:           true,
			IOStreams: genericclioptions.IOStreams{
				In:     os.Stdin,
				Out:    os.Stdout,
				ErrOut: os.Stderr,
			},
		},
	}

	err = execOptions.Validate()
	if err != nil {
		Die(err)
	}

	err = execOptions.Run()

	return nil

}
*/
