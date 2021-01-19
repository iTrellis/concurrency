// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import (
	"time"
)

// Repo functions for go routings to run tasks
type Repo interface {
	// Invoke tasks: task must be functions
	Invoke(tasks []interface{}) ([]Runner, error)
	InvokeDuration(tasks []interface{}, timeout time.Duration) ([]Runner, error)
}

// New  return a concurrency pool with number
func New(number int) Repo {
	return newPool(number)
}
