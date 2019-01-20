package cmd

import (
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"testing"
)

var TU TestUtils

type TestUtils struct {
	Initialized bool
}

// Initialized TestUtils struct
func (tu TestUtils) Init() {

	if TU.Initialized {
		return
	}

	logrus.SetFormatter(&logrus.TextFormatter{
		ForceColors: true,
	})
	logrus.SetOutput(colorable.NewColorableStdout())
	TU.Initialized = true
}

// Compare Test
func (tu TestUtils) Equal(actual, expect interface{}, t *testing.T) {
	if actual != expect {
		t.Fatal(actual, "is not", expect, "\n(at: ", t.Name(), ")\nactual: ", actual, "\nexpect: ", expect)
	}
}

// Compare Test with comparer func
func (tu TestUtils) EqualComparer(actual, expect interface{}, t *testing.T, comparer func(interface{}, interface{}) bool) {
	if !comparer(actual, expect) {
		t.Fatal(actual, "is not", expect, "\n(at: ", t.Name(), ")\nactual: ", actual, "\nexpect: ", expect)
	}
}

// Check Boolean Status Test
func (tu TestUtils) Assert(status bool, t *testing.T, msg ...interface{}) {
	if !status {
		t.Fatal(t.Name(), "had failed", msg)
	}
}

// Compare Collection Test
func (tu TestUtils) CompareCollection(actual, expect []interface{}, t *testing.T, comparer func(interface{}, interface{}) bool) {
	if len(actual) != len(expect) {
		t.Fatal("does not match length of collections\nlen(actual): ", len(actual), "\nlen(expect): ", len(expect))
	}

	for i, v := range actual {
		tu.EqualComparer(v, expect[i], t, comparer)
	}
}
