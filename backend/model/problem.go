package model

// Problem represents the metadata for a coding problem.
type Problem struct {
	ID          string  `json:"id"`
	Title       string  `json:"title"`
	Difficulty  string  `json:"difficulty"`
	Category    string  `json:"category"`
	Points      float64 `json:"points"`
	SolvedCount int     `json:"solvedCount"`
}

// ProblemDetail extends Problem with full statement details.
type ProblemDetail struct {
	Problem
	Description  string        `json:"description"`
	InputFormat  string        `json:"inputFormat"`
	OutputFormat string        `json:"outputFormat"`
	Constraints  []string      `json:"constraints"`
	Examples     []TestExample `json:"examples"`
	TimeLimit    int           `json:"timeLimit"`   // ms
	MemoryLimit  int           `json:"memoryLimit"` // MB
}

// TestExample represents a sample test case for problem description.
type TestExample struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}
