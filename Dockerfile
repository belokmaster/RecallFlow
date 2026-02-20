# syntax=docker/dockerfile:1

FROM node:20-alpine AS web-build
WORKDIR /src
COPY web/frontend/package.json web/frontend/package-lock.json ./web/frontend/
RUN cd web/frontend && npm ci
COPY web/frontend ./web/frontend
RUN cd web/frontend && npm run build

FROM golang:1.22-alpine AS go-build
WORKDIR /src
RUN apk add --no-cache build-base sqlite-dev
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.20
WORKDIR /app
RUN apk add --no-cache ca-certificates sqlite-libs
COPY --from=go-build /out/server ./server
COPY --from=web-build /src/web/frontend/dist ./web/frontend/dist
ENV RECALL_FLOW_DB=/data/recall_flow.db
EXPOSE 8080
CMD ["./server"]
