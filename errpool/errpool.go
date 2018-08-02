package errpool

import (
	"context"

	"golang.org/x/sync/errgroup"
)

type Pool struct {
	g *errgroup.Group
	tasks chan func() error
	closed bool
}

func WithContext(ctx context.Context, maxTasks int) (*Pool, context.Context) {
	g, ctx := errgroup.WithContext(ctx)
	tasks := make(chan func() error, maxTasks)
	pool := &Pool{
		g: g,
		tasks: tasks,
	}
	for i := 0; i < maxTasks; i++ {
		g.Go(func() error {
			for f := range tasks {
				err := f()
				if err != nil {
					return err
				}
			}
			return nil
		})
	}
	return pool, ctx		
}

func (p *Pool) Go(f func() error) {
	p.tasks <- f
}

func (p *Pool) Wait() error {
	if !p.closed {
		p.closed = true
		close(p.tasks)
	}
	return p.g.Wait()
}
