# Step 1: Build env
FROM golang:latest as builder
ENV NAME "antibruteforce"

WORKDIR /app
ENV CGO_ENABLED=0
COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN go build -o bin/${NAME} ./cmd/${NAME}

# Step 2: Final
FROM alpine
ENV NAME "antibruteforce"

COPY --from=builder /app/config.yml /app/config.yml
COPY --from=builder /app/bin/${NAME} /app/

WORKDIR /app

EXPOSE 50051
CMD ./${NAME}
