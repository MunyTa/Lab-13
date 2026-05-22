package common

import "testing"

func TestSelectWinningBidPrefersAvailabilityThenSkillThenCost(t *testing.T) {
	bids := []Bid{
		{AgentID: "matcher-a", Cost: 4, SkillScore: 0.74, Available: true},
		{AgentID: "matcher-b", Cost: 3, SkillScore: 0.91, Available: true},
		{AgentID: "matcher-c", Cost: 1, SkillScore: 0.95, Available: false},
	}

	winner, err := SelectWinningBid(bids)
	if err != nil {
		t.Fatalf("SelectWinningBid returned error: %v", err)
	}

	if winner.AgentID != "matcher-b" {
		t.Fatalf("expected matcher-b to win, got %s", winner.AgentID)
	}
}

func TestResumeStatsRestoreKeepsCountersAcrossRestart(t *testing.T) {
	stats := ResumeStats{ProcessedCount: 3, CachedCandidates: 7}

	restored := stats.Restore()

	if restored.ProcessedCount != 3 {
		t.Fatalf("expected processed count 3, got %d", restored.ProcessedCount)
	}

	if restored.CachedCandidates != 7 {
		t.Fatalf("expected cached candidates 7, got %d", restored.CachedCandidates)
	}
}
