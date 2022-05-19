package backend

import (
	"bytes"
	"embed"
	"encoding/json"
	"html/template"
	"io/ioutil"
	"log"
	"os"
	"path"
	"reflect"
)

//go:embed Config.gohtml
var tmplFS embed.FS

var (
	ConfigurationOptions = Config{}
	defaultConfig        = Config{
		EditorFontSize: 12,
		EditorTheme:    "system",
	}
	ConfigFileName      = "config.json"
	ConfigDirectoryName = path.Clean("JetBrains/IntelliJLogAnalyzer")
)

type Config struct {
	EditorFontSize             int    `json:"EditorFontSize"`
	EditorTheme                string `json:"EditorTheme"`
	EditorDefaultSoftWrapState bool   `json:"EditorDefaultSoftWrapState"`
}

func GetConfig() *Config {
	if ConfigurationOptions != *new(Config) {
		return &ConfigurationOptions
	} else {
		ConfigurationOptions = generateConfig()
	}
	return &ConfigurationOptions
}
func GetSettingsScreenHTML() string {
	config := GetConfig()
	var tpl bytes.Buffer
	t := template.Must(template.New("Config.gohtml").
		ParseFS(tmplFS, "Config.gohtml"))
	err := t.Execute(&tpl, config)
	if err != nil {
		log.Printf("Template Config.gohtml execution failed. Error: %s", err.Error())
	}
	return tpl.String()
}
func (c *Config) saveConfig() {
	configPath := getConfigFilePath()
	if !FileExists(configPath) {
		createConfig(configPath)
	}
	configFileContent, err := json.Marshal(c)
	if err != nil {
		log.Println("Error while marshalling config file content: ", err)
	}
	err = ioutil.WriteFile(configPath, configFileContent, 0644)
	if err != nil {
		log.Println("Error while writing config file: ", err)
	}
}
func (c *Config) SaveSetting(id string, value interface{}) {
	log.Printf("Saving setting '%s' with value '%v'", id, value)
	aValue := reflect.ValueOf(value)
	fields := reflect.VisibleFields(reflect.TypeOf(Config{}))
	for _, field := range fields {
		if aValue.CanConvert(field.Type) {
			if field.Name == id {
				reflect.ValueOf(c).Elem().FieldByName(field.Name).Set(aValue.Convert(field.Type))
			}
		}
	}
	c.saveConfig()
}

func generateConfig() Config {
	configPath := getConfigFilePath()
	if !FileExists(configPath) {
		log.Printf("Could not open configuration file: %s \n Using default config", configPath)
	}
	configValues := getConfigValues(configPath)
	if configValues != *new(Config) {
		return configValues
	} else {
		return defaultConfig
	}
}
func getConfigValues(configPath string) Config {
	configFile, err := os.Open(configPath)
	if err != nil {
		log.Printf("Could not open configuration file: %s", configPath)
		return *new(Config)
	}
	defer configFile.Close()

	byteValue, _ := ioutil.ReadAll(configFile)
	var config Config
	json.Unmarshal(byteValue, &config)
	return config
}
func getConfigFilePath() string {
	return getConfigDir() + string(os.PathSeparator) + ConfigDirectoryName + string(os.PathSeparator) + ConfigFileName
}
func getConfigDir() string {
	configPath, err := os.UserConfigDir()
	if err != nil {
		log.Printf("Could not get system configuration directory: %v \n Checking config file in home directory", err)
		configPath, _ = os.UserHomeDir()
	}
	return configPath
}
func createConfig(configFilePath string) {
	configDir := path.Dir(configFilePath)
	if err, ConfigDirectoryExists := os.Stat(configDir); ConfigDirectoryExists == nil {
		log.Printf("Creating configuration file in %s", configDir)
		createConfigFile(configFilePath)
	} else {
		log.Printf("Could not find configuration directory: %v", err)
		err := createConfigDirectory(configDir)
		if err == nil {
			createConfigFile(configFilePath)
		} else {
			log.Printf("Could not create configuration directory: %v", err)
		}
	}
}
func createConfigFile(configFilePath string) {
	configFile, err := os.Create(configFilePath)
	defer func(configFile *os.File) {
		err := configFile.Close()
		if err != nil {
			log.Printf("Could not close configuration file: %v", err)
		}
	}(configFile)

	if err != nil {
		log.Printf("Could not create configuration file: %v", err)
	}
	log.Printf("Successfully created configuration file: %s", configFilePath)
}
func createConfigDirectory(path string) error {
	err := os.MkdirAll(path, os.ModePerm)
	if err != nil {
		log.Printf("Could not create configuration directory: %v", err)
		return err
	}
	return nil
}
func FileExists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		return false
	}
	return true
}
