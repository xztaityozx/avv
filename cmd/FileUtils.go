package cmd

import (
	"github.com/sirupsen/logrus"
	"io"
	"os"
	"path/filepath"
)

// TryMkdir(path string)
// try make directory path to [path]
func TryMkDir(p string) {
	if _, err := os.Stat(p); err != nil {

		// output warn
		logrus.Warn(err)
		// mkdir all
		if err := os.MkdirAll(p, 0644); err != nil {
			logrus.Fatal(err)
		} else {
			logrus.Println(p, " had created")
		}
	}
}

// wrapper for filepath.Join
func PathJoin(p ...string) string {
	return filepath.Join(p...)
}

// wrapper for os.ChDir
func TryChDir(p string) {
	if _, err := os.Stat(p); err != nil {
		logrus.Warn(err)

		// mkdir
		TryMkDir(p)
	}

	if err := os.Chdir(p); err != nil {
		logrus.Fatal(err)
	}
}

// wrapper for io.Copy
func Copy(src, dst string) {
	sfp, err := os.Open(src)
	if err != nil {
		logrus.Fatal(err)
	}
	dfp, err := os.Open(dst)
	if err != nil {
		logrus.Fatal(err)
	}

	if _, err := io.Copy(dfp, sfp); err != nil {
		logrus.Fatal(err)
	}

	logrus.Println("Copy: ", src, "copy to", dst)
}
