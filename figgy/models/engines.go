package models

type ParsedConfig = map[string]interface{}

type Load func(filename string) (ParsedConfig, error)
type Write func(path string, parsedConfig ParsedConfig) error

type ConfigEngine struct {
	Extension string
	Load      Load
	Write     Write
}
