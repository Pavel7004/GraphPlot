package midpointimplicit

import (
	"context"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
)

type MidpointImpInt struct {
	begin  float64
	end    float64
	step   float64
	saveFn func(t float64, x *circuit.Circuit)
}

var _ integrator.Integrator = (*MidpointImpInt)(nil)

func NewMidpointImplInt(begin, end, step float64, saveFn func(t float64, x *circuit.Circuit)) integrator.Integrator {
	return &MidpointImpInt{
		begin:  begin,
		end:    end,
		step:   step,
		saveFn: saveFn,
	}
}

func (si *MidpointImpInt) Integrate(ctx context.Context, circ *circuit.Circuit) float64 {
	span, ctx := tracing.StartSpanFromContext(ctx)
	span.SetTag("StartPoint", si.begin)
	span.SetTag("EndPoint", si.end)
	span.SetTag("Step", si.step)
	span.SetTag("RK-stages", 2)

	defer span.Finish()

	var (
		t    = si.begin
		last bool
	)

	for !last {
		if t+si.step > si.end {
			last = true
			si.step = si.end - t
		}

		k1 := circ.Clone()
		k1.ApplyDerivative(si.step/2, k1.GetDerivative())
		circ.ApplyDerivative(si.step, k1.GetDerivative())

		t += si.step
		si.saveFn(t, circ)
	}

	return t
}
