package extract

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

func TestTask_GetCommand(t *testing.T) {
	w := WaveView{Path: "/path/to/wv"}

	actual := w.getCommand("/path/to/dst", "/path/to/ace")
	expect := "cd /path/to/dst && /path/to/wv -k -ace_no_gui /path/to/ace &> ./wv.log"

	assert.Equal(t, expect, actual)
}

func TestWaveView_Invoke(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	home, _ := homedir.Dir()

	t.Run("no error", func(t *testing.T) {
		w := WaveView{Path: filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "wv.sh")}

		fileutils.TryMakeDirAll(filepath.Join(home, "TestDir"))

		err := w.Invoke(context.Background(), task.Task{
			Files: parameters.Files{
				ACEScript: "/path/to/ace",
				Directories: parameters.Directories{
					DstDir: filepath.Join(home, "TestDir"),
				},
			},
		})

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		w := WaveView{Path: filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "wv.err.sh")}

		fileutils.TryMakeDirAll(filepath.Join(home, "TestDir"))

		err := w.Invoke(context.Background(), task.Task{
			Files: parameters.Files{
				ACEScript: "/path/to/ace",
				Directories: parameters.Directories{
					DstDir: filepath.Join(home, "TestDir"),
				},
			},
		})

		assert.Error(t, err)

	})

}
