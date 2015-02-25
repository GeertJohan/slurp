package slurp

import (
	"fmt"
	"strings"
	"sync"
)

type Task func(*C) error

type task struct {
	name string
	deps taskstack
	task Task

	//called bool

	lock sync.Mutex
}

type taskstack map[string]*task

func (t *task) run(c *C) error {

	t.lock.Lock()
	defer t.lock.Unlock()

	//if t.called {
	//		return nil
	//	}
	c.Infof("Starting %s.", t.name)

	failed := make(chan string)
	cancel := make(chan struct{}, len(t.deps))
	var wg sync.WaitGroup
	go func(failed chan string) {
		defer close(failed)
		for name, t := range t.deps {
			select {
			case <-cancel:
				break
			default:

				wg.Add(1)
				go func(t *task, name string) {
					defer wg.Done()
					c.Infof("Waiting for %s", name)
					c := &C{c.New(fmt.Sprintf("%s: ", name))}
					err := t.run(c)
					if err != nil {
						c.Error(err)
						failed <- name
					}
				}(t, name)
			}
		}
		wg.Wait()
	}(failed)

	var failedjobs []string

	for job := range failed {
		cancel <- struct{}{}
		failedjobs = append(failedjobs, job)
	}

	if len(failedjobs) > 0 {
		return fmt.Errorf("Task Canacled. Reason: Failed Dependency (%s).", strings.Join(failedjobs, ","))
	}

	//t.called = true
	err := t.task(c)
	if err == nil {
		c.Infof("Done %s.", t.name)
	}

	return err
}
