package handler

import (
	"crypto/rand"
	"encoding/hex"
	"log"
	"net/http"

	"github.com/GCET-Open-Source-Foundation/coding_arena/backend/adapter"
	"github.com/GCET-Open-Source-Foundation/coding_arena/backend/model"
	"github.com/gin-gonic/gin"
)

// Submit handles POST /submit.
func Submit(c *gin.Context) {
	var req model.SubmitRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "missing or invalid fields: code, language, and problem_id are required",
		})
		return
	}

	if err := validateInput(req.Code, req.Language, req.ProblemID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	id := generateID("sub_")

	log.Printf("[INFO] submission received: id=%s lang=%s problem=%s ip=%s",
		id, req.Language, req.ProblemID, c.ClientIP())

	if judgeAdapter == nil || !judgeAdapter.Available() {
		c.JSON(http.StatusAccepted, model.SubmitResponse{
			ID:        id,
			ProblemID: req.ProblemID,
			Language:  req.Language,
			Status:    "queued",
			Message:   "submission received, pending judge execution",
		})
		return
	}

	result, err := judgeAdapter.Submit(adapter.SubmissionRequest{
		ProblemID:    req.ProblemID,
		Language:     req.Language,
		Source:       req.Code,
		ShortCircuit: false,
	})
	if err != nil {
		log.Printf("[ERROR] judge failed for %s: %v", id, err)
		c.JSON(http.StatusAccepted, model.SubmitResponse{
			ID:        id,
			ProblemID: req.ProblemID,
			Language:  req.Language,
			Status:    "judge_error",
			Message:   "judge grading failed, submission queued for retry",
		})
		return
	}

	resp := model.SubmitResponse{
		ID:        id,
		ProblemID: req.ProblemID,
		Language:  req.Language,
		Status:    "graded",
		Message:   result.Status,
		Result:    mapJudgeResult(result),
	}

	log.Printf("[INFO] graded: id=%s verdict=%s points=%.1f/%.1f",
		id, result.Status, result.Points, result.TotalPoints)

	c.JSON(http.StatusOK, resp)
}

func generateID(prefix string) string {
	b := make([]byte, 16)
	rand.Read(b)
	return prefix + hex.EncodeToString(b)
}

func mapJudgeResult(r *adapter.SubmissionResult) *model.JudgeResult {
	jr := &model.JudgeResult{
		Verdict:      r.Status,
		CompileError: r.CompileError,
		TotalTime:    r.TotalTime,
		MaxMemory:    r.MaxMemory,
		Points:       r.Points,
		TotalPoints:  r.TotalPoints,
	}
	for _, c := range r.Cases {
		jr.Cases = append(jr.Cases, model.JudgeCaseResult{
			Position: c.Position,
			Status:   c.Status,
			Time:     c.Time,
			Memory:   c.Memory,
			Points:   c.Points,
			Total:    c.Total,
			Feedback: c.Feedback,
		})
	}
	return jr
}
