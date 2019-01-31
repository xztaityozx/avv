package cmd

import (
	"github.com/mitchellh/go-homedir"
	"testing"
)

func TestFilter_CountUp(t *testing.T) {
	home, _ := homedir.Dir()
	src := PathJoin(home, "TestDir")
	FU.TryMkDir(src)
	src = PathJoin(src, "src.csv")
	FU.WriteFile(src, `
10, 21, 32, 43
51, 62, 83, 94
`)

}
