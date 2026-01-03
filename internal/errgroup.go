package internal

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

var _ ErrGroup = (*errGroupImpl)(nil)

type errGroupImpl struct {
	wg   workflow.WaitGroup
	once Once
	err  error

	ctx    workflow.Context
	cancel workflow.CancelFunc
}

func (g *errGroupImpl) Wait() error {
	g.wg.Wait(g.ctx)
	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}

func (g *errGroupImpl) Go(f func(ctx workflow.Context) error) {
	g.wg.Add(1)

	workflow.Go(g.ctx, func(ctxFn workflow.Context) {
		defer g.wg.Done()

		if err := f(ctxFn); err != nil {
			if err := g.once.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			}); err != nil {
				g.err = errors.Join(g.err, err)
			}
		}
	})
}
