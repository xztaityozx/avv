package push

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"os/exec"
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
	return fmt.Sprintf("%s push --config %s --parallel %d %s", t.TaaPath, t.ConfigFile, t.Parallel,
		strings.Join(files, " "))
}

// Invoke start push command
// params:
//  - ctx: context
//  - files: files of data, these filename must be `SEED%05d.csv` format
// returns:
//  - error: error
func (t Taa) Invoke(ctx context.Context, files []string) error {
	// check filename
	for _, v := range files {
		base := filepath.Base(v)
		unexpected := errors.New(fmt.Sprintf("Unexpected filename: %s", base))

		if len(base) < len("SEED00000.csv") {
			return unexpected
		}

		if _, err := strconv.Atoi(
			strings.Replace(
				strings.Replace(base, "SEED", "", -1),
				".csv", "", -1)); err != nil {
			return unexpected
		}
	}

	ch := make(chan error)
	p := mpb.NewWithContext(ctx)
	total := len(files)
	name := color.New(color.FgHiYellow).Sprint("push")
	inProgressMSG := color.New(color.FgCyan).Sprint("Processing...")
	finishMSG := color.New(color.FgHiGreen).Sprint("done!")
	bar := p.AddBar(int64(total),
		mpb.BarStyle("|██▒|"),
		mpb.BarWidth(50),
		mpb.PrependDecorators(
			decor.Name(name, decor.WC{W: len(name) + 1, C: decor.DidentRight}),
			decor.OnComplete(decor.EwmaETA(decor.ET_STYLE_GO, 60, decor.WC{W: 4}), "done")),
		mpb.AppendDecorators(
			decor.Name("   "),
			decor.Percentage(),
			decor.Name(" | "),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			decor.Name(" | "),
			decor.OnComplete(decor.Name(inProgressMSG), finishMSG)),
	)

	go func() {
		cmd := exec.Command("bash", "-c", t.getCommand(files))

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- err
		}

		err = cmd.Start()
		if err != nil {
			ch <- err
		}

		scan := bufio.NewScanner(stdout)
		for scan.Scan() {
			bar.IncrBy(len(strings.Split(scan.Text(), "\n")))
		}

		p.Wait()
		ch <- nil
	}()

	select {
	case <-ctx.Done():
		return errors.New("Canceled\n")
	case err := <-ch:
		return err

	}
}
