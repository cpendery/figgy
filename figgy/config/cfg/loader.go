package cfg

import (
	configparser "github.com/bigkevmcd/go-configparser"
)

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
