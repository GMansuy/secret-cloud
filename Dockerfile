# Build stage
FROM golang:1.22-alpine AS builder

RUN apk add --no-cache git ca-certificates tzdata

COPY . /app/

WORKDIR /app
RUN go mod download

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o rocket-app cmd/main.go

# Final stage
FROM alpine:latest

# Add certificates and timezone data
RUN apk --no-cache add ca-certificates tzdata curl

# Create non-root user
RUN adduser -D -g '' appuser

RUN curl -sLO "https://dl.k8s.io/release/$(curl -Ls https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl" \
    && install -o root -g root -m 0755 kubectl /usr/local/bin/kubectl \
    && rm kubectl

WORKDIR /app

# Copy binary from build stage
COPY --from=builder /app/rocket-app .

# Expose the API port
EXPOSE 8080

# Run the application
CMD ["./rocket-app"]