package simulation

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/avv/parameters"
	"os/exec"

	"github.com/xztaityozx/avv/task"
	"golang.org/x/xerrors"
)

type HSPICE struct {
	Path    string
	Options string
	Tmp     parameters.Templates
	Log     *logrus.Logger
}

// getCommand generate simulation command with hspice
// returns:
//  - string: command string
func (h HSPICE) getCommand(dst, spi string) string {
	return fmt.Sprintf("cd %s && %s %s -i %s -o ./hspice &> ./hspice.log",
		dst,
		h.Path, h.Options,
		spi)
}

// Invoke start simulation with context
func (h HSPICE) Invoke(ctx context.Context, task task.Task) (task.Task, error) {

	ch := make(chan error)

	go func() {
		task.MakeFiles(h.Tmp)
		command := h.getCommand(task.Files.Directories.DstDir, task.Files.SPIScript)
		h.Log.Info(command)
		_, err := exec.Command("bash", "-c", command).Output()
		if err != nil {
			ch <- xerrors.Errorf("simulation command failed: %s", err)
		} else {
			ch <- nil
		}
	}()

	select {
	case <-ctx.Done():
		return task, errors.New("canceled")
	case err := <-ch:
		return task, err

	}
}
