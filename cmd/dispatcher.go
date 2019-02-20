package cmd

import (
	"context"
	"errors"
	"sync"
	"time"

	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type (
	TaskResult struct {
		Task   Task
		Status bool
	}

	Dispatcher struct {
		WaitGroup   *sync.WaitGroup
		Queue       chan ITask
		Size        int
		Receiver    chan TaskResult
		Name        string
		ProgressBar *mpb.Bar
	}

	PipeLine struct{}
	Pipe     struct {
		Name       string
		Converter  func(Task) ITask
		FailedConverter  func(Task) ITask
		Parallel   int
		RetryLimit int
		AutoRetry  bool
	}
)

func (p PipeLine) Start(ctx context.Context, input []ITask, pipe ...Pipe) (success []Task, failed []Task, err error) {
	success = []Task{}
	failed = []Task{}

	ch := make(chan error)
	defer close(ch)

	go func() {
		for _, pi := range pipe {
			rt, f, err := pi.Connect(ctx, input)
			if err != nil {
				log.Error(err)
			}

			failed = append(failed, f...)
			input = rt
		}
		ch <- nil
	}()

	select {
	case <-ctx.Done():
		return nil, nil, errors.New("PipeLine: canceled by context")
	case err = <-ch:
	}

	for _, t := range input {
		success = append(success, t.Self())
	}

	return
}

func (p Pipe) Connect(ctx context.Context, input []ITask) (success []ITask, failed []Task, err error) {
	tasks := len(input)
	var do = input

	ch := make(chan struct{})
	defer close(ch)

	cnt := 0
	go func() {

		for ok := true; ok; ok = len(success) != tasks && (p.AutoRetry || config.AutoRetry) {
			dis := NewDispatcher(p.Name)
			res := dis.Dispatch(ctx, p.Parallel, do)
			do = []ITask{}
			for _, r := range res {
				if r.Status {
					success = append(success, p.Converter(r.Task))
				} else {
					do = append(do, p.FailedConverter(r.Task))
				}
			}

			if len(do) != 0 && cnt >= p.RetryLimit {
				err = errors.New("Retry Limit Exceeded ")
				ch <- struct{}{}
				return
			} else if len(do) != 0 {
				cnt++
				log.Warn("Pipe.Connect: Retry(", cnt, ")")
			}

		}
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, nil, nil
	case <-ch:
	}

	for _, t := range do {
		failed = append(failed, t.Self())
	}
	return
}

func NewDispatcher(name string) Dispatcher {
	return Dispatcher{
		Name: name,
	}
}

func (d *Dispatcher) Worker(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	for j := range d.Queue {
		rec := make(chan TaskResult)
		go func() {
			rec <- j.Run(ctx)
		}()

		select {
		case <-ctx.Done():
			return
		case res := <-rec:
			d.Done()
			d.Receiver <- res
		}

		close(rec)
	}
}

func (d *Dispatcher) Dispatch(parent context.Context, workers int, t []ITask) []TaskResult {

	if len(t) < workers {
		workers=len(t)
	}

	d.WaitGroup = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	l := log.WithField("at", "dispatcher")
	l.Info(d.Name, " job Start")
	l.Info("Parallel: ", workers)
	l.Info("Begin: ", time.Now().Format(time.ANSIC))

	d.Size = len(t)
	// make result channel
	d.Receiver = make(chan TaskResult, d.Size)
	// enqueue jobs
	d.Queue = make(chan ITask, d.Size)
	for _, v := range t {
		d.Add(v)
	}
	defer close(d.Queue)

	pb := mpb.New(mpb.WithContext(ctx), mpb.WithWaitGroup(d.WaitGroup))

	barName := color.New(color.FgHiYellow).Sprint("Dispatcher:") + d.Name

	finishMSG := color.New(color.FgHiGreen).Sprint(" done!")
	workingMSG := color.New(color.FgCyan).Sprint(" In progress...")

	d.ProgressBar = pb.AddBar(int64(d.Size),
		mpb.BarStyle("┃██▒┃"),
		mpb.BarWidth(50),
		//mpb.BarClearOnComplete(),
		mpb.PrependDecorators(
			decor.Name(barName, decor.WC{W: len(barName) + 1, C: decor.DidentRight}),
		),
		mpb.AppendDecorators(
			decor.Name("   "),
			decor.Percentage(decor.WC{W: 5}),
			decor.Name(" | "),
			decor.CountersNoUnit("%d / %d", decor.WCSyncWidth),
			decor.Name(" | "),
			decor.OnComplete(decor.Name(workingMSG), finishMSG),
		),
	)
	// start workers
	for i := 0; i < workers; i++ {
		go func() { d.Worker(ctx) }()
	}

	ch := make(chan struct{})
	defer close(ch)
	go func() {
		pb.Wait()
		d.WaitGroup.Wait()
		ch <- struct{}{}
	}()

	// wait trap or finish
	select {
	case <-ctx.Done():
		close(d.Receiver)
		return nil
	case <-ch:
	}

	close(d.Receiver)

	var results []TaskResult

	for r := range d.Receiver {
		results = append(results, r)
	}

	l.Info("End: ", time.Now().Format(time.ANSIC))

	return results
}

func (d *Dispatcher) Done() {
	d.ProgressBar.Increment()
	d.WaitGroup.Done()
}

func (d *Dispatcher) Add(job ITask) {
	d.WaitGroup.Add(1)
	d.Queue <- job
}
