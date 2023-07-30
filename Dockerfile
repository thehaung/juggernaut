# IMAGE: internal-platform_tool-oauth-gsuite
FROM golang:1.20-alpine as builder

# Enable go modules
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=arm64

# Install git. (alpine image does not have git in it)
RUN apk update \
    && apk add --no-cache git \
    && apk upgrade \
    && apk add --no-cache \
    ca-certificates \
    && update-ca-certificates \
    && rm -rf /var/cache/apk/*

# Set current working directory
WORKDIR /app

# Note here: To avoid downloading dependencies every time we
# build image. Here, we are caching all the dependencies by
# first copying go.mod and go.sum files and downloading them,
# to be used every time we build the image if the dependencies
# are not changed.

# Copy go mod and sum files
COPY go.mod ./
COPY go.sum ./

# Download all dependencies.
RUN go mod download

# Now, copy the source code
COPY . .

# Note here: CGO_ENABLED is disabled for cross system compilation
# It is also a common best practise.

# Build the application.
RUN go build -o ./bin/main ./cmd

# Finally our multi-stage to build a small image
# Start a new stage from scratch
FROM scratch

ENV APP_ENV=production

# Copy the Pre-built binary file
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/bin/main .
COPY --from=builder /app/config/config.yml ./config/config.yml
COPY --from=builder /app/.env .env

# Run executable
CMD ["./main"]

EXPOSE 8081