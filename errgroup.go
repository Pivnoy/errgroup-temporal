package errgroup_temporal

import (
	"sync"

	"go.temporal.io/sdk/workflow"
)

type group struct {
	wg      workflow.WaitGroup
	cancel  func()
	errOnce sync.Once
	err     error
}

func NewErrGroup(ctx workflow.Context) (*group, workflow.Context) {
	ctx, cancel := workflow.WithCancel(ctx)
	return &group{
		cancel: cancel,
		wg:     workflow.NewWaitGroup(ctx),
	}, ctx
}

func (g *group) Wait(ctx workflow.Context) error {
	g.wg.Wait(ctx)
	if g.cancel != nil {
		g.cancel()
	}

	return g.err
}

func (g *group) Go(ctx workflow.Context, f func() error) {
	g.wg.Add(1)
	workflow.Go(ctx, func(_ workflow.Context) {
		defer g.wg.Done()

		if err := f(); err != nil {
			g.errOnce.Do(func() {
				g.err = err
				if g.cancel != nil {
					g.cancel()
				}
			})
		}
	})
}
