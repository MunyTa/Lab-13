FROM golang:1.26-alpine AS build

WORKDIR /app

ARG AGENT_DIR

COPY go.mod .
COPY agents ./agents

RUN go mod tidy
RUN CGO_ENABLED=0 go build -o /bin/agent ${AGENT_DIR}

FROM alpine:3.22

WORKDIR /app

COPY --from=build /bin/agent /bin/agent

CMD ["/bin/agent"]
