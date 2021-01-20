// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import (
	"github.com/iTrellis/common/data-structures/stack"
)

// RunnerStack runner stack
type RunnerStack struct {
	stack stack.Stack
}

func newRunnerStack() *RunnerStack {
	return &RunnerStack{
		stack: stack.New(),
	}
}

// Push push a runner to stack
func (p *RunnerStack) Push(r Runner) {
	p.stack.Push(r)
}

// Pop pop a runner from stack
func (p *RunnerStack) Pop() Runner {
	v, ok := p.stack.Pop()
	if !ok {
		return nil
	}
	if ok {
		if r, ok := v.(Runner); ok {
			return r
		}
	}
	return nil
}

// PopAll pop all runners from stack
func (p *RunnerStack) PopAll() (rs []Runner) {
	vs, ok := p.stack.PopAll()
	if !ok {
		return
	}
	for _, v := range vs {
		if r, ok := v.(Runner); ok {
			rs = append(rs, r)
		}
	}
	return
}

// Length stack length
func (p *RunnerStack) Length() int64 {
	return p.stack.Length()
}
