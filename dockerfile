FROM oven/bun:1 AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package.json frontend/bun.lockb* ./
RUN bun install --frozen-lockfile

COPY frontend/ ./
RUN bun run build

FROM golang:1.23-alpine AS backend-builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

COPY --from=frontend-builder /app/frontend/dist ./frontend/dist

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /app/server .


FROM alpine:3.20

RUN apk add --no-cache ca-certificates

WORKDIR /app
COPY --from=backend-builder /app/server .

EXPOSE 3000

ENTRYPOINT ["/app/server"]