package simulation

import (
	"context"
	"errors"
	"fmt"
	"os/exec"

	"github.com/xztaityozx/avv/task"
	"golang.org/x/xerrors"
)

type HSPICE struct {
	Path    string
	Options string
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
		command := h.getCommand(task.Files.Directories.DstDir, task.Files.SPIScript)
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
