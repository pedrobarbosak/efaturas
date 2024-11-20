##############################
#   Build Container
FROM golang:1.23.2 as build

ENV CGO_ENABLED=0
WORKDIR /app

COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy all the code
COPY . .

# Build
RUN go build -o bin .

##############################
#   Run Container
FROM alpine:latest as run

WORKDIR /app

COPY --from=build /app/bin .

CMD ["./bin"]