package errgroup_temporal

import (
	"errors"

	"go.temporal.io/sdk/workflow"
)

type group struct {
	wg      workflow.WaitGroup
	errOnce Once
	err     error

	ctx    workflow.Context
	cancel workflow.CancelFunc
}

func NewErrGroup(ctx workflow.Context) *group {
	ctx, cancel := workflow.WithCancel(ctx)
	return &group{
		ctx:    ctx,
		cancel: cancel,
		wg:     workflow.NewWaitGroup(ctx),
	}
}

func (g *group) Wait() error {
	g.wg.Wait(g.ctx)
	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}

func (g *group) Go(f func(ctx workflow.Context) error) {
	g.wg.Add(1)

	workflow.Go(g.ctx, func(ctxFn workflow.Context) {
		defer g.wg.Done()

		if err := f(ctxFn); err != nil {
			if err := g.errOnce.Do(ctxFn, func() {
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
