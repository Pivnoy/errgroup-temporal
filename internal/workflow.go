package internal

import "go.temporal.io/sdk/workflow"

type (
	Once interface {
		Do(f func()) error
	}

	ErrGroup interface {
		Go(f func(ctx workflow.Context) error)
		Wait() error
	}
)

func NewOnce(ctx workflow.Context) Once {
	return &onceImpl{
		mtx:  workflow.NewMutex(ctx),
		done: false,
		ctx:  ctx,
	}
}

func NewErrGroup(ctx workflow.Context) ErrGroup {
	ctx, cancel := workflow.WithCancel(ctx)
	once := NewOnce(ctx)
	return &errGroupImpl{
		wg:     workflow.NewWaitGroup(ctx),
		once:   once,
		ctx:    ctx,
		cancel: cancel,
	}
}
