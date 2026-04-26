package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadJudgeConfig(t *testing.T) {
	// Helper to clear environment variables between tests
	clearEnv := func() {
		os.Unsetenv("JUDGE_TIME_LIMIT")
		os.Unsetenv("JUDGE_MEMORY_LIMIT")
	}

	tests := []struct {
		name          string
		timeLimitEnv  string
		memoryLimitEnv string
		wantErr       bool
		expectedTime  time.Duration
		expectedMem   int
	}{
		{
			name:          "Valid Configuration",
			timeLimitEnv:  "10s",
			memoryLimitEnv: "1024",
			wantErr:       false,
			expectedTime:  10 * time.Second,
			expectedMem:   1024,
		},
		{
			name:          "Empty Configuration Fallback",
			timeLimitEnv:  "",
			memoryLimitEnv: "",
			wantErr:       false,
			expectedTime:  5 * time.Second, // Default fallback
			expectedMem:   512,             // Default fallback
		},
		{
			name:          "Garbage Time Limit Input",
			timeLimitEnv:  "invalid",
			memoryLimitEnv: "1024",
			wantErr:       true,
		},
		{
			name:          "Garbage Memory Limit Input",
			timeLimitEnv:  "10s",
			memoryLimitEnv: "not_a_number",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clearEnv()

			if tt.timeLimitEnv != "" {
				os.Setenv("JUDGE_TIME_LIMIT", tt.timeLimitEnv)
			}
			if tt.memoryLimitEnv != "" {
				os.Setenv("JUDGE_MEMORY_LIMIT", tt.memoryLimitEnv)
			}

			config, err := LoadJudgeConfig()
			if (err != nil) != tt.wantErr {
				t.Fatalf("LoadJudgeConfig() error = %v, wantErr %v", err, tt.wantErr)
			}
			
			if !tt.wantErr {
				if config.TimeLimit != tt.expectedTime {
					t.Errorf("TimeLimit = %v, want %v", config.TimeLimit, tt.expectedTime)
				}
				if config.MemoryLimit != tt.expectedMem {
					t.Errorf("MemoryLimit = %v, want %v", config.MemoryLimit, tt.expectedMem)
				}
			}
		})
	}
	clearEnv() // cleanup after all tests
}
