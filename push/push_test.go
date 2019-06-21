package push

import (
	"context"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
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
		err := taa.Invoke(context.Background(), parameters.Transistor{}, parameters.Transistor{}, 0, []string{
			gopath, home})

		assert.Error(t, err)
	})

	t.Run("ok", func(t *testing.T) {
		err := taa.Invoke(context.Background(), parameters.Transistor{}, parameters.Transistor{}, 0, []string{
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

	t.Run("fail", func(t *testing.T) {
		taa.ConfigFile = "err"

		err := taa.Invoke(context.Background(), parameters.Transistor{}, parameters.Transistor{}, 0, []string{
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

		assert.Error(t, err)
		assert.Equal(t, fmt.Sprintf("this is error messages: FAILED: 123456"), err.Error())
	})

	t.Run("command not found", func(t *testing.T) {
		taa.TaaPath = "/path/to/taa"
		err := taa.Invoke(context.Background(), parameters.Transistor{}, parameters.Transistor{}, 0, []string{})

		assert.Error(t, err)
	})
}
