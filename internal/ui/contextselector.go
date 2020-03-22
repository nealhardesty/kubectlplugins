package ui

import (
	"fmt"
	"regexp"
	"sort"

	"github.com/manifoldco/promptui"
	"github.com/nealhardesty/kubectlplugins/internal/k8s"
	"k8s.io/client-go/tools/clientcmd/api"
)

type contextInfo struct {
	Name    string
	Context *api.Context
}

func sortContexts(m map[string]*api.Context) []*contextInfo {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	sortedContexts := make([]*contextInfo, 0, len(m))
	for _, key := range keys {
		sortedContexts = append(sortedContexts, &contextInfo{Name: key, Context: m[key]})
	}

	return sortedContexts
}

/*
func printContexts(out io.Writer) {
	config := util.GetStartingConfigOrDie()
	currentContextName := config.CurrentContext
	contexts := config.Contexts

	maxWidth := 0
	for contextName := range contexts {
		if len(contextName) > maxWidth {
			maxWidth = len(contextName)
		}
	}
	maxWidth = maxWidth + 2

	format := fmt.Sprintf("%%v %%-%dv\t(%%v)\n", maxWidth)
	//fmt.Fprintf(out, format, " ", "CONTEXT", "NAMESPACE")
	for _, contextName := range sortedKeysForStringContextMap(contexts) {
		context := contexts[contextName]
		defaultColumn := " "
		if contextName == currentContextName {
			defaultColumn = "*"
		}
		fmt.Fprintf(out, format, defaultColumn, contextName, context.Namespace)
	}
}
*/

func SelectContextWithUI(contextNameRegex string) (newContextName string, err error) {
	config := k8s.GetStartingConfigOrDie()

	// Is the search undefined enough to just warrant starting with the current context
	isCurrentContextMatch := false
	if contextNameRegex == "." {
		contextNameRegex = fmt.Sprintf("^%v$", config.CurrentContext)
		isCurrentContextMatch = true
	}
	re := regexp.MustCompile(fmt.Sprintf("^%v", contextNameRegex))

	matchingContextNames := []string{}
	sortedContextInfos := sortContexts(config.Contexts)
	firstMatchingIndex := -1

	for i, contextInfo := range sortedContextInfos {
		if re.MatchString(contextInfo.Name) {
			if firstMatchingIndex == -1 {
				firstMatchingIndex = i
			}
			matchingContextNames = append(matchingContextNames, contextInfo.Name)
		}
	}

	if len(matchingContextNames) == 1 && !isCurrentContextMatch {
		return matchingContextNames[0], nil
	}

	//promptui.IconSelect = promptui.Styler(promptui.FGBold)("*")
	promptui.IconSelect = "*"
	promptui.IconInitial = ""
	promptui.SearchPrompt = "FILTER >>> "
	prompt := promptui.Select{
		Label:        "SELECT A CONTEXT",
		Size:         25, //len(sortedContextInfos),
		Items:        sortedContextInfos,
		CursorPos:    firstMatchingIndex,
		IsVimMode:    false,
		HideHelp:     true,
		HideSelected: true,
		Templates: &promptui.SelectTemplates{
			Active:   "* {{ .Name | bold | blue | bgWhite }}",
			Inactive: "   {{ .Name | blue }}",
			Details:  "\n NAMESPACE: {{ .Context.Namespace | bold }}\n   CLUSTER: {{ .Context.Cluster | bold }}\n    ORIGIN: {{ .Context.LocationOfOrigin | bold }}",
		},
	}
	selectedIndex, _, err := prompt.RunCursorAt(firstMatchingIndex, firstMatchingIndex-5)
	if err != nil {
		return "", err
	}

	return sortedContextInfos[selectedIndex].Name, nil
}
