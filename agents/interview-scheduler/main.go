package main

import (
	"context"
	"errors"
	"fmt"

	"github.com/MunyTa/Lab-13/agents/common"
)

type scheduler struct{}

func (scheduler) Process(_ context.Context, task common.Task) (common.Result, error) {
	if task.Decision == "hold" {
		return common.Result{}, errors.New("candidate is not eligible for interview scheduling")
	}

	slot := common.DefaultInterviewSlot()
	if task.Context == nil {
		task.Context = map[string]any{}
	}
	task.Context["scheduled_at"] = slot

	return common.Result{
		TaskID:    task.ID,
		Stage:     task.Stage,
		Success:   true,
		Output:    fmt.Sprintf("Interview scheduled at %s", slot),
		TraceID:   task.TraceID,
		NextStage: "feedback_agent",
		Context:   task.Context,
		Metadata:  task.Metadata,
		Candidate: task.Candidate,
		Vacancy:   task.Vacancy,
		Score:     task.Score,
		Decision:  task.Decision,
		Scheduled: slot,
	}, nil
}

func main() {
	ctx := context.Background()
	cfg := common.NewAgentConfig("interview_scheduler")
	if err := common.RunAgent(ctx, cfg, scheduler{}); err != nil {
		panic(err)
	}
}
