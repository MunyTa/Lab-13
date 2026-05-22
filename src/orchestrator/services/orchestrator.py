from __future__ import annotations

from datetime import UTC, datetime, timedelta

from orchestrator.domain.auction import AgentBid, TaskAuction
from orchestrator.domain.pipeline import PipelineBuilder
from orchestrator.models import TaskExecutionResult, TaskSubmission
from orchestrator.services.llm_agent import LLMFeedbackAgent
from orchestrator.services.monitoring import MonitoringService


class HROrchestrator:
    def __init__(self) -> None:
        self._pipeline_builder = PipelineBuilder()
        self._auction = TaskAuction()
        self._llm_feedback_agent = LLMFeedbackAgent()
        self._monitoring = MonitoringService()

    def execute(self, submission: TaskSubmission) -> TaskExecutionResult:
        pipeline = self._pipeline_builder.build_hr_screening_pipeline()
        bids = [
            AgentBid(agent_id="matcher-primary", cost=3, skill_score=0.92, available=True),
            AgentBid(agent_id="matcher-backup", cost=2, skill_score=0.81, available=True),
            AgentBid(agent_id="matcher-cold", cost=1, skill_score=0.95, available=False),
        ]
        winner = self._auction.select_winner(bids)

        required = {item.lower() for item in submission.vacancy.required_skills}
        candidate_skills = {item.lower() for item in submission.candidate.skills}
        matched = len(required & candidate_skills)
        score = matched / len(required) if required else 1.0
        if submission.candidate.experience >= submission.vacancy.minimum_years:
            score = min(1.0, score + 0.15)

        decision = "invite" if score >= 0.75 else "hold"
        scheduled = ""
        if decision == "invite":
            scheduled = (datetime.now(UTC) + timedelta(days=2)).isoformat()

        feedback = self._llm_feedback_agent.generate_feedback(
            candidate=submission.candidate,
            vacancy=submission.vacancy,
            decision=decision,
            scheduled=scheduled,
        )

        result = TaskExecutionResult(
            pipeline=pipeline,
            winner=winner.agent_id,
            status="completed",
            score=score,
            decision=decision,
            scheduled=scheduled,
            feedback=feedback,
            stage_results=[
                {"stage": "resume_parser", "status": "completed"},
                {"stage": "candidate_matcher", "status": "completed", "score": score},
                {"stage": "interview_scheduler", "status": "completed" if scheduled else "skipped"},
                {"stage": "feedback_agent", "status": "completed"},
                {"stage": "llm_feedback_agent", "status": "completed"},
            ],
        )
        self._monitoring.record_task(result.task_id)
        return result

    def monitoring_snapshot(self) -> MonitoringService:
        return self._monitoring

