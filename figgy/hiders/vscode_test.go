package hiders

import (
	"errors"
	"io"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func fixtureGetHider() VSCodeHider {
	return VSCodeHider{
		folderPath:         "../../testdata",
		settingsFile:       "/settings.json",
		settingsPath:       "../../testdata/settings.json",
		excludeFileSetting: "files.exclude",
		fileMode:           444,
		funcWalk:           filepath.Walk,
		funcMkdir:          os.MkdirAll,
		funcOsCreate:       os.Create,
		funcReadAll:        ioutil.ReadAll,
		funcOpen:           os.Open,
	}
}

func TestNewVSCodeHider(t *testing.T) {
	//GIVEN + WHEN
	hider := NewVSCodeHider()

	//THEN
	require.NotNil(t, hider)
}

func TestHide_FailsRead(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	hider.funcWalk = func(s string, wf filepath.WalkFunc) error { return errors.New("") }

	//WHEN
	err := hider.Hide()

	//THEN
	require.NotNil(t, err)
}

func TestHide_FailsWrite(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	hider.settingsPath = "../../testdataSink/fake.txt"
	os.Remove(hider.settingsPath)
	os.MkdirAll("../../testdataSink", hider.fileMode)             //nolint:errcheck
	os.WriteFile(hider.settingsPath, []byte("a"), hider.fileMode) //nolint:errcheck

	//WHEN
	err := hider.Hide()

	//THEN
	require.NotNil(t, err)
}

func TestGetFiggyConfigs(t *testing.T) {
	//GIVEN
	tests := []struct {
		name       string
		configName string
		expected   []string
	}{
		{"no configs found", "", []string{}},
		{"config found", "hider.go", []string{"hider.go"}},
	}
	hider := fixtureGetHider()

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			//WHEN
			hider.configName = test.configName
			paths, err := hider.getFiggyConfigs()

			//THEN
			require.Equal(t, test.expected, *paths)
			require.Nil(t, err)
		})
	}
}

func TestGetFiggyConfigs_Fails(t *testing.T) {
	//GIVEN
	expectedErr := errors.New("Fails Walk")
	hider := fixtureGetHider()
	hider.funcWalk = func(s string, wf filepath.WalkFunc) error {
		return wf("", nil, expectedErr)
	}

	//WHEN
	paths, err := hider.getFiggyConfigs()

	//THEN
	require.Nil(t, paths)
	require.Equal(t, expectedErr, err)
}

func TestReadVSCodeSettings(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	expected := make(map[string]interface{})
	expected["window.zoomLevel"] = float64(-1)

	//WHEN
	data, err := hider.readVSCodeSettings()

	//THEN
	require.Equal(t, expected, *data)
	require.Nil(t, err)
}

func TestReadVSCodeSettings_Fails(t *testing.T) {
	//GIVEN
	tests := []struct {
		name         string
		funcMkdir    funcMkdir
		funcOsCreate funcOsCreate
		funcReadAll  funcReadAll
		funcOpen     funcOpen
	}{
		{"mkdir fails", func(string, fs.FileMode) error { return errors.New("mkdir fails") }, os.Create, ioutil.ReadAll, os.Open},
		{"open & create fails", os.MkdirAll, func(string) (*os.File, error) { return nil, errors.New("create fails") }, ioutil.ReadAll, func(name string) (*os.File, error) { return nil, errors.New("open fails") }},
		{"read all fails", os.MkdirAll, os.Create, func(io.Reader) ([]byte, error) { return nil, errors.New("read fails") }, os.Open},
		{"json parsing fails", os.MkdirAll, os.Create, ioutil.ReadAll, os.Open},
	}
	hider := fixtureGetHider()
	hider.settingsPath = "../../testdata/invalid_settings.json"

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			hider.funcMkdir = test.funcMkdir
			hider.funcOsCreate = test.funcOsCreate
			hider.funcReadAll = test.funcReadAll
			hider.funcOpen = test.funcOpen

			//WHEN
			data, err := hider.readVSCodeSettings()

			//THEN
			require.Nil(t, data)
			require.NotNil(t, err)
		})
	}
}

func TestGetFilesToHide(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	files := []string{"../../testdata/.fig.yaml"}
	expected := []string{"../../testdata/setup.cfg", "../../testdata/tox.cfg"}

	//WHEN
	data, err := hider.getFilesToHide(&files)

	//THEN
	require.Equal(t, expected, *data)
	require.Nil(t, err)
}

func TestGetFilesToHide_Fails(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	files := []string{"../../testdata/invalid.yaml"}

	//WHEN
	data, err := hider.getFilesToHide(&files)

	//THEN
	require.Nil(t, data)
	require.NotNil(t, err)
}

func TestWriteVSCodeSettings(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	hider.settingsPath = "../../testdataSink/fakeFile.txt"
	files := []string{"../../testdata/.fig.yaml"}
	os.MkdirAll("../../testdataSink", hider.fileMode) //nolint:errcheck
	os.Remove(hider.settingsPath)

	//WHEN
	err := hider.writeVSCodeSettings(&files)

	//THEN
	require.Nil(t, err)
}

func TestWriteVSCodeSettings_FailsToReadSettings(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	hider.funcMkdir = func(path string, perm fs.FileMode) error { return errors.New("mkdir fails") }
	files := []string{"../../testdata/invalid.yaml"}

	//WHEN
	err := hider.writeVSCodeSettings(&files)

	//THEN
	require.NotNil(t, err)
}

func TestWriteVSCodeSettings_FailsToFindFilesToFind(t *testing.T) {
	//GIVEN
	hider := fixtureGetHider()
	files := []string{"../../testdata/invalid.yaml"}

	//WHEN
	err := hider.writeVSCodeSettings(&files)

	//THEN
	require.NotNil(t, err)
}
