// GNU GPL v3 License
// Copyright (c) 2016 github.com:iTrellis

package concurrency

import "reflect"

var (
	errType = reflect.TypeOf((*error)(nil)).Elem()
)

// Result run task return values
type Result struct {
	values []reflect.Value
}

// MapV map values into function parmas
func (p *Result) MapV(params interface{}) error {
	return p.mapTo(params)
}

func (p *Result) mapTo(params interface{}) (err error) {
	if err = p.lastError(); err != nil {
		return
	}

	if params != nil {
		t := reflect.TypeOf(params)
		if t.Kind() != reflect.Func {
			return ErrFailedGetResult
		}

		if t.NumIn() > len(p.values) {
			return ErrFailedGetResult
		}

		tValue := reflect.ValueOf(params)
		tValue.Call(p.values[0:t.NumIn()])
	}
	return nil
}

func (p *Result) lastError() error {
	if lenVals := len(p.values); lenVals > 0 {
		v := p.values[lenVals-1]
		if v.IsValid() && !v.IsNil() && v.Type().ConvertibleTo(errType) {
			return v.Interface().(error)
		}
	}
	return nil
}
