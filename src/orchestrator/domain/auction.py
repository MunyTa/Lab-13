from dataclasses import dataclass


@dataclass(slots=True)
class AgentBid:
    agent_id: str
    cost: int
    skill_score: float
    available: bool


class TaskAuction:
    def select_winner(self, bids: list[AgentBid]) -> AgentBid:
        if not bids:
            raise ValueError("At least one bid is required")

        return sorted(
            bids,
            key=lambda bid: (
                not bid.available,
                -bid.skill_score,
                bid.cost,
                bid.agent_id,
            ),
        )[0]

