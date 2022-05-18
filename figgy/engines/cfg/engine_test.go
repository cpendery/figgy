package cfg

import (
	"os"
	"testing"

	"github.com/bigkevmcd/go-configparser"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"
)

const testSinkDirPath = "../../../testdataSink"

func TestLoad(t *testing.T) {
	//GIVEN
	path := "../../../testdata/config.cfg"

	//WHEN
	res, err := Load(path)

	//THEN
	cupaloy.SnapshotT(t, res)
	require.Nil(t, err)
}

func TestLoad_Fails(t *testing.T) {
	//GIVEN
	path := "invalid/"

	//WHEN
	_, err := Load(path)

	//THEN
	require.NotNil(t, err)
}

func TestWrite(t *testing.T) {
	//GIVEN
	path := testSinkDirPath + "/cfg_test_write.cfg"
	os.MkdirAll(testSinkDirPath, 0444) //nolint:errcheck
	config := make(map[string]interface{})
	section := make(map[string]interface{})
	section["key"] = "value"
	config["section"] = section

	//WHEN
	err := Write(path, config)

	//THEN
	require.Nil(t, err)
	res, err := Load(path)
	require.Nil(t, err)

	expected_section := config["section"]
	expected_value := expected_section.(map[string]interface{})["key"].(string)
	found_section := res["section"]
	found_value := found_section.(configparser.Dict)["key"]
	require.Equal(t, expected_value, found_value)
}

func TestWrite_Fails(t *testing.T) {
	//GIVEN
	path := testSinkDirPath + "/cfg_test_write_fails.cfg"
	os.MkdirAll(testSinkDirPath, 0444) //nolint:errcheck
	config := make(map[string]interface{})
	section := make(map[string]interface{})
	section["key"] = ""
	config["DEFAULT"] = section

	//WHEN
	err := Write(path, config)

	//THEN
	require.NotNil(t, err)
}
