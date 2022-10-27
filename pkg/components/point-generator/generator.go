package pointgenerator

import (
	"context"

	"github.com/Pavel7004/Common/tracing"

	"github.com/Pavel7004/GraphPlot/pkg/adapters/circuit"
)

type endPoint string

const EndPoint endPoint = "end"

func GeneratePoints(ctx context.Context, args *Args) {
	span, ctx := tracing.StartSpanFromContext(ctx)
	defer span.Finish()

	var (
		circuit *circuit.Circuit = args.Circuit.Clone()

		left  float64
		right float64
	)

	right, ok := ctx.Value(EndPoint).(float64)
	if !ok {
		right = 60
	}

	for left < right {
		int := args.NewIntFn(left, right, args.Step, args.SaveFn)

		left = int.Integrate(ctx, circuit)

		circuit.ToggleState()
	}
}