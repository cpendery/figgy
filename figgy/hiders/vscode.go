package hiders

import (
	"encoding/json"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/cpendery/figgy/figgy/config"
)

type VSCodeHider struct {
	folderPath         string
	settingsFile       string
	settingsPath       string
	excludeFileSetting string
	fileMode           fs.FileMode
	configName         string

	funcWalk     funcWalk
	funcMkdir    funcMkdir
	funcOsCreate funcOsCreate
	funcReadAll  funcReadAll
	funcOpen     funcOpen
}

type funcWalk func(string, filepath.WalkFunc) error
type funcMkdir func(string, fs.FileMode) error
type funcOsCreate func(string) (*os.File, error)
type funcReadAll func(io.Reader) ([]byte, error)
type funcOpen func(name string) (*os.File, error)

func NewVSCodeHider() Hider {
	return &VSCodeHider{
		folderPath:         "./.vscode",
		settingsFile:       "/settings.json",
		settingsPath:       "./.vscode/settings.json",
		excludeFileSetting: "files.exclude",
		fileMode:           0444,

		configName:   config.FiggyConfigName,
		funcWalk:     filepath.Walk,
		funcMkdir:    os.MkdirAll,
		funcOsCreate: os.Create,
		funcReadAll:  ioutil.ReadAll,
		funcOpen:     os.Open,
	}
}

//Updates a .vscode settings file to hide all files found in
//any figgy files at any depth from the current directory
func (v *VSCodeHider) Hide() error {
	figgyConfigs, err := v.getFiggyConfigs()
	if err != nil {
		return err
	}
	return v.writeVSCodeSettings(figgyConfigs)
}

func (v *VSCodeHider) getFiggyConfigs() (*[]string, error) {
	figgyConfigs := []string{}
	err := v.funcWalk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && info.Name() == v.configName {
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
	err := v.funcMkdir(v.folderPath, v.fileMode)
	if err != nil {
		return nil, err
	}
	settingsJson := make(map[string]interface{})
	file, err := v.funcOpen(v.settingsPath)
	if err != nil {
		_, err := v.funcOsCreate(v.settingsPath)
		if err != nil {
			return nil, err
		}
	} else {
		data, err := v.funcReadAll(file)
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

	excludedFiles, exists := (*settingsJson)[v.excludeFileSetting]
	if !exists {
		excludedFiles = make(map[string]interface{})
		(*settingsJson)[v.excludeFileSetting] = excludedFiles
	}
	for _, fileToHide := range *filesToHide {
		excludedFiles.(map[string]interface{})[fileToHide] = true
	}

	bytes, _ := json.MarshalIndent(settingsJson, "", "    ")
	return ioutil.WriteFile(v.settingsPath, bytes, v.fileMode)
}
