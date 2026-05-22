package main

import (
	"context"
	"fmt"

	"github.com/MunyTa/Lab-13/agents/common"
)

type feedbackAgent struct{}

func (feedbackAgent) Process(_ context.Context, task common.Task) (common.Result, error) {
	feedback := fmt.Sprintf(
		"Candidate %s is moving forward for %s. Interview at %s.",
		task.Candidate.Name,
		task.Vacancy.Title,
		task.Scheduled,
	)

	if task.Context == nil {
		task.Context = map[string]any{}
	}
	task.Context["feedback_message"] = feedback

	return common.Result{
		TaskID:    task.ID,
		Stage:     task.Stage,
		Success:   true,
		Output:    "Feedback prepared",
		TraceID:   task.TraceID,
		NextStage: "llm_feedback_agent",
		Context:   task.Context,
		Metadata:  task.Metadata,
		Candidate: task.Candidate,
		Vacancy:   task.Vacancy,
		Score:     task.Score,
		Decision:  task.Decision,
		Scheduled: task.Scheduled,
		Feedback:  feedback,
	}, nil
}

func main() {
	ctx := context.Background()
	cfg := common.NewAgentConfig("feedback_agent")
	if err := common.RunAgent(ctx, cfg, feedbackAgent{}); err != nil {
		panic(err)
	}
}
