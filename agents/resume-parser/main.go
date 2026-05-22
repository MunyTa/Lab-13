package main

import (
	"context"
	"fmt"
	"strings"

	"github.com/MunyTa/Lab-13/agents/common"
)

type resumeParser struct {
	store common.StateStore
}

func (parser resumeParser) Process(ctx context.Context, task common.Task) (common.Result, error) {
	stats, err := parser.store.Load(ctx, "resume-parser")
	if err != nil {
		return common.Result{}, err
	}

	stats.ProcessedCount++
	stats.CachedCandidates++
	if err := parser.store.Save(ctx, "resume-parser", stats); err != nil {
		return common.Result{}, err
	}

	task.Candidate.Summary = strings.TrimSpace(task.Candidate.Summary)
	if task.Candidate.Summary == "" {
		task.Candidate.Summary = fmt.Sprintf(
			"%s: %d years of experience in %s",
			task.Candidate.Name,
			task.Candidate.Experience,
			strings.Join(task.Candidate.Skills, ", "),
		)
	}

	if task.Context == nil {
		task.Context = map[string]any{}
	}
	task.Context["resume_stats"] = stats
	task.Context["parsed_resume"] = task.Candidate.Summary

	return common.Result{
		TaskID:    task.ID,
		Stage:     task.Stage,
		Success:   true,
		Output:    "Resume parsed successfully",
		TraceID:   task.TraceID,
		NextStage: "candidate_matcher",
		Context:   task.Context,
		Metadata:  task.Metadata,
		Candidate: task.Candidate,
		Vacancy:   task.Vacancy,
	}, nil
}

func main() {
	ctx := context.Background()
	cfg := common.NewAgentConfig("resume_parser")
	if err := common.RunAgent(ctx, cfg, resumeParser{store: common.NewMemoryStateStore()}); err != nil {
		panic(err)
	}
}
