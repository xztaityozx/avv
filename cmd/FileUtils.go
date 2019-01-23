package cmd

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"unsafe"
)

type FileUtils struct {
	logger *logrus.Entry
}

var FU = FileUtils{
	logger: log.WithField("name", "FileUtils"),
}

// Cat(p string)
// Read strings from [p]
// returns: [strings]
func (fu FileUtils) Cat(p string) string {
	b, err := ioutil.ReadFile(p)
	if err != nil {
		fu.logger.Error("can not open", p)
		fu.logger.Fatal(err)
	}

	return *(*string)(unsafe.Pointer(&b))
}

// WriteFile(p string)
// Wrapper for ioutil.WriteFile
func (fu FileUtils) WriteFile(p string, data string) {
	if _, err := os.Create(p); err != nil {
		fu.logger.Fatal(err)
	}

	if err := ioutil.WriteFile(p, []byte(data), 0644); err != nil {
		fu.logger.Fatal(err)
	}
}

// TryMkdir(path string)
// try make directory path to [path]
func (fu FileUtils) TryMkDir(p string) {
	if _, err := os.Stat(p); err != nil {

		// output warn
		fu.logger.Warn(err)
		// mkdir all
		if err := os.MkdirAll(p, 0644); err != nil {
			fu.logger.Fatal(err)
		} else {
			fu.logger.Info(p, " had created")
		}
	}
}

// wrapper for filepath.Join
func PathJoin(p ...string) string {
	return filepath.Join(p...)
}

// wrapper for os.ChDir
func (fu FileUtils) TryChDir(p string) {
	if _, err := os.Stat(p); err != nil {
		log.Warn(err)

		// mkdir
		fu.TryMkDir(p)
	}

	if err := os.Chdir(p); err != nil {
		fu.logger.Fatal(err)
	}
}

// wrapper for io.Copy
func (fu FileUtils) Copy(src, dst string) {
	sfp, err := os.Open(src)
	if err != nil {
		fu.logger.Fatal(err)
	}
	dfp, err := os.Open(dst)
	if err != nil {
		fu.logger.Fatal(err)
	}

	if _, err := io.Copy(dfp, sfp); err != nil {
		fu.logger.Fatal(err)
	}

	fu.logger.Info("Copy: ", src, "copy to", dst)
}
