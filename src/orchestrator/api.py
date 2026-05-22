from __future__ import annotations

from fastapi import FastAPI

from orchestrator.models import TaskSubmission
from orchestrator.services.monitoring import MonitoringService
from orchestrator.services.orchestrator import HROrchestrator
from orchestrator.services.scaler import QueueScaler

app = FastAPI(title="Lab 13 HR Multi-Agent System")
orchestrator = HROrchestrator()
monitoring = MonitoringService()
scaler = QueueScaler()


@app.get("/health")
def healthcheck() -> dict[str, str]:
    return {"status": "ok"}


@app.post("/tasks")
def submit_task(payload: TaskSubmission) -> dict:
    result = orchestrator.execute(payload)
    monitoring.record_task(result.task_id)
    return result.model_dump()


@app.get("/monitoring")
def monitoring_snapshot() -> dict:
    snapshot = monitoring.snapshot(
        active_agents={
            "resume_parser": 1,
            "candidate_matcher": 1,
            "interview_scheduler": 1,
            "feedback_agent": 1,
            "llm_feedback_agent": 1,
        },
        queue_depth=0,
    )
    return snapshot.model_dump()


@app.get("/scaling/{agent_name}")
def scaling_plan(agent_name: str, queue_depth: int = 0, current_replicas: int = 1) -> dict:
    decision = scaler.plan(agent_name, queue_depth, current_replicas)
    return decision.__dict__

