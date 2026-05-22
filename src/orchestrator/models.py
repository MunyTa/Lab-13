from __future__ import annotations

from typing import Any
from uuid import uuid4

from pydantic import BaseModel, Field


class CandidatePayload(BaseModel):
    name: str
    email: str
    skills: list[str] = Field(default_factory=list)
    experience: int = 0
    summary: str = ""


class VacancyPayload(BaseModel):
    title: str
    required_skills: list[str] = Field(default_factory=list)
    minimum_years: int = 0


class TaskSubmission(BaseModel):
    candidate: CandidatePayload
    vacancy: VacancyPayload
    metadata: dict[str, str] = Field(default_factory=dict)


class TaskExecutionResult(BaseModel):
    task_id: str = Field(default_factory=lambda: str(uuid4()))
    trace_id: str = Field(default_factory=lambda: str(uuid4()))
    pipeline: list[str]
    winner: str
    status: str
    score: float = 0
    decision: str = ""
    scheduled: str = ""
    feedback: str = ""
    stage_results: list[dict[str, Any]] = Field(default_factory=list)


class MonitoringSnapshot(BaseModel):
    active_agents: dict[str, int]
    queue_depth: int
    processed_tasks: int
    last_task_id: str = ""

