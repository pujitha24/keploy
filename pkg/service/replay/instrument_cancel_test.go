package replay

import (
	"context"
	"errors"
	"testing"

	"go.keploy.io/server/v3/config"
	"go.keploy.io/server/v3/pkg/models"
	"go.uber.org/zap"
)

// setupInstrumentation fails Setup with a fixed error.
type setupInstrumentation struct {
	Instrumentation
	err error
}

func (s setupInstrumentation) Setup(context.Context, string, models.SetupOptions) error {
	return s.err
}

// A replay cancelled while the environment is being set up must surface the
// cancellation from Instrument, so callers treat the run as interrupted and
// do not log it at ERROR. Instrument used to return a fresh error that
// dropped the cause.
func TestInstrumentKeepsSetupCancellation(t *testing.T) {
	boom := errors.New("boom")
	for _, tc := range []struct {
		name string
		err  error
	}{
		{"canceled", context.Canceled},
		{"other failure", boom},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := &Replayer{logger: zap.NewNop(), config: &config.Config{}, instrument: true, instrumentation: setupInstrumentation{err: tc.err}}
			_, err := r.Instrument(context.Background())
			if !errors.Is(err, tc.err) {
				t.Fatalf("Instrument error %v does not wrap %v", err, tc.err)
			}
		})
	}
}
