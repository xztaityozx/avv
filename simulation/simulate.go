package simulation

import (
	"context"
	"fmt"
	"github.com/xztaityozx/avv/task"
	"os/exec"
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
func (h HSPICE) Invoke(ctx context.Context, task task.Task) error {

	ch := make(chan error)

	go func() {
		_, err := exec.Command("bash", "-c",
			h.getCommand(task.Files.Directories.DstDir, task.Files.SPIScript)).Output()
		ch <- err
	}()

	select {
	case <-ctx.Done():
		return nil
	case err := <-ch:
		return err

	}
}
