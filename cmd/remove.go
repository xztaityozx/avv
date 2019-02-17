package cmd

import (
	"bufio"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func (rt RunTask) GarbageCollect() error {

	return nil
}

func (t Task) Remove() error {
	if !t.AutoRemove {
		return nil
	}
	return nil
}

// Remove remove DstDir
// returns: error
func (d SimulationDirectories) Remove() error {
	err := os.RemoveAll(d.DstDir)
	if err != nil {
		log.WithError(err).Error("Failed Remove Directory: ",d.DstDir)
		return err
	}
	return nil
}

func (s SimulationFiles) Remove() error {
	err := os.Remove(s.SPIScript)
	if err != nil {
		return err
	}

	err = os.Remove(s.AddFile.Path)
	if err != nil {
		return err
	}

	return nil
}

type RemoveTask struct {
	Task Task
}

func (r RemoveTask) Run(ctx context.Context) TaskResult {
	ch := make(chan error)
	defer close(ch)
	go func() {
		err := r.Task.SimulationDirectories.Remove()
		if err != nil {
			ch<-err
			return
		}
		err = r.Task.SimulationFiles.Remove()
		ch<-err
	}()

	select {
	case <-ctx.Done():
		return TaskResult{
			Task:r.Task,
			Status:false,
		}
	case err := <-ch:
		if err != nil {
			log.WithError(err).Error("Failed Remove")
		}
		return TaskResult{
			Status:err == nil,
			Task:Task{},
		}
	}
}

func (RemoveTask) Self() Task {
	return Task{}
}

func (RemoveTask) String() string {
	return ""
}

var removeCmd = &cobra.Command{
	Use:   "remove",
	Short: "Remove TaskFile",
	Run: func(cmd *cobra.Command, args []string) {
		remove := func(b, yes bool, p string) error {
			if !b {
				return nil
			}
			if !yes {
				fmt.Println("Remove all file in ", p, "?")
				fmt.Println("[y]:はい [n]:いいえ")
				fmt.Print(">>> ")

				s := bufio.NewScanner(os.Stdin)
				s.Scan()
				if s.Text() != "y" {
					return nil
				}
			}

			err := os.RemoveAll(p)
			if err != nil {
				return err
			}

			FU.TryMkDir(p)

			return nil
		}

		yes, _ := cmd.Flags().GetBool("yes")
		r, _ := cmd.Flags().GetBool("reserve")
		err := remove(r, yes, ReserveDir())
		if err != nil {
			log.WithError(err).Fatal("Failed Remove File")
		}
		f, _ := cmd.Flags().GetBool("failed")
		err = remove(f, yes, FailedDir())
		if err != nil {
			log.WithError(err).Fatal("Failed Remove File")
		}
		d, _ := cmd.Flags().GetBool("done")
		err = remove(d, yes, DoneDir())
		if err != nil {
			log.WithError(err).Fatal("Failed Remove File")
		}
		d, _ = cmd.Flags().GetBool("dust")
		err = remove(d, yes, DustDir())
		if err != nil {
			log.WithError(err).Fatal("Failed Remove File")
		}

		return
	},
}

func init() {
	rootCmd.AddCommand(removeCmd)
	removeCmd.Flags().Bool("reserve", false, "Remove TaskFiles in reserve directory")
	removeCmd.Flags().Bool("failed", false, "Remove TaskFiles in failed directory")
	removeCmd.Flags().Bool("dust", true, "Remove TaskFiles in dust directory")
	removeCmd.Flags().Bool("done", true, "Remove TaskFiles in done directory")

	removeCmd.Flags().BoolP("yes", "y", false, "Skip Confirm")
}
