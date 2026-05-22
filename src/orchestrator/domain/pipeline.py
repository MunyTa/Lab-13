class PipelineBuilder:
    def build_hr_screening_pipeline(self) -> list[str]:
        return [
            "resume_parser",
            "candidate_matcher",
            "interview_scheduler",
            "feedback_agent",
            "llm_feedback_agent",
        ]

