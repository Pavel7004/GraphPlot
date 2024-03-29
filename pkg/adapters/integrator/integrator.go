package integrator

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
)

type NewIntFunc func(begin, end float64, step float64, saveFn func(t float64, x *circuit.Circuit) error) Integrator

type Integrator interface {
	Integrate(ctx context.Context, st *circuit.Circuit) float64
}
