FROM golang:1.20-alpine AS build_base

RUN apk update && apk add --no-cache git

# Set the Current Working Directory inside the container
WORKDIR /tmp/gochat

# We want to populate the module cache based on the go.{mod,sum} files.
COPY go.mod .
COPY go.sum .

RUN go mod download && go mod verify

COPY . .

# Build the Go app
RUN CGO_ENABLED=0 go build -o ./out/gochat ./cmd/main.go

# Start fresh from a smaller image
FROM gcr.io/distroless/static

COPY --from=build_base /tmp/gochat/out/gochat /app/gochat

# Expose the container to the outside world
EXPOSE 3000

# Run the program
CMD ["/app/gochat"]
