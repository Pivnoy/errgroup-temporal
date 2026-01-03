package internal

import "go.temporal.io/sdk/workflow"

var _ Once = (*onceImpl)(nil)

type onceImpl struct {
	mtx  workflow.Mutex
	done bool

	ctx workflow.Context
}

func (o *onceImpl) Do(f func()) error {
	if err := o.mtx.Lock(o.ctx); err != nil {
		return err
	}
	defer o.mtx.Unlock()

	if !o.done {
		f()
		o.done = true
	}

	return nil
}
