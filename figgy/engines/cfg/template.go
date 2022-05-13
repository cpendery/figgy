package cfg

import (
	"github.com/cpendery/figgy/figgy/engines"
	"github.com/cpendery/figgy/figgy/models"

	configparser "github.com/bigkevmcd/go-configparser"
)

func init() {
	engines.Register(models.ConfigEngine{
		Extension: ".cfg",
		Load:      Load,
		Write:     Write,
	})
}

func Load(filename string) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	p, err := configparser.NewConfigParserFromFile(filename)
	if err != nil {
		return nil, err
	}

	for _, section := range p.Sections() {
		res, err := p.Items(section)
		if err != nil {
			return nil, err
		}
		result[section] = res
	}
	return result, nil
}

func Write(path string, fig map[string]interface{}) error {
	p := configparser.New()
	for sectionHeader, sectionData := range fig {
		if err := p.AddSection(sectionHeader); err != nil {
			return err
		}
		for option, value := range sectionData.(map[string]interface{}) {
			if err := p.Set(sectionHeader, option, value.(string)); err != nil {
				return err
			}
		}
	}
	return p.SaveWithDelimiter(path, "=")
}
