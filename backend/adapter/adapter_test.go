package adapter

import (
	"testing"
	"time"

	"github.com/Aerosane/coding_arena/backend/bridge"
	"github.com/Aerosane/coding_arena/backend/config"
)

func TestJudgeAdapterConfig(t *testing.T) {
	b := bridge.New(":9999", "test", "key")
	
	cfg := &config.JudgeConfig{
		TimeLimit:   5 * time.Second,
		MemoryLimit: 512,
	}

	adapt := New(b, cfg)
	
	if adapt.cfg.TimeLimit != 5*time.Second {
		t.Errorf("expected TimeLimit 5s, got %v", adapt.cfg.TimeLimit)
	}
	
	if adapt.cfg.MemoryLimit != 512 {
		t.Errorf("expected MemoryLimit 512, got %v", adapt.cfg.MemoryLimit)
	}
}
