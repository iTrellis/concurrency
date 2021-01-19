// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import (
	"context"
	"reflect"
	"sync"
	"time"
)

type taskStatus int

const (
	taskDoneStatus taskStatus = iota + 1
	taskRunningStatus
)

// Task task
type Task struct {
	sync.Mutex

	fnCall interface{}
	result chan *Result
	cancel chan cancelTask

	cancelled bool
	status    taskStatus
}

type cancelTask struct {
	interruptFlag bool
}

func newTask(fn interface{}) Runner {
	return &Task{
		fnCall: fn,
		result: make(chan *Result, 1),
		cancel: make(chan cancelTask, 1),
	}
}

// Get get result from context
func (p *Task) Get() (result *Result) {
	return p.getWithContext(context.Background())
}

// GetDuration get result
func (p *Task) GetDuration(d time.Duration) (result *Result) {
	ctx, cf := context.WithTimeout(context.Background(), d)
	defer cf()
	return p.getWithContext(ctx)
}

// Cancel cancel a task
func (p *Task) Cancel() bool {
	p.Lock()
	defer p.Unlock()

	if p.cancelled {
		return true
	}

	if p.IsRunning() {
		return false
	}

	p.cancel <- cancelTask{interruptFlag: true}

	return true
}

func (p *Task) changeStatus(ts taskStatus) {
	p.Lock()
	defer p.Unlock()
	p.status = ts
}

// IsRunning judge a task if running return true
func (p *Task) IsRunning() bool {
	return p.status == taskRunningStatus
}

// IsDone judge a task if was done return true
func (p *Task) IsDone() bool {
	return p.status == taskDoneStatus
}

// IsCancelled judge a task if was cancelled return true
func (p *Task) IsCancelled() bool {
	p.Lock()
	defer p.Unlock()
	return p.cancelled
}

func (p *Task) getWithContext(ctx context.Context) (result *Result) {
	if p.IsCancelled() {
		return
	}

	select {
	case <-ctx.Done():
		return nil
	case result = <-p.result:
		return
	default:
		return
	}
}

// Run run task
func (p *Task) Run() {

	if !p.checkRun() {
		return
	}

	p.checkCancel()

	p.changeStatus(taskRunningStatus)

	p.run()

	p.changeStatus(taskDoneStatus)
}

func (p *Task) checkCancel() bool {
	select {
	case cancel := <-p.cancel:
		if cancel.interruptFlag {
			p.Lock()
			p.cancelled = true
			p.Unlock()
			return true
		}
	default:
	}
	return false
}

func (p *Task) checkRun() bool {
	p.Lock()
	defer p.Unlock()
	if p.cancelled || p.IsDone() {
		return false
	}
	return true
}

func (p *Task) run() {
	var refVals []reflect.Value

	switch fn := p.fnCall.(type) {
	case Runner:
		fn.Run()
	default:
		if reflect.TypeOf(fn).Kind() == reflect.Func {
			fnVal := reflect.ValueOf(fn)
			refVals = fnVal.Call(nil)
		}
	}

	p.result <- &Result{values: refVals}
}
