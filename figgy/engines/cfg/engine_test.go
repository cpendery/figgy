package cfg

import (
	"testing"

	"github.com/bigkevmcd/go-configparser"
	"github.com/bradleyjkemp/cupaloy"
	"github.com/stretchr/testify/require"
)

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
	expected := "open invalid/: The system cannot find the file specified."

	//WHEN
	_, err := Load(path)

	//THEN
	require.Equal(t, expected, err.Error())
}

func TestWrite(t *testing.T) {
	//GIVEN
	path := "../../../testdataSink/cfg_test_write.cfg"
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
	path := "../../../testdataSink/cfg_test_write.cfg"
	config := make(map[string]interface{})
	section := make(map[string]interface{})
	section["key"] = ""
	config["DEFAULT"] = section

	//WHEN
	err := Write(path, config)

	//THEN
	require.NotNil(t, err)
}
