FROM golang:1.19-buster as builder
WORKDIR /app

COPY . ./
RUN go mod download

RUN go build ./cmd/main.go


FROM debian:buster-slim
RUN set -x && apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y \
    ca-certificates && \
    rm -rf /var/lib/apt/lists/*

# Copy the binary to the production image from the builder stage.
COPY --from=builder /app/main /app/server
ADD /configs /configs

# Run the web service on container startup.
ENTRYPOINT ["/app/server"]