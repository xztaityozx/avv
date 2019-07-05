package parameters

import (
	"path/filepath"

	"github.com/xztaityozx/avv/logger"
)

type Config struct {
	Templates   Templates
	Default     Default
	SlackConfig logger.SlackConfig
	AutoRemove  bool
	HSPICE      HSPICEConfig
	WaveView    WaveViewConfig
	Taa         TaaConfig
	MaxRetry    int
}

type Default struct {
	NetListDir string
	BaseDir    string
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
