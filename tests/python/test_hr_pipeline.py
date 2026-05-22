from orchestrator.domain.auction import AgentBid, TaskAuction
from orchestrator.domain.pipeline import PipelineBuilder


def test_pipeline_builder_creates_full_hr_flow():
    pipeline = PipelineBuilder().build_hr_screening_pipeline()

    assert pipeline == [
        "resume_parser",
        "candidate_matcher",
        "interview_scheduler",
        "feedback_agent",
        "llm_feedback_agent",
    ]


def test_auction_prefers_highest_skill_lowest_cost_and_available_agent():
    bids = [
        AgentBid(agent_id="matcher-a", cost=4, skill_score=0.74, available=True),
        AgentBid(agent_id="matcher-b", cost=3, skill_score=0.91, available=True),
        AgentBid(agent_id="matcher-c", cost=1, skill_score=0.95, available=False),
    ]

    winner = TaskAuction().select_winner(bids)

    assert winner.agent_id == "matcher-b"

