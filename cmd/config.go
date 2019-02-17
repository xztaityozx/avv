package cmd

import (
	"time"
)

type Config struct {
	Default          Task
	DefaultSEEDRange SEED
	SlackConfig      SlackConfig
	MachineName      string
	LogDir           string
	TaskDir          string
	Templates        Templates
	ParallelConfig   ParallelConfig
	HSPICE           HSPICEConfig
	WaveView         WaveViewConfig
	RetryConfig      RetryConfig
	AutoRetry        bool
	AutoDBBackUp     bool
	BackUpDir        string
}

type HSPICEConfig struct {
	Command string
	Option  string
}

type RetryConfig struct {
	HSPICE   int
	WaveView int
	CountUp  int
	DBAccess int
}

type WaveViewConfig struct {
	Command string
}

type SlackConfig struct {
	Channel string
	User    string
	Token   string
}

type SlackMessage struct {
	Success      int
	Failed       int
	StartTime    time.Time
	FinishedTime time.Time
	SubMessage   string
	ErrorMassage string
}

type ParallelConfig struct {
	DB       int
	HSPICE   int
	WaveView int
	CountUp  int
}

// Compare func for ParallelConfig struct
func (s ParallelConfig) Compare(t ParallelConfig) bool {
	return s.HSPICE == t.HSPICE && s.WaveView == t.WaveView && s.CountUp == t.CountUp && s.DB == t.DB
}
