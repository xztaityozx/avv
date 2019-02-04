package cmd

import (
	"context"
	"errors"
	"sync"

	"github.com/fatih/color"
	"github.com/vbauerster/mpb"
	"github.com/vbauerster/mpb/decor"
)

type (
	Result struct {
		Task   Task
		Status bool
	}

	Dispatcher struct {
		WaitGroup   *sync.WaitGroup
		Queue       chan ITask
		Size        int
		Receiver    chan Result
		Name        string
		ProgressBar *mpb.Bar
	}

	PipeLine struct{}
)

func (p PipeLine) Start(ctx context.Context, tasks []SimulationTask) (success []Task, failed []Task, err error) {
	simDis := NewDispatcher("HSPICE")
	wvDis := NewDispatcher("WaveView")
	cuDis := NewDispatcher("CountUp")
	err = nil

	var wg sync.WaitGroup
	wg.Add(3)

	go func() {
		var do []ITask
		for _, v := range tasks {
			do = append(do, v)
		}

		do, retry := p.PipeStart(ctx, simDis, config.ParallelConfig.HSPICE, do)
		failed = append(failed, retry...)
		wg.Done()

		do, retry = p.PipeStart(ctx, wvDis, config.ParallelConfig.WaveView, do)
		failed = append(failed, retry...)
		wg.Done()

		do, retry = p.PipeStart(ctx, cuDis, config.ParallelConfig.CountUp, do)
		failed = append(failed, retry...)

		for _, r := range do {
			success = append(success, r.Self())
		}

		wg.Done()
	}()

	ch := make(chan struct{})
	defer close(ch)
	go func() {
		wg.Wait()
		ch <- struct{}{}
	}()

	select {
	case <-ctx.Done():
		return nil, nil, errors.New("Canceled by context\n")
	case <-ch:
	}

	return
}

func (r Result) Next() ITask {
	if !r.Status {
		return r.Task.GetWrapper()
	}

	if r.Task.Stage == HSPICE {
		r.Task.Stage = WaveView
		return ExtractTask{
			Task: r.Task,
		}
	}

	r.Task.Stage = CountUp
	return CountTask{
		Task: r.Task,
	}
}

func (p PipeLine) PipeStart(ctx context.Context, d Dispatcher, parallel int, tasks []ITask) (done []ITask, failed []Task) {
	var do = tasks
	done = []ITask{}
	var retry []ITask
	failed = []Task{}

	for len(done) != len(tasks) && config.AutoRetry {
		res := d.Dispatch(ctx, parallel, do)
		for _, r := range res {
			next := r.Next()
			if r.Status {
				done = append(done, next)
			} else {
				retry = append(retry, next)
			}
		}
		do = retry
	}

	for _, r := range retry {
		failed = append(failed, r.Self())
	}

	return done, failed
}

func NewDispatcher(name string) Dispatcher {
	return Dispatcher{
		WaitGroup: &sync.WaitGroup{},
		Name:      name,
	}
}

func (d *Dispatcher) Worker(parent context.Context) {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	for j := range d.Queue {
		rec := make(chan Result)
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

func (d *Dispatcher) Dispatch(parent context.Context, workers int, t []ITask) []Result {
	ctx, cancel := context.WithCancel(parent)
	defer cancel()

	l := log.WithField("at", "dispatcher")

	// start workers
	l.Info("Start Dispatcher")
	for i := 0; i < workers; i++ {
		go func() { d.Worker(ctx) }()
	}
	l.Info(workers, " workers was started")

	d.Size = len(t)
	// make result channel
	d.Receiver = make(chan Result, d.Size)
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

	var results []Result

	for r := range d.Receiver {
		results = append(results, r)
	}

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
