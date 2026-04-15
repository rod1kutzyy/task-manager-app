FROM golang:1.26.2-bookworm AS builder

WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /app/cmd/notesapp/exe /app/cmd/notesapp


FROM alpine:3.23

WORKDIR /app

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

COPY --from=builder /app/cmd/notesapp/exe /app

COPY --from=builder /app/public /app/public

RUN chown -R appuser:appgroup /app

USER appuser

ENV PROJECT_ROOT=/app

EXPOSE 5050

CMD [ "/app/exe" ]
