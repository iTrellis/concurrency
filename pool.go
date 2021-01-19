// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import (
	"context"
	"sync"
	"time"
)

type poolExector struct {
	sync.Mutex
	once sync.Once

	number     int
	doneNumber int

	runnerStack *RunnerStack
	runnerChan  chan Runner
	runners     []Runner

	interrupted bool
	interrupt   chan struct{}

	workers    map[int]chan struct{}
	dispatcher chan struct{}
}

func newPool(number int) Repo {
	exector := &poolExector{
		number:      number,
		runnerStack: newRunnerStack(),
		runnerChan:  make(chan Runner, number*2),

		workers:   make(map[int]chan struct{}, number),
		interrupt: make(chan struct{}),
	}

	exector.init()

	return exector
}

func (p *poolExector) init() {
	p.once.Do(func() {
		p.beginWorkers()
		p.beginDispatcher()
	})
}

func (p *poolExector) beginWorkers() {

	for i := 0; i < p.number; i++ {
		workerChan := make(chan struct{})
		p.workers[i] = workerChan
		go p.handler(i, workerChan)
	}
}

func (p *poolExector) handler(id int, workerChan chan struct{}) {
	for {
		select {
		case <-workerChan:
			return
		case runner, ok := <-p.runnerChan:
			if !ok {
				return
			}

			if p.interrupted {
				runner.Cancel()
				runner.Run()
				continue
			}

			runner.Run()
		}
	}
}

func (p *poolExector) beginDispatcher() {

	p.dispatcher = make(chan struct{})

	go func() {
		for {
			select {
			case <-p.dispatcher:
				return
			default:
			}

			if runner := p.runnerStack.Pop(); runner != nil {
				p.runnerChan <- runner
			}
		}
	}()
}

// Invoke
func (p *poolExector) Invoke(tasks []interface{}) ([]Runner, error) {
	return p.InvokeDuration(tasks, -1)
}

// InvokeDuration
func (p *poolExector) InvokeDuration(tasks []interface{}, timeout time.Duration) (runners []Runner, err error) {
	if p.isInterrupted() {
		return
	}

	var rs []Runner

	defer func() {
		if err != nil || len(rs) > 0 {
			for _, v := range rs {
				v.Cancel()
			}
		}
	}()

	wg := &sync.WaitGroup{}
	for _, v := range tasks {
		var runner Runner
		if runner, err = p.addTask(v); err != nil {
			return
		}
		rs = append(rs, runner)

		wg.Add(1)

		ctx := context.Background()
		var cancelFunc context.CancelFunc

		if timeout > 0 {
			ctx, cancelFunc = context.WithTimeout(ctx, timeout)
			defer cancelFunc()
		}

		go func(r Runner, c context.Context) {
			defer wg.Done()
			for {
				select {
				case <-c.Done():
					if !r.IsDone() {
						r.Cancel()
					}
					return
				default:
					if r.IsDone() || p.interrupted {
						return
					}
				}
			}
		}(runner, ctx)
	}

	wg.Wait()
	runners = rs

	return
}

func (p *poolExector) Interrupt() {
	p.Lock()
	defer p.Unlock()

	if p.isInterrupted() {
		return
	}

	close(p.interrupt)

	go func() {
		totalCount := len(p.runners)
		for p.doneNumber != totalCount {
			time.Sleep(time.Millisecond * 100)
			p.doneNumber = 0
			for _, v := range p.runners {
				if v.IsDone() || v.IsCancelled() {
					p.doneNumber++
				}
			}
		}
		for _, c := range p.workers {
			close(c)
		}
		close(p.dispatcher)
	}()
}

func (p *poolExector) addTask(fn interface{}) (Runner, error) {
	p.Lock()
	defer p.Unlock()
	if p.isInterrupted() {
		return nil, ErrInterruptedAddTask
	}

	r := newTask(fn)
	p.runnerStack.Push(r)
	p.runners = append(p.runners, r)
	return r, nil
}

func (p *poolExector) isInterrupted() (b bool) {
	select {
	case <-p.interrupt:
		b = true
	default:
	}
	return
}
