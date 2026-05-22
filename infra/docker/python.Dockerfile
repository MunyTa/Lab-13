FROM python:3.12-slim

WORKDIR /app

COPY pyproject.toml .
COPY src ./src

RUN pip install --no-cache-dir .

CMD ["python", "-m", "uvicorn", "orchestrator.api:app", "--host", "0.0.0.0", "--port", "8000"]

