package cmd

import (
	"time"
)

var config Config

type Config struct {
	TaskDir        string
	DstDir         string
	SrcDir         string
	SlackConfig    SlackConfig
	ParallelConfig ParallelConfig
	LogDir         string
	MachineName    string
	PlotPoint      PlotPoint
	Templates      Templates
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
	HSPICE   int
	WaveView int
	CountUp  int
}

// Compare func for ParallelConfig struct
func (s ParallelConfig) Compare(t ParallelConfig) bool {
	return s.HSPICE == t.HSPICE && s.WaveView == t.WaveView && s.CountUp == t.CountUp
}
