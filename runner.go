// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import (
	"time"
)

// Runner functions for a runner
type Runner interface {
	// get execute result
	Get() *Result
	// get execute result in duration
	GetDuration(d time.Duration) (result *Result)
	// judge runner is done
	IsDone() bool
	// judge runner is running
	IsRunning() bool
	// cancel runner
	Cancel() bool
	// is cancelled
	IsCancelled() bool
	// runner to run
	Run()
}
