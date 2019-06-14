package parameters

import (
	"github.com/xztaityozx/avv/logger"
	"github.com/xztaityozx/avv/simulation"
	"path/filepath"
)

type Config struct {
	Templates   simulation.Templates
	Default     Default
	SlackConfig logger.SlackConfig
	AutoRemove  bool
}

type Default struct {
	SearchDir  string
	NetListDir string
	BaseDir    string
	SEED       SEED
	Parameters Parameters
}

func (d Default) TaskDir() string {
	return filepath.Join(d.BaseDir, "task")
}
