# syntax=docker/dockerfile:1

# Build stage
ARG GO_VERSION=1.21.6
FROM golang:${GO_VERSION} AS build

# Install libvips
RUN apt-get update && apt-get install -y libvips-dev

WORKDIR /src

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -o /bin/server .

# Final stage
FROM debian:latest AS final

# Install runtime dependencies
RUN apt-get update && apt-get install -y \
    libvips \
    ca-certificates \
    tzdata \
    && update-ca-certificates && rm -rf /var/lib/apt/lists/*

# Create a non-privileged user
ARG UID=10001
RUN adduser \
    --disabled-password \
    --gecos "" \
    --home "/nonexistent" \
    --no-create-home \
    --uid "${UID}" \
    appuser

RUN mkdir -p /data/images && chown -R appuser /data/images

USER appuser

# Copy the executable from the build stage
COPY --from=build /bin/server /bin/

# Expose the port that the application listens on
EXPOSE 8080

# Run the application
ENTRYPOINT [ "/bin/server" ]