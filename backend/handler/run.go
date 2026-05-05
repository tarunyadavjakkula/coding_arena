package handler

import (
	"fmt"
	"log"
	"net/http"

	"github.com/GCET-Open-Source-Foundation/coding_arena/backend/adapter"
	"github.com/GCET-Open-Source-Foundation/coding_arena/backend/model"
	"github.com/gin-gonic/gin"
)

// Run handles POST /run.
func Run(c *gin.Context) {
	var req model.RunRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing or invalid fields: source, language, and problem_id are required",
		})
		return
	}

	if err := validateInput(req.Source, req.Language, req.ProblemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	runID := generateID("run_")

	log.Printf("[INFO] run request: id=%s lang=%s problem=%s ip=%s",
		runID, req.Language, req.ProblemID, c.ClientIP())

	if judgeAdapter == nil || !judgeAdapter.Available() {
		c.JSON(http.StatusServiceUnavailable, model.RunResponse{
			RunID:   runID,
			Status:  "unavailable",
			Message: "no judge connected",
		})
		return
	}

	result, err := judgeAdapter.Submit(adapter.SubmissionRequest{
		ProblemID:    req.ProblemID,
		Language:     req.Language,
		Source:       req.Source,
		ShortCircuit: false,
	})
	if err != nil {
		log.Printf("[ERROR] run failed for %s: %v", runID, err)
		c.JSON(http.StatusInternalServerError, model.RunResponse{
			RunID:   runID,
			Status:  "error",
			Message: "judge execution failed",
		})
		return
	}

	testCases := make([]model.RunCaseResult, 0, len(result.Cases))
	for _, cr := range result.Cases {
		testCases = append(testCases, model.RunCaseResult{
			Name:         fmt.Sprintf("Test Case %d", cr.Position),
			Status:       cr.Status,
			Time:         cr.Time,
			MemoryKB:     cr.Memory,
			ActualOutput: cr.Feedback,
		})
	}

	c.JSON(http.StatusOK, model.RunResponse{
		RunID:     runID,
		Status:    result.Status,
		Message:   result.Status,
		TestCases: testCases,
	})
}
