package simulation

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/fileutils"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"os"
	"path/filepath"
	"testing"
)

func TestHSPICE_Invoke(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	home,_:=homedir.Dir()

	t.Run("no error", func(t *testing.T) {
		w := HSPICE{Path: filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "hspice.sh"), Options:""}

		fileutils.TryMakeDirAll(filepath.Join(home, "TestDir"))

		err := w.Invoke(context.Background(), task.Task{
			Files:parameters.Files{
				SPIScript:"/path/to/spi",
				Directories:parameters.Directories{
					DstDir:filepath.Join(home, "TestDir"),
				},
			},
		})

		assert.NoError(t, err)

	})

	t.Run("error", func(t *testing.T) {
		w := HSPICE{Path: filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "hspice.sh"), Options:"err"}

		fileutils.TryMakeDirAll(filepath.Join(home, "TestDir"))

		err := w.Invoke(context.Background(), task.Task{
			Files:parameters.Files{
				SPIScript:"/path/to/spi",
				Directories:parameters.Directories{
					DstDir:filepath.Join(home, "TestDir"),
				},
			},
		})

		assert.Error(t, err)
	})
}
