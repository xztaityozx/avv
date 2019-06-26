package simulation

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"os"
	"path/filepath"
	"testing"
)

func TestHSPICE_Invoke(t *testing.T) {
	gopath := os.Getenv("GOPATH")
	h := HSPICE{
		Path: filepath.Join(gopath, "src", "github.com", "xztaityozx", "avv", "test", "hspice.sh"),
	}

	x := task.Task{
		Parameters: parameters.Parameters{
			Seed: 1,
		},
	}

	as := assert.New(t)

	t.Run("fail", func(t *testing.T) {
		h.Options = "err"
		_, err := h.Invoke(context.Background(), x)
		as.Error(err)
	})

	t.Run("success", func(t *testing.T) {
		h.Options = ""
		_, err := h.Invoke(context.Background(), x)
		as.NoError(err)
	})
}
