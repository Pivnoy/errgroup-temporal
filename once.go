package errgroup_temporal

import (
	"go.temporal.io/sdk/workflow"
)

type Once struct {
	mtx  workflow.Mutex
	done bool
}

func (o *Once) Do(ctx workflow.Context, f func()) error {
	if err := o.mtx.Lock(ctx); err != nil {
		return err
	}
	defer o.mtx.Unlock()

	if !o.done {
		f()
		o.done = true
	}

	return nil
}
