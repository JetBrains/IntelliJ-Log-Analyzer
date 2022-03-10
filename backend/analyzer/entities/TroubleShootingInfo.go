package entities

import (
	"bufio"
	"io"
	"log"
	"log_analyzer/backend/analyzer"
	"os"
	"regexp"
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
	r, _ := regexp.Compile("^Custom plugins:.*")
	s := r.FindString(currentString)
	if len(s) < 1 {
		return nil
	}
	s = strings.TrimPrefix(s, "Custom plugins:")
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "[")
	s = strings.TrimSuffix(s, "]")
	pluginsListAsString := strings.Split(s, ",")
	for _, plugin := range pluginsListAsString {
		version := plugin[strings.LastIndex(plugin, "(") : strings.LastIndex(plugin, ")")+1]
		plugin = strings.Replace(plugin, version, "", -1)
		plugin = strings.TrimSpace(plugin)
		version = strings.TrimPrefix(version, "(")
		version = strings.TrimSuffix(version, ")")
		log.Println(plugin + version)
		//todo: retreive plugin's link
		pluginsList = append(pluginsList, analyzer.IDEPlugin{
			Version: version,
			Name:    plugin,
			Link:    "https://plugins.jetbrains.com",
		})
	}

	return pluginsList
}

func findJRE(currentString string) string {
	r, _ := regexp.Compile("^JRE.*")
	s := r.FindString(currentString)
	return strings.TrimPrefix(s, "JRE:")
}

func findBuild(currentString string) string {
	r, _ := regexp.Compile("Build:\\s#([[:alnum:]|\\.|\\-]*)")
	return strings.TrimPrefix(r.FindString(currentString), "Build:")
}
