package metric

import (
	"context"
	"math"
	"os"

	"github.com/shirou/gopsutil/v4/process"
	oMetric "go.opentelemetry.io/otel/metric"

	"go.opentelemetry.io/otel/semconv/v1.43.0/processconv"

	"github.com/lcnascimento/go-kit/errors"
)

const scopeName = "github.com/lcnascimento/go-kit/o11y/process"

// StartProcess registers process.cpu.time and process.memory.usage on mp.
func StartProcess(mp oMetric.MeterProvider) error {
	pid := os.Getpid()
	if pid < 0 || pid > math.MaxInt32 {
		return errors.New("invalid process ID: %d", pid)
	}

	proc, err := process.NewProcess(int32(pid))
	if err != nil {
		return errors.Wrap(err, "failed to inspect current process")
	}

	meter := mp.Meter(scopeName)

	cpuTime, err := processconv.NewCPUTime(meter)
	if err != nil {
		return errors.Wrap(err, "failed to create process.cpu.time")
	}

	memUsage, err := processconv.NewMemoryUsageObservable(meter)
	if err != nil {
		return errors.Wrap(err, "failed to create process.memory.usage")
	}

	userMode := oMetric.WithAttributes(cpuTime.AttrCPUMode(processconv.CPUModeUser))
	systemMode := oMetric.WithAttributes(cpuTime.AttrCPUMode(processconv.CPUModeSystem))

	_, err = meter.RegisterCallback(func(ctx context.Context, o oMetric.Observer) error {
		times, err := proc.TimesWithContext(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to read process cpu times")
		}

		mem, err := proc.MemoryInfoWithContext(ctx)
		if err != nil {
			return errors.Wrap(err, "failed to read process memory")
		}

		o.ObserveFloat64(cpuTime.Inst(), times.User, userMode)
		o.ObserveFloat64(cpuTime.Inst(), times.System, systemMode)
		o.ObserveInt64(memUsage.Inst(), int64(min(mem.RSS, math.MaxInt64)))
		return nil
	}, cpuTime.Inst(), memUsage.Inst())
	if err != nil {
		return errors.Wrap(err, "failed to register process metrics")
	}

	return nil
}
