package push

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"os"
	"path/filepath"
	"testing"
)

func TestTaa_Invoke(t *testing.T) {
	gohome := os.Getenv("GOPATH")
	taa := Taa{
		TaaPath: filepath.Join(gohome, "src", "github.com", "xztaityozx", "avv", "test", "taa.exe"),
	}

	x := task.Task{
		Files: parameters.Files{
			ResultFile: "/path/to/SEED00001.csv",
		},
	}

	y := task.Task{
		Files: parameters.Files{
			ResultFile: "/path/to/XXXX.csv",
		},
	}

	as := assert.New(t)

	ctx := context.Background()

	t.Run("unexpected", func(t *testing.T) {
		_, err := taa.Invoke(ctx, y)
		as.Error(err)
	})

	t.Run("fail", func(t *testing.T) {
		taa.ConfigFile = "err"
		_, err := taa.Invoke(ctx, x)
		as.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		taa.ConfigFile = "/path/to/config"
		_, err := taa.Invoke(ctx, x)
		as.NoError(err)
	})
}
