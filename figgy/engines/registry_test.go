package engines

import (
	"testing"

	"github.com/cpendery/figgy/figgy/models"
	"github.com/stretchr/testify/require"
)

func TestRegister(t *testing.T) {
	//GIVEN
	ext := ".answer"

	//WHEN
	Register(models.ConfigEngine{
		Extension: ext,
	})

	//THEN
	_, exists := allConfigIOs[ext]
	require.True(t, exists)
}

func TestRegister_Fails(t *testing.T) {
	//GIVEN
	ext := ".answer"
	allConfigIOs[ext] = models.ConfigEngine{Extension: ext}

	defer func() {
		//THEN
		r := recover()
		require.NotNil(t, r)
	}()

	//WHEN
	Register(models.ConfigEngine{
		Extension: ext,
	})
}

func TestGet(t *testing.T) {
	//GIVEN
	ext := ".answer"
	allConfigIOs[ext] = models.ConfigEngine{Extension: ext}

	//WHEN
	engine, exists := Get(ext)

	//THEN
	require.NotNil(t, engine)
	require.True(t, exists)
}
