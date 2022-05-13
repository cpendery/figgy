package engines

import (
	"fmt"

	"github.com/cpendery/figgy/figgy/models"
)

var allConfigIOs = make(map[string]models.ConfigEngine)

// Register registers a template with the given name.
// Intended to be called at program init time.
func Register(t models.ConfigEngine) {
	if _, exists := allConfigIOs[t.Extension]; exists {
		panic(fmt.Sprintf("duplicate template: %v", t.Extension))
	}
	allConfigIOs[t.Extension] = t
}

func Get(name string) (models.ConfigEngine, bool) {
	t, exists := allConfigIOs[name]
	return t, exists
}
