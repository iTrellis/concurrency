// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import "errors"

// concurrency errors
var (
	ErrInterruptedAddTask = errors.New("interrupted to add task")
	ErrFailedGetResult    = errors.New("failed get result")
	ErrLastValueIsError   = errors.New("last value is error")
)
