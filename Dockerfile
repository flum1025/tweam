FROM golang:1.14 as builder

COPY go.* /src/
WORKDIR /src
RUN go mod download

ADD . /src

ENV CGO_ENABLED=0
RUN go build -o scheduler cmd/scheduler/main.go
RUN go build -o worker cmd/worker/main.go

# ------

FROM alpine:latest

COPY --from=builder /src/scheduler /usr/local/bin/scheduler
COPY --from=builder /src/worker /usr/local/bin/worker

EXPOSE 3000
