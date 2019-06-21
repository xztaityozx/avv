package remove

import (
	"context"
	"github.com/mitchellh/go-homedir"
	"github.com/stretchr/testify/assert"
	"github.com/xztaityozx/avv/fileutils"
	"os"
	"path/filepath"
	"testing"
)

func TestRemoveSimulationFiles(t *testing.T) {
	home, _ := homedir.Dir()
	path := filepath.Join(home, "TestDir")
	fileutils.TryMakeDirAll(path)

	path = filepath.Join(path, "A")

	_, err := os.Create(path)

	assert.NoError(t, err)

	err = Do(context.Background(), path)
	assert.NoError(t, err)
}
