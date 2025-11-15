FROM golang:1.22 AS build
WORKDIR /app

# 1) bring in go.mod first (cacheable) and prefetch modules
COPY go.mod .
RUN go mod download

# 2) now copy the rest and let Go compute go.sum from imports
COPY . .
RUN go mod tidy

# 3) build static binary
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o watcher ./cmd/partition_watcher

FROM gcr.io/distroless/base-debian12
WORKDIR /app
COPY --from=build /app/watcher /usr/local/bin/watcher
USER nonroot:nonroot
EXPOSE 9108
ENTRYPOINT ["/usr/local/bin/watcher"]
