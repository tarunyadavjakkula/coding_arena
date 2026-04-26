package config

import "time"

// JudgeConfig holds the limits for code execution.
type JudgeConfig struct {
	TimeLimit   time.Duration
	MemoryLimit int // in MB
}
