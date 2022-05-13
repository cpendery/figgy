package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/cpendery/figgy/figgy/engines"
	_ "github.com/cpendery/figgy/figgy/engines/all"
	"gopkg.in/yaml.v3"
)

func WriteFiggyConfig(fig ParsedConfig) error {
	f, err := os.OpenFile(FiggyConfigName, os.O_CREATE, 444)
	if err != nil {
		return err
	}
	encoder := yaml.NewEncoder(f)
	defer f.Close()
	defer encoder.Close()
	for _, configObj := range fig {
		if err := encoder.Encode(configObj); err != nil {
			return err
		}
	}
	return nil
}

func WriteConfig(fig ParsedConfig) error {
	extension := ""
	path := ""
	if currPath, exists := fig[FiggyYamlKey]; exists {
		path = currPath.(string)
		extension = filepath.Ext(path)
	}

	delete(fig, FiggyYamlKey)

	engine, exists := engines.Get(extension)
	if exists {
		return engine.Write(path, fig)
	}
	return fmt.Errorf("unsupported file extension: %s", extension)
}
