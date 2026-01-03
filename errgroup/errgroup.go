package errgroup

import (
	"github.com/Pivnoy/errgroup-temporal/internal"
	"go.temporal.io/sdk/workflow"
)

type (
	ErrGroup = internal.ErrGroup
)

func NewErrGroup(ctx workflow.Context) ErrGroup {
	return internal.NewErrGroup(ctx)
}
