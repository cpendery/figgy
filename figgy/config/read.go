package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cpendery/figgy/figgy/config/cfg"
	yaml "gopkg.in/yaml.v3"
)

const FiggyConfigName = ".fig.yaml"
const FiggyYamlKey = ".figgy_file"

func GetAllFiggiedConfigs(path string) (*[]string, error) {
	yfile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(yfile)
	figgiedConfigs := []string{}
	var node yaml.Node
	for {
		if err := decoder.Decode(&node); err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			return &figgiedConfigs, nil
		}
		data := make(map[string]interface{})
		err := node.Decode(&data)
		if err != nil {
			return nil, err
		}
		//use node.HeadComment once fixed upstream
		if figgiedConfig, exists := data[FiggyYamlKey]; exists {
			figgiedConfigs = append(figgiedConfigs, figgiedConfig.(string))
		}
	}
}

func ReadFiggyConfig(path string) (ParsedConfig, error) {
	yfile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	decoder := yaml.NewDecoder(yfile)
	figgyConfig := make(map[string]interface{})
	var node yaml.Node
	for {
		if err := decoder.Decode(&node); err != nil {
			if err.Error() != "EOF" {
				return nil, err
			}
			return figgyConfig, nil
		}
		data := make(map[string]interface{})
		err := node.Decode(&data)
		if err != nil {
			return nil, err
		}
		//use node.HeadComment once fixed upstream
		if figgiedConfig, exists := data[".figgy_file"]; exists {
			figgyConfig[figgiedConfig.(string)] = data
		}
	}
}

func ReadConfig(filename string) (ParsedConfig, error) {
	extension := filepath.Ext(filename)
	var fig ParsedConfig
	var err error
	switch extension {
	case ".cfg":
		fig, err = cfg.Load(filename)
	default:
		fig, err = nil, fmt.Errorf("unsupported file extension: %s", extension)
	}
	if fig != nil {
		fig[FiggyYamlKey] = filename
	}
	return fig, err
}
