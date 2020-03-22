package ui

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/briandowns/spinner"
	"github.com/manifoldco/promptui"
	"github.com/nealhardesty/kubectlplugins/internal/global"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"github.com/nealhardesty/kubectlplugins/internal/util"
)

func SelectNamespaceWithUI(namespaceRegex string, contextName string, options *global.GlobalOptions) (string, error) {
	currentContext := strings.ToUpper(k8s.CurrentContextName())
	spin := spinner.New(spinner.CharSets[9], 100*time.Millisecond)
	spin.Prefix = fmt.Sprintf("loading namespaces in %v ", contextName)
	spin.Start()
	namespaces, err := k8s.GetNamespacesWithContext(contextName)
	if err != nil {
		util.Die(err)
	}
	spin.Stop()

	isCurrentContextMatch := false

	currentNamespace := k8s.CurrentNamespaceName()
	if namespaceRegex == "" || namespaceRegex == "." {
		namespaceRegex = fmt.Sprintf("^%v$", currentNamespace)
		isCurrentContextMatch = true
	}

	re := regexp.MustCompile(fmt.Sprintf("^%v", namespaceRegex))

	firstMatchedIndex := -1
	matchingNamespaceNames := []string{}
	for i, namespace := range namespaces {
		if re.MatchString(namespace.ObjectMeta.Name) {
			matchingNamespaceNames = append(matchingNamespaceNames, namespace.ObjectMeta.Name)
			if firstMatchedIndex == -1 {
				firstMatchedIndex = i
			}
		}
	}

	// If only one match, return it now
	if len(matchingNamespaceNames) == 1 && !isCurrentContextMatch {
		return matchingNamespaceNames[0], nil
	}

	promptui.IconSelect = "*"
	promptui.IconInitial = ""
	promptui.SearchPrompt = "FILTER >>> "
	prompt := promptui.Select{
		Label:        "SELECT A NAMESPACE IN " + currentContext,
		Size:         25, //len(namespaces),
		Items:        namespaces,
		CursorPos:    firstMatchedIndex,
		IsVimMode:    false,
		HideHelp:     true,
		HideSelected: true,
		Searcher: func(input string, i int) bool {
			namespaceName := namespaces[i].ObjectMeta.Name
			return strings.HasPrefix(namespaceName, input)
		},
		Templates: &promptui.SelectTemplates{
			Active:   "* {{ .ObjectMeta.Name | bold | blue | bgWhite }}",
			Inactive: "   {{ .ObjectMeta.Name | blue }}",
			Details:  "\n      STATUS: {{ .Status.Phase | bold }} {{ DefaultIfCurrentNamespace .ObjectMeta.Name | bold | blue }}\n     CREATED: {{ .ObjectMeta.CreationTimestamp | bold }}\n      LABELS: {{ .ObjectMeta.Labels | bold }}\n ANNOTATIONS: {{ .ObjectMeta.Annotations | bold }}\n",
			FuncMap:  promptui.FuncMap,
		},
	}
	prompt.Templates.FuncMap["DefaultIfCurrentNamespace"] = func(ns string) string {
		if ns == currentNamespace {
			return "DEFAULT"
		}
		return ""
	}
	//i, _, err := prompt.Run()
	i, _, err := prompt.RunCursorAt(firstMatchedIndex, firstMatchedIndex-5)
	if err != nil {
		return "", err
	}

	selectedNamespace := namespaces[i]

	//fmt.Printf("you selected %v", selectedNamespace)

	return selectedNamespace.ObjectMeta.Name, nil

}

func SwitchNamespaceByRegex(namespaceRegex string, options *global.GlobalOptions) (string, error) {
	namespace, err := SelectNamespaceWithUI(namespaceRegex, k8s.CurrentContextName(), options)
	if err != nil {
		return "", err
	}

	err = k8s.SwitchNamespace(namespace)
	if err != nil {
		return "", err
	}

	return namespace, nil

}
