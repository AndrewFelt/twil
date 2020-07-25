FROM golang:alpine AS build_base

RUN apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/twil

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .

RUN go mod download

COPY collector.go .
COPY main.go .

# Unit tests
# Unit tests currently do not exist
# RUN CGO_ENABLED=0 go test -v

# Build the Go app
RUN go build -o ./out/twil .

# Start fresh from a smaller image
FROM alpine:latest
RUN apk add ca-certificates

COPY --from=build_base /tmp/twil/out/twil /app/twil

# This container exposes port 8888 to the outside world
EXPOSE 8888

# Run the binary program produced by `go build`
CMD ["/app/twil", "-port=:8888", "-account=<twil_account_key>", "-token=<twil_token>"]