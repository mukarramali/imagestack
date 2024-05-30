# Build stage
ARG GO_VERSION=1.22.3
FROM golang:${GO_VERSION}-alpine AS build

# Install dependencies for Alpine
RUN apk add --no-cache gcc musl-dev

# Copy library
WORKDIR /lib
COPY lib /lib

WORKDIR /src

# Download dependencies
COPY services/api/go.mod services/api/go.sum ./
RUN go mod tidy

# Copy the source code
COPY services/api .

# Build the application
RUN CGO_ENABLED=1 go build -o /bin/server .

# Final stage
FROM alpine:latest AS final

# Install runtime dependencies
RUN apk add --no-cache \
    ca-certificates \
    tzdata

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