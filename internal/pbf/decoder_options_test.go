package pbf

import "testing"

func TestDefaultNCpu(t *testing.T) {
	t.Parallel()

	if got := DefaultNCpu(); got < 1 {
		t.Errorf("DefaultNCpu() = %d, want at least 1", got)
	}
}

func TestDecoderOptions(t *testing.T) {
	t.Parallel()

	cfg := defaultDecoderConfig
	WithProtoBatchSize(42)(&cfg)
	WithNCpus(3)(&cfg)

	if cfg.protoBatchSize != 42 {
		t.Errorf("protoBatchSize: got %d, want 42", cfg.protoBatchSize)
	}
	if cfg.nCPU != 3 {
		t.Errorf("nCPU: got %d, want 3", cfg.nCPU)
	}
}
