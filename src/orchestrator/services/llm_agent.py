from __future__ import annotations

from orchestrator.models import CandidatePayload, VacancyPayload


class LLMFeedbackAgent:
    def generate_feedback(
        self,
        candidate: CandidatePayload,
        vacancy: VacancyPayload,
        decision: str,
        scheduled: str,
    ) -> str:
        if decision != "invite":
            return (
                f"Candidate {candidate.name} is placed on hold for {vacancy.title}. "
                "Recommend collecting more signals before sending final feedback."
            )

        top_skills = ", ".join(candidate.skills[:3]) or "relevant experience"
        return (
            f"Candidate {candidate.name} is a strong fit for {vacancy.title}. "
            f"Key strengths: {top_skills}. "
            f"Interview slot reserved for {scheduled}."
        )

