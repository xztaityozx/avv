package push

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"
)

type Taa struct {
	// Path to taa.exe
	TaaPath string
	// path to config.yml for taa.exe
	ConfigFile string
	// number of parallel
	Parallel int
}

// getCommand generate command for pushing data to database with taa command
// returns:
//  - string: command string
func (t Taa) getCommand(files []string) string {
	return fmt.Sprintf("%s push --config %s --parallel %d %s", t.TaaPath, t.ConfigFile,t.Parallel,
			strings.Join(files," "))
}

// Invoke start push command
// params:
//  - ctx: context
//  - files: files of data, these filename must be `SEED%05d.csv` format
// returns:
//  - error: error
func (t Taa) Invoke(ctx context.Context, files []string) error {
	// check filename
	for _,v := range files {
		base := filepath.Base(v)
		unexpected := errors.New(fmt.Sprintf("Unexpected filename: %s", base))

		if len(base) < len("SEED00000.csv") {
			return unexpected
		}

		if _, err := strconv.Atoi(
			strings.Replace(
				strings.Replace(base, "SEED","",-1),
				".csv","",-1)); err != nil {
			return unexpected
		}
	}

	ch := make(chan error)

	select {
	case <-ctx.Done():
		return errors.New("Canceled\n")
	case err := <-ch:
		return err

	}
}