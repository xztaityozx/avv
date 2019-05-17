package cmd

import "github.com/xztaityozx/avv/repository"

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
	Server           repository.Repository
}

type HSPICEConfig struct {
	Command string
	Option  string
}

type RetryConfig struct {
	HSPICE   int
	WaveView int
	CountUp  int
}

type WaveViewConfig struct {
	Command string
}

type ParallelConfig struct {
	Master int
}
