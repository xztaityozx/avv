package remove

import (
	"context"
	"errors"
	"github.com/xztaityozx/avv/task"
	"os"
	"path/filepath"
)

type Remove struct{}

func (r Remove) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	Do(ctx, t.Files.SPIScript)
	Do(ctx, t.Files.TaskFile)
	Do(ctx, filepath.Join(t.Files.Directories.DstDir, "../"))
	return t, nil
}

// Do remove files
// params:
//  - ctx: context
// returns:
//	- error:
func Do(ctx context.Context, path string) error {
	ch := make(chan error, 1)
	defer close(ch)

	go func() {
		ch <- os.RemoveAll(path)
	}()

	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-ch:
		return err
	}
}
