package push

import (
	"context"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/xztaityozx/avv/parameters"
	"github.com/xztaityozx/avv/task"
	"golang.org/x/xerrors"
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
	// logrus.Logger
	Log *logrus.Logger
}

// getCommand generate command for pushing data to database with taa command
// returns:
//  - string: command string
func (taa Taa) getCommand(vtn, vtp parameters.Transistor, sweeps int, file string) string {
	return fmt.Sprintf("%s push --config %s --parallel %d --VtpVoltage %f --vtpSigma %f --vtpDeviation %f --VtnVoltage %f --vtnSigma %f --vtnDeviation %f --sweeps %d %s",
		taa.TaaPath, taa.ConfigFile, taa.Parallel,
		vtp.Threshold, vtp.Sigma, vtp.Deviation,
		vtn.Threshold, vtn.Sigma, vtn.Deviation,
		sweeps,
		file)
}

// Invoke start push command
// params:
//  - ctx: context
//  - files: files of data, these filename must be `SEED%05d.csv` format
// returns:
//  - error: error
func (taa Taa) Invoke(ctx context.Context, t task.Task) (task.Task, error) {
	// check filename
	base := filepath.Base(t.Files.ResultFile)
	unexpected := errors.New(fmt.Sprintf("Unexpected filename: %s", base))

	if len(base) < len("SEED00000.csv") {
		return t, unexpected
	}

	if _, err := strconv.Atoi(
		strings.Replace(
			strings.Replace(base, "SEED", "", -1),
			".csv", "", -1)); err != nil {
		return t, unexpected
	}

	command := taa.getCommand(t.Vtn, t.Vtp, t.Sweeps, t.Files.ResultFile)
	taa.Log.Info(command)
	_, err := exec.Command("bash", "-c", command).Output()

	if err != nil {
		return t, xerrors.Errorf("failed taa push: %s", err)
	}

	return t, nil
}
