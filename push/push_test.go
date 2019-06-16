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
		TaaPath: filepath.Join(gopath, "src","github.com","xztaityozx","avv","test","taa.exe"),
		Parallel:10,
		ConfigFile:"/path/to/config",
	}

	t.Run("unexpected file path", func(t *testing.T) {
		err := taa.Invoke(context.Background(), []string{
			gopath, home})

		assert.Error(t, err)
	})
}
