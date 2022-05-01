package config

import (
	"os"

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
