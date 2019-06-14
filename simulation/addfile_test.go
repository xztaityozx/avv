package simulation

import (
	"bytes"
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/fileutils"
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"
)

func TestAddFile_GenerateAddFile(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, "TestDir")
	fileutils.TryMakeDirAll(path)

	path = filepath.Join(path, "addfile")
	a := AddFile{
		SEED:1,
		GndVoltage:0.0,
		ICCommand:"ic command",
		Options:[]string{"option1", "option2"},
		VddVoltage:1.2,
	}

	err := a.GenerateAddFile(path)
	as := assert.New(t)

	as.NoError(err)

	data,err := ioutil.ReadFile(path)
	as.NoError(err)

	expect := `VDD VDD! 0 1.2V
VGND GND! 0 0.0V
ic command
option1
option2
.option SEED=1`

	as.True(bytes.Equal(data, []byte(expect)))

	err = os.Remove(path)
	if err != nil {
		logrus.Warn(err)
	}
}

