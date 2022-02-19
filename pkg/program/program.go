package program

import (
	"context"
	"image/color"
	"math"
	"os"
	"path"

	"github.com/Pavel7004/GraphPlot/pkg/circuit"
	"github.com/Pavel7004/GraphPlot/pkg/cli"
	misc "github.com/Pavel7004/GraphPlot/pkg/common/misc"
	"github.com/Pavel7004/GraphPlot/pkg/common/tracing"
	"github.com/Pavel7004/GraphPlot/pkg/graph"
	"github.com/Pavel7004/GraphPlot/pkg/integrator"
	bogatskiyshampin "github.com/Pavel7004/GraphPlot/pkg/integrator/bogatskiy-Shampin"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/euler"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint"
	midpointimplicit "github.com/Pavel7004/GraphPlot/pkg/integrator/midpoint-implicit"
	threeeighth "github.com/Pavel7004/GraphPlot/pkg/integrator/three-eighth"
	"github.com/Pavel7004/GraphPlot/pkg/integrator/trapeziod"
)

func Run(ctx context.Context, circ *circuit.Circuit, folderName string, buffSize, dpi int) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	integrators := []integrator.NewIntFunc{
		euler.NewEulerInt,
		midpoint.NewMidpointInt,
		midpointimplicit.NewMidpointImplInt,
		bogatskiyshampin.NewShampinInt,
		threeeighth.NewThreeEighthInt,
		trapeziod.NewTrapezoidInt,
	}

	if err := os.MkdirAll(folderName, os.ModePerm); err != nil {
		panic(err)
	}

	for _, int := range integrators {
		gr := graph.NewInfoPlotter(buffSize, dpi)

		PlotSystem(ctx, gr, circ, int)
		PlotTheory(ctx, gr, circ)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_theory.png"))
	}

	for _, int := range integrators {
		gr := graph.NewInfoPlotter(buffSize, dpi)

		PlotDiffFunc(ctx, gr, circ, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_diffErr.png"))
	}

	for _, int := range integrators {
		gr := graph.NewInfoPlotter(buffSize, dpi)

		ctx := context.WithValue(ctx, "end", 200.0)
		PlotSystem(ctx, gr, circ, int)

		gr.SaveToFile(ctx, path.Join(folderName, misc.GetFuncModule(int)+"_multiTicks.png"))
	}
}

func PlotSystem(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit, newInt integrator.NewIntFunc) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		st     = circ.Clone()
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
	)

	end, ok := ctx.Value("end").(float64)
	if !ok {
		end = 60
	}

	for right <= end {
		int := newInt(left, right, cli.Step, func(t float64, x *circuit.Circuit) {
			gr.AddPoint(t, x.GetLoadVoltage())
		})

		int.Integrate(ctx, st)
		st.ToggleState()
		left = right
		right += period
	}
}

func PlotTheory(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	st := circ.Clone()
	gr.PlotFunc(color.RGBA{R: 255, A: 255}, st.GetLoadVoltageFunc())
}

func PlotDiffFunc(ctx context.Context, gr *graph.InfoPlotter, circ *circuit.Circuit, newInt integrator.NewIntFunc) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		st     = circ.Clone()
		period = st.GetSystemPeriod()
		left   = 0.0
		right  = period
		theory = st.GetLoadVoltageFunc()
	)

	gr.SetYLabel("x(t), %")
	for right <= 60 {
		int := newInt(left, right, cli.Step, func(t float64, x *circuit.Circuit) {
			vol := x.GetLoadVoltage()
			if vol < 0.0001 {
				gr.AddPoint(t, 0.0)
			} else {
				gr.AddPoint(t, math.Abs(vol-theory(t))/vol*100)
			}
		})

		int.Integrate(ctx, st)
		st.ToggleState()
		left = right
		right += period
	}
}
