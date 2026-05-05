package handler

import (
	"errors"
	"regexp"

	"github.com/GCET-Open-Source-Foundation/coding_arena/backend/adapter"
)

var (
	supportedLanguages = map[string]bool{
		"python": true,
		"cpp":    true,
		"c":      true,
		"java":   true,
		"go":     true,
	}

	problemIDPattern = regexp.MustCompile(`^[a-z0-9][a-z0-9\-]{0,63}$`)
)

const maxCodeLength = 512 * 1024

var judgeAdapter *adapter.JudgeAdapter

// SetAdapter injects the judge adapter into the handler package.
func SetAdapter(a *adapter.JudgeAdapter) {
	judgeAdapter = a
}

func validateInput(source, language, problemID string) error {
	if len(source) > maxCodeLength {
		return errors.New("code exceeds maximum allowed size")
	}
	if !supportedLanguages[language] {
		return errors.New("unsupported language")
	}
	if !problemIDPattern.MatchString(problemID) {
		return errors.New("invalid problem_id: must be 1-64 lowercase alphanumeric characters or hyphens")
	}
	return nil
}
