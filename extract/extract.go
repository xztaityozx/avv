package extract

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"os/exec"

	"github.com/xztaityozx/avv/task"
	"golang.org/x/xerrors"
)

type WaveView struct {
	Path string
	Log  *logrus.Logger
}

// getCommand generate command string for extracting
// returns:
//  - string: command string
func (w WaveView) getCommand(dst, ace string) string {
	return fmt.Sprintf("cd %s && %s -k -ace_no_gui %s &> ./wv.log",
		dst,
		w.Path, ace)
}

// Invoke start extract task with custom waveview
func (w WaveView) Invoke(ctx context.Context, task task.Task) (task.Task, error) {
	ch := make(chan error, 1)

	command := w.getCommand(task.Files.Directories.DstDir, task.Files.ACEScript)
	w.Log.Info(command)
	go func() {
		defer close(ch)
		_, err := exec.Command("bash", "-c", command).Output()

		if err != nil {
			ch <- xerrors.Errorf("failed wv: %s", err)
		} else {
			//ch <- remove.Do(ctx, task.Files.Directories.DstDir)
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
