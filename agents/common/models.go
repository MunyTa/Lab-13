package common

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/nats-io/nats.go"
)

type Task struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Stage      string                 `json:"stage"`
	ReplyTo    string                 `json:"reply_to"`
	TraceID    string                 `json:"trace_id"`
	Candidate  CandidateProfile       `json:"candidate"`
	Vacancy    VacancyProfile         `json:"vacancy"`
	Context    map[string]any         `json:"context"`
	Metadata   map[string]string      `json:"metadata"`
	Pipeline   []string               `json:"pipeline"`
	Score      float64                `json:"score"`
	Decision   string                 `json:"decision"`
	Feedback   string                 `json:"feedback"`
	Scheduled  string                 `json:"scheduled"`
	Attributes map[string]interface{} `json:"attributes"`
}

type CandidateProfile struct {
	Name       string   `json:"name"`
	Email      string   `json:"email"`
	Skills     []string `json:"skills"`
	Experience int      `json:"experience"`
	Summary    string   `json:"summary"`
}

type VacancyProfile struct {
	Title          string   `json:"title"`
	RequiredSkills []string `json:"required_skills"`
	MinimumYears   int      `json:"minimum_years"`
}

type Result struct {
	TaskID     string                 `json:"task_id"`
	Stage      string                 `json:"stage"`
	Success    bool                   `json:"success"`
	Output     string                 `json:"output"`
	TraceID    string                 `json:"trace_id"`
	NextStage  string                 `json:"next_stage,omitempty"`
	Context    map[string]any         `json:"context,omitempty"`
	Metadata   map[string]string      `json:"metadata,omitempty"`
	Candidate  CandidateProfile       `json:"candidate"`
	Vacancy    VacancyProfile         `json:"vacancy"`
	Score      float64                `json:"score,omitempty"`
	Decision   string                 `json:"decision,omitempty"`
	Feedback   string                 `json:"feedback,omitempty"`
	Scheduled  string                 `json:"scheduled,omitempty"`
	Attributes map[string]interface{} `json:"attributes,omitempty"`
}

type Bid struct {
	AgentID    string  `json:"agent_id"`
	Cost       int     `json:"cost"`
	SkillScore float64 `json:"skill_score"`
	Available  bool    `json:"available"`
}

func SelectWinningBid(bids []Bid) (Bid, error) {
	if len(bids) == 0 {
		return Bid{}, errors.New("no bids provided")
	}

	winner := bids[0]
	for _, bid := range bids[1:] {
		if bid.Available != winner.Available {
			if bid.Available {
				winner = bid
			}
			continue
		}

		if bid.SkillScore > winner.SkillScore {
			winner = bid
			continue
		}

		if bid.SkillScore == winner.SkillScore && bid.Cost < winner.Cost {
			winner = bid
		}
	}

	return winner, nil
}

type ResumeStats struct {
	ProcessedCount   int `json:"processed_count"`
	CachedCandidates int `json:"cached_candidates"`
}

func (stats ResumeStats) Restore() ResumeStats {
	return ResumeStats{
		ProcessedCount:   stats.ProcessedCount,
		CachedCandidates: stats.CachedCandidates,
	}
}

type StateStore interface {
	Load(ctx context.Context, key string) (ResumeStats, error)
	Save(ctx context.Context, key string, stats ResumeStats) error
}

type MemoryStateStore struct {
	mu    sync.Mutex
	items map[string]ResumeStats
}

func NewMemoryStateStore() *MemoryStateStore {
	return &MemoryStateStore{
		items: make(map[string]ResumeStats),
	}
}

func (store *MemoryStateStore) Load(_ context.Context, key string) (ResumeStats, error) {
	store.mu.Lock()
	defer store.mu.Unlock()

	stats, ok := store.items[key]
	if !ok {
		return ResumeStats{}, nil
	}

	return stats.Restore(), nil
}

func (store *MemoryStateStore) Save(_ context.Context, key string, stats ResumeStats) error {
	store.mu.Lock()
	defer store.mu.Unlock()

	store.items[key] = stats.Restore()
	return nil
}

type AgentProcessor interface {
	Process(ctx context.Context, task Task) (Result, error)
}

type AgentConfig struct {
	Name          string
	Subject       string
	DefaultReply  string
	NATSURL       string
	LogFile       string
	MetricLogFile string
}

func NewAgentConfig(name string) AgentConfig {
	safeName := strings.ReplaceAll(name, " ", "-")

	return AgentConfig{
		Name:          name,
		Subject:       envOrDefault("AGENT_SUBJECT", "tasks."+safeName),
		DefaultReply:  envOrDefault("RESULT_SUBJECT", "tasks.completed"),
		NATSURL:       envOrDefault("NATS_URL", nats.DefaultURL),
		LogFile:       envOrDefault("AGENT_LOG_FILE", fmt.Sprintf("logs/%s.log", safeName)),
		MetricLogFile: envOrDefault("AGENT_METRIC_FILE", fmt.Sprintf("logs/%s-metrics.log", safeName)),
	}
}

func RunAgent(ctx context.Context, cfg AgentConfig, processor AgentProcessor) error {
	logger, metricLogger, closeLogs, err := buildLoggers(cfg)
	if err != nil {
		return err
	}
	defer closeLogs()

	nc, err := nats.Connect(cfg.NATSURL)
	if err != nil {
		return fmt.Errorf("connect nats: %w", err)
	}
	defer nc.Close()

	var processed int

	_, err = nc.QueueSubscribe(cfg.Subject, cfg.Name+"-workers", func(msg *nats.Msg) {
		var task Task
		if err := json.Unmarshal(msg.Data, &task); err != nil {
			logger.Error("failed to decode task", slog.String("error", err.Error()))
			return
		}

		result, err := processor.Process(ctx, task)
		if err != nil {
			logger.Error("task processing failed", slog.String("task_id", task.ID), slog.String("error", err.Error()))
			result = Result{
				TaskID:    task.ID,
				Stage:     task.Stage,
				Success:   false,
				Output:    err.Error(),
				TraceID:   task.TraceID,
				Candidate: task.Candidate,
				Vacancy:   task.Vacancy,
				Context:   task.Context,
				Metadata:  task.Metadata,
			}
		}

		if result.TaskID == "" {
			result.TaskID = task.ID
		}
		if result.Stage == "" {
			result.Stage = task.Stage
		}
		if result.TraceID == "" {
			result.TraceID = task.TraceID
		}
		if result.Metadata == nil {
			result.Metadata = task.Metadata
		}
		if result.Context == nil {
			result.Context = task.Context
		}

		payload, err := json.Marshal(result)
		if err != nil {
			logger.Error("failed to encode result", slog.String("task_id", task.ID), slog.String("error", err.Error()))
			return
		}

		replyTo := task.ReplyTo
		if replyTo == "" {
			replyTo = cfg.DefaultReply
		}

		if err := nc.Publish(replyTo, payload); err != nil {
			logger.Error("failed to publish result", slog.String("task_id", task.ID), slog.String("error", err.Error()))
			return
		}

		processed++
		metricLogger.Info(
			"agent_metrics",
			slog.String("agent", cfg.Name),
			slog.Int("processed_tasks", processed),
			slog.String("task_id", task.ID),
			slog.Bool("success", result.Success),
		)
	})
	if err != nil {
		return fmt.Errorf("subscribe: %w", err)
	}

	logger.Info("agent started", slog.String("name", cfg.Name), slog.String("subject", cfg.Subject))
	<-ctx.Done()
	logger.Info("agent stopping", slog.String("name", cfg.Name))

	return nil
}

func envOrDefault(key, fallback string) string {
	value := strings.TrimSpace(os.Getenv(key))
	if value == "" {
		return fallback
	}
	return value
}

func buildLoggers(cfg AgentConfig) (*slog.Logger, *slog.Logger, func(), error) {
	if err := os.MkdirAll("logs", 0o755); err != nil {
		return nil, nil, nil, err
	}

	logFile, err := os.OpenFile(cfg.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		return nil, nil, nil, err
	}

	metricFile, err := os.OpenFile(cfg.MetricLogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0o644)
	if err != nil {
		_ = logFile.Close()
		return nil, nil, nil, err
	}

	logger := slog.New(slog.NewTextHandler(logFile, &slog.HandlerOptions{Level: slog.LevelInfo}))
	metricLogger := slog.New(slog.NewTextHandler(metricFile, &slog.HandlerOptions{Level: slog.LevelInfo}))

	closeLogs := func() {
		_ = logFile.Close()
		_ = metricFile.Close()
	}

	return logger, metricLogger, closeLogs, nil
}

func OverlapScore(left []string, right []string) float64 {
	if len(right) == 0 {
		return 1
	}

	index := make(map[string]struct{}, len(left))
	for _, item := range left {
		index[strings.ToLower(strings.TrimSpace(item))] = struct{}{}
	}

	var matched int
	for _, item := range right {
		if _, ok := index[strings.ToLower(strings.TrimSpace(item))]; ok {
			matched++
		}
	}

	return float64(matched) / float64(len(right))
}

func DefaultInterviewSlot() string {
	return time.Now().UTC().Add(48 * time.Hour).Format(time.RFC3339)
}
