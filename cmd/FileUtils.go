package cmd

import (
	"github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
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
		fu.logger.WithField("at", "Cat").Error("can not open", p)
		fu.logger.Fatal(err)
	}

	return *(*string)(unsafe.Pointer(&b))
}

// WriteFile(p string)
// Wrapper for ioutil.WriteFile
func (fu FileUtils) WriteFile(p string, data string) {
	if err := ioutil.WriteFile(p, []byte(data), 0644); err != nil {
		fu.logger.WithField("at", "WriteFile").Fatal(err)
	}
}

// TryMkdir(path string)
// try make directory path to [path]
func (fu FileUtils) TryMkDir(p string) {
	if _, err := os.Stat(p); err != nil {

		// output warn
		fu.logger.Warn(err)
		// mkdir all
		if err := os.MkdirAll(p, 0755); err != nil {
			fu.logger.WithField("at", "TryMkDir").Fatal(err)
		} else {
			fu.logger.WithField("at", "TryMkDir").Info(p, " had created")
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
		log.WithField("at", "TryChDir").Warn(err)

		// mkdir
		fu.TryMkDir(p)
	}

	if err := os.Chdir(p); err != nil {
		fu.logger.WithField("at", "TryChDir").Fatal(err)
	}
}

// wrapper for io.Copy
func (fu FileUtils) Copy(src, dst string) {
	sfp, err := os.Open(src)
	if err != nil {
		fu.logger.WithField("at", "Copy").Fatal(err)
	}
	dfp, err := os.Open(dst)
	if err != nil {
		fu.logger.WithField("at", "Copy").Fatal(err)
	}

	if _, err := io.Copy(dfp, sfp); err != nil {
		fu.logger.WithField("at", "Copy").Fatal(err)
	}

	fu.logger.WithField("at", "Copy").Info("Copy: ", src, "copy to", dst)
}

func (fu FileUtils) WriteSlice(dst string, box []string, sep string) {
	fu.WriteFile(dst, strings.Join(box, sep))
}

func (fu FileUtils) Ls(path string) []os.FileInfo {
	rt, err := ioutil.ReadDir(path)
	if err != nil {
		fu.logger.Fatal(err)
	}
	return rt
}
