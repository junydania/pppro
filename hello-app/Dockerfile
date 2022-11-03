FROM golang:1.18-alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /build

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Build the Go app
RUN go build -ldflags="-s -w" -o helloapp cmd/app/main.go

# Start fresh from a smaller image
FROM alpine:3.16

RUN apk update && apk upgrade && apk add ca-certificates curl git bash \
iputils busybox-extras

COPY --from=build_base /build/helloapp /app/helloapp

# This container exposes port 8082 to the outside world
EXPOSE 8082

# Run the binary program produced by `go install`
CMD ["/app/helloapp"]