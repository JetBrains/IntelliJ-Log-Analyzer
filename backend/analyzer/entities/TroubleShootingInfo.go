package entities

import (
	"bufio"
	"io"
	"log"
	"log_analyzer/backend/analyzer"
	"os"
	"strings"
)

func init() {
	CurrentAnalyzer.AddStaticEntity(analyzer.StaticEntity{
		Name:                "troubleshooting.txt",
		ConvertToStaticInfo: parseTroubleshootingInfo,
		CheckPath:           isTroubleshootingInfo,
	})
}
func isTroubleshootingInfo(path string) bool {
	if strings.Contains(path, "troubleshooting.txt") {
		return true
	}
	return false
}

func parseTroubleshootingInfo(path string) (a analyzer.StaticInfo) {
	reader, _ := os.Open(path)
	bufReader := bufio.NewReader(reader)
	for {
		currentString, err := bufReader.ReadString('\n')
		var build string
		if build = findBuild(currentString); len(build) > 0 {
			a.Build = build
		}
		if jre := findJRE(currentString); len(jre) > 0 {
			a.JRE = jre
		}
		if customPLuginsList := findCustomPlugins(currentString); len(customPLuginsList) > 0 {
			a.PluginsList = customPLuginsList
		}
		if os := findOS(currentString); len(os) > 0 {
			a.OS = os
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("ERROR: %s", err)
		}
	}
	return a
}

func findCustomPlugins(currentString string) (pluginsList []analyzer.IDEPlugin) {
	s := analyzer.GetRegexNamedCapturedGroups(`^Custom plugins: \[(?P<PluginsString>.*)\]`, currentString)["PluginsString"]
	if len(s) == 0 {
		return nil
	}
	pluginsListAsSlice := strings.Split(s, ",")
	for _, pluginAsString := range pluginsListAsSlice {
		s := analyzer.GetRegexNamedCapturedGroups(`^\s*(?P<Plugin>.*)\s+\((?P<Version>.*)\)$`, pluginAsString)
		version := s["Version"]
		name := s["Plugin"]
		//todo: retreive plugin's link
		pluginsList = append(pluginsList, analyzer.IDEPlugin{
			Version: version,
			Name:    name,
			Link:    "https://plugins.jetbrains.com/search?search=" + name,
		})
	}

	return pluginsList
}

func findJRE(currentString string) string {
	s := analyzer.GetRegexNamedCapturedGroups(`^JRE:\s*(?P<JRE>.*)`, currentString)["JRE"]
	return s
}

func findBuild(currentString string) string {
	s := analyzer.GetRegexNamedCapturedGroups(`Build:\s(?P<Build>#.*)`, currentString)["Build"]
	return s
}

func findOS(currentString string) string {
	s := analyzer.GetRegexNamedCapturedGroups(`^Operating System:\s*(?P<OS>.*)`, currentString)["OS"]
	return s
}
