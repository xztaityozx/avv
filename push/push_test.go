package push

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

func TestTaa_Invoke(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	home, _ := homedir.Dir()
	taa := Taa{
		TaaPath:    filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "taa.exe"),
		Parallel:   10,
		ConfigFile: "/path/to/config",
	}

	t.Run("unexpected file path", func(t *testing.T) {
		err := taa.Invoke(context.Background(), []string{
			gopath, home})

		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		err := taa.Invoke(context.Background(), []string{
			"/path/to/SEED00001.csv",
			"/path/to/SEED00002.csv",
			"/path/to/SEED00003.csv",
			"/path/to/SEED00004.csv",
			"/path/to/SEED00005.csv",
			"/path/to/SEED00006.csv",
			"/path/to/SEED00007.csv",
			"/path/to/SEED00008.csv",
			"/path/to/SEED00009.csv",
			"/path/to/SEED00010.csv",
		})

		assert.NoError(t, err)
	})
}
