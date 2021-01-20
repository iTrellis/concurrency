# concurrency
go concurrency library 

* [![GoDoc](http://godoc.org/github.com/iTrellis/concurrency?status.svg)](http://godoc.org/github.com/iTrellis/concurrency)

## Installation

```golang
go get -u github.com/iTrellis/concurrency
```

## Usage

### concurrency repo

```golang
// ConcurrencyRepo functions for go routings to run tasks
type ConcurrencyRepo interface {
	// Invoke tasks: task must be functions
	Invoke(tasks []interface{}) ([]Runner, error)
	InvokeDuration(tasks []interface{}, timeout time.Duration) ([]Runner, error)
}
```

### new and input a namespace's transaction

```golang
	c := concurrency.New(100)

	f := func(i int) (n int) {
		fmt.Println(i, time.Now())
		return i
	}

	var tasks []interface{}
	for i := 0; i < 50; i++ {
		tasks = append(tasks, f)
	}

	runners, err := c.Invoke(tasks)

	for i := 0; i < 50; i++ {
		fmt.Println(runners[i].Get.MapV(func(n int){fmt.Println("%3.d", n)}))
	}
```