package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type JudgeConfig struct {
	TimeLimit   time.Duration
	MemoryLimit int // in MB
}

// LoadJudgeConfig reads and parses the judge configuration.
// It uses default values if the env vars are empty, and returns an error
// if the strings are malformed.
func LoadJudgeConfig() (*JudgeConfig, error) {
	config := &JudgeConfig{
		TimeLimit:   5 * time.Second, // Default time limit
		MemoryLimit: 512,             // Default memory limit in MB
	}

	timeStr := os.Getenv("JUDGE_TIME_LIMIT")
	if timeStr != "" {
		parsedTime, err := time.ParseDuration(timeStr)
		if err != nil {
			return nil, fmt.Errorf("invalid JUDGE_TIME_LIMIT: %w", err)
		}
		config.TimeLimit = parsedTime
	}

	memStr := os.Getenv("JUDGE_MEMORY_LIMIT")
	if memStr != "" {
		parsedMem, err := strconv.Atoi(memStr)
		if err != nil {
			return nil, fmt.Errorf("invalid JUDGE_MEMORY_LIMIT: %w", err)
		}
		config.MemoryLimit = parsedMem
	}

	return config, nil
}
