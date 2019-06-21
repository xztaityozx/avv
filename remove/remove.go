package remove

import (
	"context"
	"errors"
	"os"
)

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
