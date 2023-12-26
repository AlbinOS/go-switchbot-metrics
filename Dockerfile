# Build the application from source
FROM golang:1.21 AS build-stage

# Set destination for COPY
WORKDIR /app

# Download Go modules
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags "-s -w" -o /go-switchbot-influx

# Run the tests in the container
FROM build-stage AS run-test-stage
RUN go test -v ./...

# Deploy the application binary into a lean image
FROM gcr.io/distroless/static-debian12 AS build-release-stage

WORKDIR /

COPY --from=build-stage /go-switchbot-influx /go-switchbot-influx

EXPOSE 3000

USER nonroot:nonroot

ENTRYPOINT ["/go-switchbot-influx"]
