# ─── Stage 1: Build Vue frontend ──────────────────────────────────
FROM node:20-alpine AS frontend-builder

WORKDIR /app/frontend

COPY frontend/package*.json ./
RUN npm ci --silent

COPY frontend/ ./
RUN npm run build          # outputs to ../static via vite.config.js outDir

# ─── Stage 2: Build Go binary ─────────────────────────────────────
FROM golang:1.24-alpine AS go-builder
ENV GOTOOLCHAIN=auto

WORKDIR /app

# Cache dependency downloads separately from source
COPY go.mod go.sum ./
RUN go mod download

COPY main.go ./
# Copy built frontend assets (needed only at runtime, but embedded via Static route)
COPY --from=frontend-builder /app/static ./static

RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o shortenurl main.go

# ─── Stage 3: Minimal runtime image ───────────────────────────────
FROM gcr.io/distroless/static-debian12:nonroot

WORKDIR /app

COPY --from=go-builder /app/shortenurl ./
COPY --from=go-builder /app/static     ./static

EXPOSE 8080

ENTRYPOINT ["/app/shortenurl"]
