package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestGetProblems(t *testing.T) {
	gin.SetMode(gin.TestMode)

	t.Run("successful listing", func(t *testing.T) {
		// Create a temporary directory for tests
		tempDir, err := os.MkdirTemp("", "judge-config-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Create problem subdirectories and init.yml files
		prob1Dir := filepath.Join(tempDir, "apple-picking")
		if err := os.Mkdir(prob1Dir, 0755); err != nil {
			t.Fatalf("failed to create problem directory %q: %v", prob1Dir, err)
		}
		if err := os.WriteFile(filepath.Join(prob1Dir, "init.yml"), []byte(`
id: apple-picking
title: Apple Picking
difficulty: medium
tags:
  - greedy
  - tree
`), 0644); err != nil {
			t.Fatalf("failed to write init.yml for %q: %v", prob1Dir, err)
		}

		prob2Dir := filepath.Join(tempDir, "zebra-puzzle")
		if err := os.Mkdir(prob2Dir, 0755); err != nil {
			t.Fatalf("failed to create problem directory %q: %v", prob2Dir, err)
		}
		if err := os.WriteFile(filepath.Join(prob2Dir, "init.yml"), []byte(`
name: Zebra Puzzle
`), 0644); err != nil {
			t.Fatalf("failed to write init.yml for %q: %v", prob2Dir, err)
		}

		prob3Dir := filepath.Join(tempDir, "no-yaml")
		if err := os.Mkdir(prob3Dir, 0755); err != nil {
			t.Fatalf("failed to create problem directory %q: %v", prob3Dir, err)
		}
		// no init.yml for this one, should fallback to inference

		// Override configPath
		originalConfigPath := configPath
		configPath = tempDir
		defer func() { configPath = originalConfigPath }()

		// Setup request
		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v", err)
		}

		if len(resp.Problems) != 3 {
			t.Fatalf("expected 3 problems, got %d", len(resp.Problems))
		}

		// They should be sorted alphabetically by ID: apple-picking, no-yaml, zebra-puzzle
		if resp.Problems[0].ID != "apple-picking" || resp.Problems[0].Title != "Apple Picking" {
			t.Errorf("unexpected first problem: %+v", resp.Problems[0])
		}
		if resp.Problems[0].Difficulty != "Medium" {
			t.Errorf("expected difficulty 'Medium', got '%s'", resp.Problems[0].Difficulty)
		}

		if resp.Problems[1].ID != "no-yaml" || resp.Problems[1].Title != "No Yaml" {
			t.Errorf("unexpected second problem (inferred): %+v", resp.Problems[1])
		}

		if resp.Problems[2].ID != "zebra-puzzle" || resp.Problems[2].Title != "Zebra Puzzle" {
			t.Errorf("unexpected third problem: %+v", resp.Problems[2])
		}
	})

	t.Run("empty directory", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "empty-dir")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		originalConfigPath := configPath
		configPath = tempDir
		defer func() { configPath = originalConfigPath }()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, w.Body.String())
		}

		if len(resp.Problems) != 0 {
			t.Errorf("expected empty problem list, got %d items", len(resp.Problems))
		}
	})

	t.Run("missing directory handling", func(t *testing.T) {
		originalConfigPath := configPath
		configPath = "/this/path/does/not/exist/12345"
		defer func() { configPath = originalConfigPath }()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		// Expected to handle gracefully and return empty list
		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, w.Body.String())
		}

		if len(resp.Problems) != 0 {
			t.Errorf("expected empty problem list for missing directory, got %d items", len(resp.Problems))
		}
	})

	t.Run("malformed config file skipped", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "malformed-dir")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		probDir := filepath.Join(tempDir, "bad-yaml")
		if err := os.Mkdir(probDir, 0755); err != nil {
			t.Fatalf("failed to create problem dir: %v", err)
		}
		// Write invalid YAML
		if err := os.WriteFile(filepath.Join(probDir, "init.yml"), []byte(`
id: bad-yaml
  invalid_indent: true
- broken_list
`), 0644); err != nil {
			t.Fatalf("failed to write malformed init.yml: %v", err)
		}

		originalConfigPath := configPath
		configPath = tempDir
		defer func() { configPath = originalConfigPath }()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, w.Body.String())
		}

		if len(resp.Problems) != 1 {
			t.Fatalf("expected 1 problem, got %d", len(resp.Problems))
		}

		// It should fallback to inference because YAML parsing failed
		if resp.Problems[0].ID != "bad-yaml" || resp.Problems[0].Title != "Bad Yaml" {
			t.Errorf("expected inferred problem, got %+v", resp.Problems[0])
		}
	})

	t.Run("invalid directory name skipped", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "invalid-dir-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		// Valid problem directory — should appear in results
		validDir := filepath.Join(tempDir, "valid-problem")
		if err := os.Mkdir(validDir, 0755); err != nil {
			t.Fatalf("failed to create valid problem dir: %v", err)
		}

		// Invalid directory names (uppercase, underscore) — must be skipped
		for _, name := range []string{"Invalid_Problem", "UPPERCASE", "has spaces"} {
			if err := os.Mkdir(filepath.Join(tempDir, name), 0755); err != nil {
				t.Fatalf("failed to create invalid dir %q: %v", name, err)
			}
		}

		originalConfigPath := configPath
		configPath = tempDir
		defer func() { configPath = originalConfigPath }()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		if w.Code != http.StatusOK {
			t.Errorf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, w.Body.String())
		}

		if len(resp.Problems) != 1 {
			t.Fatalf("expected 1 problem (invalid dirs skipped), got %d", len(resp.Problems))
		}
		if resp.Problems[0].ID != "valid-problem" {
			t.Errorf("expected valid-problem, got %q", resp.Problems[0].ID)
		}
	})

	t.Run("memory limit conversion", func(t *testing.T) {
		tempDir, err := os.MkdirTemp("", "memlimit-test")
		if err != nil {
			t.Fatalf("failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tempDir)

		probDir := filepath.Join(tempDir, "mem-problem")
		if err := os.Mkdir(probDir, 0755); err != nil {
			t.Fatalf("failed to create problem dir: %v", err)
		}
		if err := os.WriteFile(filepath.Join(probDir, "init.yml"), []byte(`
id: mem-problem
title: Memory Problem
difficulty: easy
memory_limit: 262144
time_limit: 2.0
`), 0644); err != nil {
			t.Fatalf("failed to write init.yml: %v", err)
		}

		originalConfigPath := configPath
		configPath = tempDir
		defer func() { configPath = originalConfigPath }()

		w := httptest.NewRecorder()
		c, _ := gin.CreateTestContext(w)
		c.Request, _ = http.NewRequest(http.MethodGet, "/api/problems", nil)

		GetProblems(c)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status 200, got %d", w.Code)
		}

		var resp struct {
			Problems []Problem `json:"problems"`
		}
		if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
			t.Fatalf("failed to unmarshal response: %v\nbody: %s", err, w.Body.String())
		}

		if len(resp.Problems) != 1 {
			t.Fatalf("expected 1 problem, got %d", len(resp.Problems))
		}

		// 262144 KB / 1024 = 256 MB
		if resp.Problems[0].MemoryLimit != 256 {
			t.Errorf("expected memory limit 256 MB, got %d", resp.Problems[0].MemoryLimit)
		}

		// 2.0 seconds = 2000 milliseconds
		if resp.Problems[0].TimeLimit != 2000 {
			t.Errorf("expected time limit 2000 ms, got %d", resp.Problems[0].TimeLimit)
		}
	})
}
