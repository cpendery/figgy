package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cpendery/figgy/figgy/config/cfg"
	yaml "gopkg.in/yaml.v3"
)

const FiggyConfigName = ".fig.yaml"

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
		if figgiedConfig, exists := data[".figgy_file"]; exists {
			figgiedConfigs = append(figgiedConfigs, figgiedConfig.(string))
		}
	}
}

func ReadFiggyConfig(path string) (ParsedConfig, error) {
	return nil, nil
}

func ReadConfig(filename string) (ParsedConfig, error) {
	extension := filepath.Ext(filename)
	switch extension {
	case ".cfg":
		return cfg.Load(filename)
	default:
		return nil, fmt.Errorf("unsupported file extension: %s", extension)
	}
}
