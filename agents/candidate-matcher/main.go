package main

import (
	"context"
	"fmt"

	"github.com/MunyTa/Lab-13/agents/common"
)

type matcher struct{}

func (matcher) Process(_ context.Context, task common.Task) (common.Result, error) {
	score := common.OverlapScore(task.Candidate.Skills, task.Vacancy.RequiredSkills)
	if task.Candidate.Experience >= task.Vacancy.MinimumYears {
		score += 0.15
	}
	if score > 1 {
		score = 1
	}

	decision := "hold"
	if score >= 0.75 {
		decision = "invite"
	}

	if task.Context == nil {
		task.Context = map[string]any{}
	}
	task.Context["match_score"] = score
	task.Context["match_decision"] = decision

	return common.Result{
		TaskID:    task.ID,
		Stage:     task.Stage,
		Success:   true,
		Output:    fmt.Sprintf("Candidate matched with score %.2f", score),
		TraceID:   task.TraceID,
		NextStage: "interview_scheduler",
		Context:   task.Context,
		Metadata:  task.Metadata,
		Candidate: task.Candidate,
		Vacancy:   task.Vacancy,
		Score:     score,
		Decision:  decision,
	}, nil
}

func main() {
	ctx := context.Background()
	cfg := common.NewAgentConfig("candidate_matcher")
	if err := common.RunAgent(ctx, cfg, matcher{}); err != nil {
		panic(err)
	}
}
