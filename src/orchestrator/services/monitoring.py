from __future__ import annotations

from orchestrator.models import MonitoringSnapshot


class MonitoringService:
    def __init__(self) -> None:
        self._processed_tasks = 0
        self._last_task_id = ""

    def record_task(self, task_id: str) -> None:
        self._processed_tasks += 1
        self._last_task_id = task_id

    def snapshot(self, active_agents: dict[str, int], queue_depth: int) -> MonitoringSnapshot:
        return MonitoringSnapshot(
            active_agents=active_agents,
            queue_depth=queue_depth,
            processed_tasks=self._processed_tasks,
            last_task_id=self._last_task_id,
        )

