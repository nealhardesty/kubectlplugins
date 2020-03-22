package ui

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/util"
	v1 "k8s.io/api/core/v1"
)

type containerInfo struct {
	DisplayName     string
	Pod             v1.Pod
	ContainerStatus v1.ContainerStatus
}

func filterPodsIntoContainerInfos(podRegex string, pods []v1.Pod) []containerInfo {
	if podRegex == "" {
		podRegex = "."
	}

	re := regexp.MustCompile(fmt.Sprintf("^%v", podRegex))

	filteredContainerInfos := []containerInfo{}
	for _, pod := range pods {
		if re.MatchString(pod.Name) {
			if pod.Status.Phase == "Running" {
				for _, status := range pod.Status.InitContainerStatuses {
					if status.State.Running != nil {
						filteredContainerInfos = append(filteredContainerInfos, containerInfo{
							DisplayName:     fmt.Sprintf("%v (init container: %v)", pod.Name, status.Name),
							Pod:             pod,
							ContainerStatus: status,
						})
					}
				}
				for _, status := range pod.Status.ContainerStatuses {
					if status.State.Running != nil {
						filteredContainerInfos = append(filteredContainerInfos, containerInfo{
							DisplayName:     fmt.Sprintf("%v (%v)", pod.Name, status.Name),
							Pod:             pod,
							ContainerStatus: status,
						})
					}
				}
			}
		}
	}

	return filteredContainerInfos

}

// SelectPodContainerWithUI ...
func SelectPodContainerWithUI(podRegex string, contextName string, namespace string, options *global.GlobalOptions) (pod string, container string, err error) {
	//fmt.Printf("loading pods in %v:%v ", contextName, namespace)
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Prefix = fmt.Sprintf("loading pods in %v:%v ", contextName, namespace)
	spin.Start()
	pods, err := k8s.GetPodsWithContext(contextName, namespace)
	if err != nil {
		util.Die(err)
	}
	spin.Stop()
	util.ClearLine(os.Stdout)

	if len(pods) == 0 {
		return "", "", fmt.Errorf("No pods found in %v:%v", namespace, contextName)
	}

	filteredContainerInfos := filterPodsIntoContainerInfos(podRegex, pods)

	if len(filteredContainerInfos) == 1 {
		// Only one pod matched, no need to ask
		return filteredContainerInfos[0].Pod.Name, filteredContainerInfos[0].ContainerStatus.Name, nil
	}
	//fmt.Printf("matched pods %v", filteredContainerInfos)

	if len(filteredContainerInfos) == 0 {
		// No match, widen the search to everything
		filteredContainerInfos = filterPodsIntoContainerInfos(".", pods)
	}

	contextNameUpper := strings.ToUpper(contextName)
	namespaceUpper := strings.ToUpper(namespace)

	promptui.IconSelect = "*"
	promptui.IconInitial = ""
	promptui.SearchPrompt = "FILTER >>> "
	prompt := promptui.Select{
		Label:        fmt.Sprintf("SELECT A POD (CONTAINER) IN %s:%s", contextNameUpper, namespaceUpper),
		Size:         25, //len(namespaces),
		Items:        filteredContainerInfos,
		CursorPos:    0,
		IsVimMode:    false,
		HideHelp:     true,
		HideSelected: true,
		Searcher: func(input string, i int) bool {
			displayName := filteredContainerInfos[i].DisplayName
			return strings.HasPrefix(displayName, input)
		},
		Templates: &promptui.SelectTemplates{
			Active:   "* {{ .DisplayName | bold | blue | bgWhite }}",
			Inactive: "   {{ .DisplayName | blue }}",
			//Details:  "\n     STARTED: {{ .ObjectMeta.StartTime }}\tCREATED: {{ .ObjectMeta.CreationTimestamp }}\tRESTARTS: {{ .ContainerStatus.RestartCount }}\n      LABELS: {{ .ObjectMeta.Labels | bold }}\n ANNOTATIONS: {{ .ObjectMeta.Annotations | bold }}\n    IMAGE: {{ .ContainerStatus.Image }}\n",
			Details: "\n     STARTED: {{ .Pod.Status.StartTime }}\tCREATED: {{ .Pod.ObjectMeta.CreationTimestamp }}\tRESTARTS: {{ .ContainerStatus.RestartCount }}", //\n      LABELS: {{ .Pod.ObjectMeta.Labels | bold }}\n ANNOTATIONS: {{ .Pod.ObjectMeta.Annotations | bold }}\n    IMAGE: {{ .ContainerStatus.Image }}\n",
			FuncMap: promptui.FuncMap,
		},
	}

	i, _, err := prompt.Run()
	if err != nil {
		return "", "", err
	}

	return filteredContainerInfos[i].Pod.Name, filteredContainerInfos[i].ContainerStatus.Name, nil
}
