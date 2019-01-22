// Copyright © 2019 xztaityozx
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package cmd

import (
	"bufio"
	"fmt"
	"github.com/mattn/go-colorable"
	"github.com/sirupsen/logrus"
	"github.com/snowzach/rotatefilehook"
	"os"
	"path/filepath"
	"time"

	"github.com/mitchellh/go-homedir"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var cfgFile string
var config = Config{}

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "avv",
	Short: "",
	Long:  ``,
	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.config/avv/.avv.json)")

	//Parallel Options
	rootCmd.PersistentFlags().IntP("Parallel", "p", 2, "HSPICE,WV,数え上げ全ての並列数です")
	rootCmd.PersistentFlags().Int("pHSPICE", 2, "HSPICEの並列数です")
	rootCmd.PersistentFlags().Int("pWV", 2, "WaveViewの並列数です")
	rootCmd.PersistentFlags().Int("pCountUp", 2, "数え上げの並列数です")

	// BindFlags
	viper.BindPFlag("Default.ParallelConfig.HSPICE", rootCmd.PersistentFlags().Lookup("pHSPICE"))
	viper.BindPFlag("Default.ParallelConfig.WaveView", rootCmd.PersistentFlags().Lookup("pWV"))
	viper.BindPFlag("Default.ParallelConfig.CountUp", rootCmd.PersistentFlags().Lookup("pCountUp"))

	cobra.OnInitialize(initConfig)
}

var log = logrus.New()

func initLogger() {

	// Hook to log file
	path := PathJoin(config.LogDir, time.Now().Format("2006-01-02-15-04-05"))
	hook, err := rotatefilehook.NewRotateFileHook(rotatefilehook.RotateFileConfig{
		Filename: path,
		MaxAge:   28,
		MaxSize:  500,
		Level:    logrus.InfoLevel,
		Formatter: &logrus.TextFormatter{
			DisableColors:   true,
			TimestampFormat: time.RFC3339,
			FullTimestamp:   true,
		},
	})
	if err != nil {
		logrus.Fatal(err)
	}

	log.SetLevel(logrus.InfoLevel)
	log.SetOutput(colorable.NewColorableStderr())
	// init logrus System
	log.SetFormatter(&logrus.TextFormatter{
		ForceColors:     true,
		FullTimestamp:   true,
		TimestampFormat: time.RFC3339,
	})

	log.AddHook(hook)
}

// initConfig reads in config file and ENV variables if set.
func initConfig() {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedir.Dir()
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}

		// Search config in home directory with name ".avv" (without extension).
		viper.AddConfigPath(filepath.Join(home, ".config", "avv"))
		viper.SetConfigName(".avv")
		viper.SetConfigType("json")
	}

	viper.AutomaticEnv() // read in environment variables that match

	if err := viper.ReadInConfig(); err != nil {
		log.Warn(err)

		// Invalid config file ?
		fmt.Print("設定ファイルがなんか変だけど大丈夫ですか？ (y/n) >>>")

		s := bufio.NewScanner(os.Stdin)
		s.Scan()
		res := s.Text()
		if res != "y" && res != "yes" {
			return
		}
	}

	// Unmarshal config file
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatal(err)
	}

	// init Logger System
	initLogger()
}
