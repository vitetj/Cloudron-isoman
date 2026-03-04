FROM oven/bun:1 AS ui-builder

WORKDIR /build/ui

COPY ui/package.json ./
RUN bun install

COPY ui/ ./
RUN bun run build

FROM golang:1.24 AS backend-builder

ARG VERSION=dev

WORKDIR /build/backend

COPY backend/go.mod backend/go.sum ./
RUN go mod download

COPY backend/ ./

RUN CGO_ENABLED=0 GOOS=linux go build \
    -mod=mod \
    -ldflags="-w -s -X main.Version=${VERSION}" \
    -o /build/server .

FROM cloudron/base:5.0.0@sha256:04fd70dbd8ad6149c19de39e35718e024417c3e01dc9c6637eaf4a41ec4e596c

WORKDIR /app/code

COPY --from=backend-builder /build/server /app/code/server
COPY --from=backend-builder /build/backend/migrations /app/code/migrations
COPY --from=ui-builder /build/ui/dist /app/code/ui/dist
COPY start.sh /app/code/start.sh

RUN chmod +x /app/code/start.sh

EXPOSE 8080

CMD ["/app/code/start.sh"]
