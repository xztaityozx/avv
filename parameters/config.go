package parameters

import (
	"github.com/xztaityozx/avv/logger"
	"path/filepath"
)

type Config struct {
	Templates   Templates
	Default     Default
	SlackConfig logger.SlackConfig
	AutoRemove  bool
	HSPICE      HSPICEConfig
	WaveView    WaveViewConfig
	Taa         TaaConfig
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

type HSPICEConfig struct {
	Path    string
	Options string
}

type TaaConfig struct {
	Path       string
	ConfigFile string
	Parallel   int
}

type WaveViewConfig struct {
	Path string
}
