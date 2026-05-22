from __future__ import annotations

from dataclasses import dataclass


@dataclass(slots=True)
class ScaleDecision:
    agent_name: str
    current_replicas: int
    desired_replicas: int
    reason: str


class QueueScaler:
    def __init__(self, threshold: int = 5, max_replicas: int = 4) -> None:
        self.threshold = threshold
        self.max_replicas = max_replicas

    def plan(self, agent_name: str, queue_depth: int, current_replicas: int) -> ScaleDecision:
        if queue_depth <= self.threshold:
            return ScaleDecision(
                agent_name=agent_name,
                current_replicas=current_replicas,
                desired_replicas=current_replicas,
                reason="Queue depth is within the safe threshold.",
            )

        desired = min(self.max_replicas, current_replicas + 1)
        return ScaleDecision(
            agent_name=agent_name,
            current_replicas=current_replicas,
            desired_replicas=desired,
            reason="Queue depth exceeded the scaling threshold.",
        )

