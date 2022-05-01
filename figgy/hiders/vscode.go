package hiders

import (
	"encoding/json"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cpendery/figgy/figgy/config"
)

const vscodeFolderPath string = "./.vscode"
const vscodeSettingsFile string = "/settings.json"
const vscodeSettingsPath string = vscodeFolderPath + vscodeSettingsFile
const vscodeExcludeFileSetting string = "files.exclude"
const vscodeSettingFileMode fs.FileMode = 444

type VSCodeHider struct{}

func (v *VSCodeHider) Hide() error {
	figgyConfigs, err := v.getFiggyConfigs()
	if err != nil {
		return err
	}
	return v.writeVSCodeSettings(figgyConfigs)
}

func (v *VSCodeHider) getFiggyConfigs() (*[]string, error) {
	figgyConfigs := []string{}
	err := filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == config.FiggyConfigName {
			figgyConfigs = append(figgyConfigs, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return &figgyConfigs, nil
}

func (v *VSCodeHider) readVSCodeSettings() (*map[string]interface{}, error) {
	err := os.MkdirAll(vscodeFolderPath, vscodeSettingFileMode)
	if err != nil {
		return nil, err
	}
	settingsJson := make(map[string]interface{})
	file, err := os.Open(vscodeSettingsPath)
	if err != nil {
		_, err := os.Create(vscodeSettingsPath)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := ioutil.ReadAll(file)
		if err != nil {
			return nil, err
		}
		err = json.Unmarshal(data, &settingsJson)
		if err != nil {
			return nil, err
		}
	}
	return &settingsJson, nil
}

func (v *VSCodeHider) getFilesToHide(figgyConfigs *[]string) (*[]string, error) {
	filesToHide := []string{}
	for _, figgyConfig := range *figgyConfigs {
		path := figgyConfig[:len(figgyConfig)-len(config.FiggyConfigName)]
		foundConfigs, err := config.GetAllFiggiedConfigs(figgyConfig)
		if err != nil {
			return nil, err
		}
		for _, foundConfig := range *foundConfigs {
			filesToHide = append(filesToHide, path+foundConfig)
		}
	}
	return &filesToHide, nil
}

func (v *VSCodeHider) writeVSCodeSettings(figgyConfigs *[]string) error {
	settingsJson, err := v.readVSCodeSettings()
	if err != nil {
		return err
	}
	filesToHide, err := v.getFilesToHide(figgyConfigs)
	if err != nil {
		return err
	}

	excludedFiles, exists := (*settingsJson)[vscodeExcludeFileSetting]
	if !exists {
		excludedFiles = make(map[string]interface{})
		(*settingsJson)[vscodeExcludeFileSetting] = excludedFiles
	}
	for _, fileToHide := range *filesToHide {
		excludedFiles.(map[string]interface{})[fileToHide] = true
	}

	bytes, err := json.MarshalIndent(settingsJson, "", "    ")
	if err != nil {
		return err
	}
	return ioutil.WriteFile(vscodeSettingsPath, bytes, vscodeSettingFileMode)
}
