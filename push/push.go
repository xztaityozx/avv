package push

import (
	"bufio"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
	"github.com/xztaityozx/avv/parameters"
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

type TaaResultKey struct {
	Vtn    parameters.Transistor
	Vtp    parameters.Transistor
	Sweeps int
}

func (t TaaResultKey) Hash() string {
	return fmt.Sprint(sha256.Sum256([]byte(fmt.Sprintf("%s%s%d", t.Vtn.String(),t.Vtp.String(),t.Sweeps))))
}

// getCommand generate command for pushing data to database with taa command
// returns:
//  - string: command string
func (t Taa) getCommand(vtn, vtp parameters.Transistor, sweeps int, files []string) string {
	return fmt.Sprintf("%s push --config %s --parallel %d --VtpVoltage %f --vtpSigma %f --vtpDeviation %f --VtnVoltage %f --vtnSigma %f --vtnDeviation %f --sweeps %d %s",
		t.TaaPath, t.ConfigFile, t.Parallel,
		vtp.Threshold, vtp.Sigma, vtp.Deviation,
		vtn.Threshold, vtn.Sigma, vtn.Deviation,
		sweeps,
		strings.Join(files, " "))
}

// Invoke start push command
// params:
//  - ctx: context
//  - files: files of data, these filename must be `SEED%05d.csv` format
// returns:
//  - error: error
func (t Taa) Invoke(ctx context.Context, vtn, vtp parameters.Transistor, sweeps int, files []string) error {
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


	ch := make(chan error)
	act := func() error {
		cmd := exec.Command("bash", "-c", t.getCommand(vtn, vtp, sweeps, files))

		stdout, err := cmd.StdoutPipe()
		if err != nil {
			ch <- err
		}
		stderr, err := cmd.StderrPipe()
		if err != nil {
			ch<-err
		}

		err = cmd.Start()
		if err != nil {
			ch <- err
		}

		box := [2]error{nil, nil}

		go func() {
			scan := bufio.NewScanner(stdout)
			for scan.Scan() {
				bar.IncrBy(len(strings.Split(scan.Text(), "\n")))
			}

		}()

		go func() {
			scan := bufio.NewScanner(stderr)
			agg := ""
			for scan.Scan() {
				agg+=scan.Text()
			}
			ch <- errors.New(agg)
		}()


		p.Wait()

	}



	select {
	case <-ctx.Done():
		return errors.New("canceled")
	case err := <-ch:
		return err

	}
}
