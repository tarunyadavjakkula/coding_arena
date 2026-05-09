package handler

import (
	"errors"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/goccy/go-yaml"
)

// Problem defines the response object for a problem, matching the frontend's expected schema.
type Problem struct {
	ID           string        `json:"id"`
	Title        string        `json:"title"`
	Difficulty   string        `json:"difficulty"`
	Category     string        `json:"category"`
	Points       int           `json:"points"`
	SolvedCount  int           `json:"solvedCount"`
	Description  string        `json:"description"`
	InputFormat  string        `json:"inputFormat"`
	OutputFormat string        `json:"outputFormat"`
	Constraints  []string      `json:"constraints"`
	Examples     []TestExample `json:"examples"`
	TimeLimit    int           `json:"timeLimit"`
	MemoryLimit  int           `json:"memoryLimit"`
}

type TestExample struct {
	Input       string `json:"input"`
	Output      string `json:"output"`
	Explanation string `json:"explanation,omitempty"`
}

// ProblemConfig defines the schema for metadata found in init.yml
type ProblemConfig struct {
	ID           string        `yaml:"id"`
	Slug         string        `yaml:"slug"`
	Title        string        `yaml:"title"`
	Name         string        `yaml:"name"`
	Difficulty   string        `yaml:"difficulty"`
	Category     string        `yaml:"category"`
	Points       int           `yaml:"points"`
	SolvedCount  int           `yaml:"solved_count"`
	Description  string        `yaml:"description"`
	InputFormat  string        `yaml:"input_format"`
	OutputFormat string        `yaml:"output_format"`
	Constraints  []string      `yaml:"constraints"`
	Examples     []TestExample `yaml:"examples"`
	TimeLimit    float64       `yaml:"time_limit"`
	MemoryLimit  int           `yaml:"memory_limit"`
}

// titleCase is a simple utility to convert a slug like "two-sum" to "Two Sum"
func titleCase(s string) string {
	words := strings.Split(s, "-")
	for i, word := range words {
		if len(word) > 0 {
			words[i] = strings.ToUpper(word[:1]) + word[1:]
		}
	}
	return strings.Join(words, " ")
}

// normalizeDifficulty ensures difficulty levels are in the canonical TitleCase format
// expected by the frontend (Easy, Medium, Hard).
func normalizeDifficulty(d string) string {
	switch strings.ToLower(strings.TrimSpace(d)) {
	case "easy":
		return "Easy"
	case "medium":
		return "Medium"
	case "hard":
		return "Hard"
	default:
		// Return original if unknown, but title-cased if possible
		if len(d) > 0 {
			return strings.ToUpper(d[:1]) + strings.ToLower(d[1:])
		}
		return "Medium" // Default fallback
	}
}

// configPath overrides getProblemsPath for testing.
// configPath must not be mutated concurrently; tests must not call t.Parallel().
var configPath = ""

func getProblemsPath() string {
	if configPath != "" {
		return configPath
	}
	path := os.Getenv("PROBLEMS_PATH")
	if path == "" {
		path = "judge-config/problems"
	}
	return path
}

func parseProblem(problemID, initPath string) (Problem, error) {
	prob := Problem{
		ID:         problemID,
		Title:      titleCase(problemID), // Fallback title
		Difficulty: "Medium",            // Default difficulty
	}

	data, err := os.ReadFile(initPath)
	if err != nil {
		return prob, err
	}

	var config ProblemConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return prob, err
	}

	overrideID := ""
	if config.ID != "" {
		overrideID = config.ID
	} else if config.Slug != "" {
		overrideID = config.Slug
	}

	if overrideID != "" {
		if problemIDPattern.MatchString(overrideID) {
			prob.ID = overrideID
			prob.Title = titleCase(overrideID) // Update fallback to match new ID
		} else {
			log.Printf("[WARN] ignored invalid problem id override for problem %s: %q", problemID, overrideID)
		}
	}

	if config.Title != "" {
		prob.Title = config.Title
	} else if config.Name != "" {
		prob.Title = config.Name
	}
	if config.Difficulty != "" {
		prob.Difficulty = normalizeDifficulty(config.Difficulty)
	} else {
		prob.Difficulty = "Medium" // Default if missing
	}
	if config.Category != "" {
		prob.Category = config.Category
	}
	if config.Points > 0 {
		prob.Points = config.Points
	}
	if config.SolvedCount > 0 {
		prob.SolvedCount = config.SolvedCount
	}
	if config.Description != "" {
		prob.Description = config.Description
	}
	if config.InputFormat != "" {
		prob.InputFormat = config.InputFormat
	}
	if config.OutputFormat != "" {
		prob.OutputFormat = config.OutputFormat
	}
	if len(config.Constraints) > 0 {
		prob.Constraints = config.Constraints
	}
	if len(config.Examples) > 0 {
		prob.Examples = config.Examples
	}
	if config.TimeLimit > 0 {
		prob.TimeLimit = int(config.TimeLimit * 1000)
	}
	if config.MemoryLimit > 0 {
		prob.MemoryLimit = config.MemoryLimit / 1024
	}

	return prob, nil
}

// GetProblems handles GET /api/problems.
// It reads the judge-config/problems directory at runtime and returns a list of problems.
func GetProblems(c *gin.Context) {
	problemsPath := getProblemsPath()
	
	entries, err := os.ReadDir(problemsPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			log.Printf("[WARN] problems directory not found: %s — check PROBLEMS_PATH env var", problemsPath)
			c.JSON(http.StatusOK, gin.H{"problems": []Problem{}})
			return
		}
		log.Printf("[ERROR] failed to read config dir: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	problems := []Problem{}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}

		problemID := entry.Name()
		if !problemIDPattern.MatchString(problemID) {
			log.Printf("[WARN] skipping problem with invalid directory name: %q", problemID)
			continue
		}
		problemPath := filepath.Join(problemsPath, problemID)
		initPath := filepath.Join(problemPath, "init.yml")

		prob, err := parseProblem(problemID, initPath)
		if err != nil && !errors.Is(err, os.ErrNotExist) {
			log.Printf("[WARN] failed to read config for problem %s: %v", problemID, err)
		}

		problems = append(problems, prob)
	}

	// Deterministic ordering by ID
	sort.Slice(problems, func(i, j int) bool {
		return problems[i].ID < problems[j].ID
	})

	c.JSON(http.StatusOK, gin.H{"problems": problems})
}

// GetProblem handles GET /api/problems/:id
func GetProblem(c *gin.Context) {
	id := c.Param("id")
	if !problemIDPattern.MatchString(id) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid problem id"})
		return
	}

	problemsPath := getProblemsPath()
	problemPath := filepath.Join(problemsPath, id)
	
	// Check if directory exists
	info, err := os.Stat(problemPath)
	if err != nil || !info.IsDir() {
		c.JSON(http.StatusNotFound, gin.H{"error": "problem not found"})
		return
	}

	initPath := filepath.Join(problemPath, "init.yml")
	prob, err := parseProblem(id, initPath)
	if err != nil && !errors.Is(err, os.ErrNotExist) {
		log.Printf("[ERROR] failed to parse problem %s: %v", id, err)
	}
	
	// Override description if problem.md exists
	mdPath := filepath.Join(problemPath, "problem.md")
	mdData, err := os.ReadFile(mdPath)
	if err == nil {
		prob.Description = string(mdData)
	}

	c.JSON(http.StatusOK, prob)
}
