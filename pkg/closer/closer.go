package closer

import (
	"log"
	"os"
	"os/signal"
	"sync"
)

type closeFn func() error

type Closer struct {
	mu    sync.Mutex
	once  sync.Once
	funcs []closeFn
	done  chan struct{}
}

var globalCloser = New()

func Add(f ...closeFn) {
	globalCloser.Add(f...)
}

func CloseAll() {
	globalCloser.CloseAll()
}

func Wait() {
	globalCloser.Wait()
}

func New(signals ...os.Signal) *Closer {
	c := &Closer{done: make(chan struct{})}
	if len(signals) > 0 {
		go func() {
			ch := make(chan os.Signal, 1)
			signal.Notify(ch, signals...)
			<-ch
			signal.Stop(ch)
			c.CloseAll()
		}()
	}

	return c
}

func (c *Closer) Add(f ...closeFn) {
	c.mu.Lock()
	c.funcs = append(c.funcs, f...)
	c.mu.Unlock()
}

func (c *Closer) CloseAll() {
	c.once.Do(func() {
		defer close(c.done)

		c.mu.Lock()
		funcs := c.funcs
		c.funcs = nil
		c.mu.Unlock()

		errs := make(chan error, len(funcs))
		for _, f := range funcs {
			go func(f closeFn) {
				errs <- f()
			}(f)
		}

		for i := 0; i < cap(errs); i++ {
			if err := <-errs; err != nil {
				log.Println("close err: ", err)
			}
		}
	})
}

func (c *Closer) Wait() {
	<-c.done
}
