package task

import (
	"crypto/sha256"
	"fmt"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/fileutils"
	"github.com/xztaityozx/avv/parameters"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestTask_MakeFiles(t *testing.T) {
	as := assert.New(t)

	home, _ := homedir.Dir()
	testDir := filepath.Join(home, "TestDir")
	net := filepath.Join(testDir, "net")
	base := filepath.Join(testDir, "base")
	search := filepath.Join(testDir, "search")

	err := fileutils.TryMakeDirAll(testDir)
	as.NoError(err)

	err = fileutils.TryMakeDirAll(base)
	as.NoError(err)
	err = fileutils.TryMakeDirAll(net)
	as.NoError(err)
	err = fileutils.TryMakeDirAll(search)
	as.NoError(err)

	task := Task{
		AutoRemove: false,
		Parameters: parameters.Parameters{
			Sweeps:    100,
			Vtn:       parameters.Transistor{Threshold: 0.6, Deviation: 1.0, Sigma: 0.046},
			Vtp:       parameters.Transistor{Threshold: -0.6, Deviation: 1.0, Sigma: 0.046},
			Seed:      1,
			ModelFile: filepath.Join(testDir, "ModelFile"),
			AddFile: parameters.AddFile{
				VddVoltage: 0.8,
				GndVoltage: 0.0,
				ICCommand:  "ICCommand",
				Options:    []string{},
				SEED:       1,
			},
			PlotPoint: parameters.PlotPoint{
				Start:   1.11111,
				Step:    2.22222,
				Stop:    3.33333,
				Signals: []string{"A"},
			},
		},
	}

	task.Files, err = parameters.Generate(base, net, task.Parameters)
	as.NoError(err)

	tmp := parameters.Templates{
		SPIScript: filepath.Join(testDir, "SPI"),
	}

	format :=
		`param: vtn=AGAUSS(%.4f,%.4f,%.4f) vtp=AGAUSS(%.4f,%.4f,%.4f)
include: %s
include: %s
monte: %d`
	err = fileutils.WriteFile(tmp.SPIScript, format)

	as.NoError(err)

	err = task.MakeFiles(tmp)
	as.NoError(err)

	hash := fmt.Sprintf("%s%s%s%010d",
		task.PlotPoint.String(), task.Vtn.String(), task.Vtp.String(), task.Sweeps)
	hashWith := fmt.Sprintf("%x", sha256.Sum256([]byte(fmt.Sprintf("%s%s%s%010d%010d",
		task.PlotPoint.String(), task.Vtn.String(), task.Vtp.String(), task.Sweeps, task.Seed))))

	t.Run("AddFile", func(t *testing.T) {
		as.Equal(filepath.Join(base, "sim", hash, fmt.Sprint(task.Seed), "sim", "add"), task.Files.AddFile)
		as.FileExists(task.Files.AddFile)
		b, err := ioutil.ReadFile(task.Files.AddFile)
		as.NoError(err)
		as.Equal([]byte(`VDD VDD! 0 0.8V
VGND GND! 0 0.0V
ICCommand

.option SEED=1`), b)
	})

	t.Run("SPI", func(t *testing.T) {
		as.Equal(filepath.Join(net, hashWith+".spi"), task.Files.SPIScript)
		as.FileExists(task.Files.SPIScript)
		b, err := ioutil.ReadFile(task.Files.SPIScript)
		as.NoError(err)
		as.Equal([]byte(fmt.Sprintf(format,
			task.Vtn.Threshold, task.Vtn.Sigma, task.Vtn.Deviation,
			task.Vtp.Threshold, task.Vtp.Sigma, task.Vtp.Deviation,
			task.Files.AddFile,
			task.ModelFile,
			task.Sweeps,
		)), b)
	})

	t.Run("ACE", func(t *testing.T) {
		as.Equal(filepath.Join(base, "sim", hash, fmt.Sprint(task.Seed), "sim", "ace"), task.Files.ACEScript)
		as.FileExists(task.Files.ACEScript)
		as.Equal(filepath.Join(base, "sim", hash, "result", "00001"), task.Files.ResultFile)
		b, err := ioutil.ReadFile(task.Files.ACEScript)
		as.NoError(err)
		as.Equal([]byte(fmt.Sprintf(`
set xml [ sx_open_wdf "resultsMap.xml" ]
sx_current_sim_file $xml
set www [ sx_signal "A" ]
sx_export_csv on
sx_export_range 1.11111 3.33333 2.22222
sx_export_data "%s" $www
`, task.Files.ResultFile)), b)
	})

	t.Run("XML", func(t *testing.T) {
		as.Equal(filepath.Join(base, "sim", hash, fmt.Sprint(task.Seed), "sim", "resultsMap.xml"), task.Files.ResultsMapXML)
		as.Equal(filepath.Join(base, "sim", hash, fmt.Sprint(task.Seed), "sim", "results.xml"), task.Files.ResultsXML)

		as.FileExists(task.Files.ResultsMapXML)
		as.FileExists(task.Files.ResultsXML)
	})

	err = os.RemoveAll(base)
	as.NoError(err)

}
