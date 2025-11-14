FROM golang:1.22 as build
WORKDIR /app
COPY go.mod .
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o watcher ./cmd/partition_watcher

FROM gcr.io/distroless/base-debian12
WORKDIR /
COPY --from=build /app/watcher /watcher
EXPOSE 9108
ENTRYPOINT ["/watcher"]
