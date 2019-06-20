package extract

import (
	"context"
	"errors"
	"fmt"
	"github.com/xztaityozx/avv/task"
	"os/exec"
)

type WaveView struct {
	Path string
}


// getCommand generate command string for extracting
// returns:
//  - string: command string
func (w WaveView) getCommand(dst,ace string) string {
	return fmt.Sprintf("cd %s && %s -k -ace_no_gui %s &> ./wv.log",
		dst,
		w.Path, ace)
}

func (w WaveView) Invoke(ctx context.Context, task task.Task) error {
	ch := make(chan error)

	go func() {
		_, err := exec.Command("bash","-c",
			w.getCommand(task.Files.Directories.DstDir, task.Files.ACEScript)).Output()
		ch<-err
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err :=<-ch:
		return err
	}
}
