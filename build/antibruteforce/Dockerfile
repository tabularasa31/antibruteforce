# Step 1: Modules caching
FROM golang:latest as modules
COPY go.mod go.sum /modules/
WORKDIR /modules
RUN go mod download

# Step 2: Builder
FROM modules as builder
COPY --from=modules /go/pkg /go/pkg
COPY .. /app
WORKDIR /app

RUN CGO_ENABLED=0 go build \
        -o /bin/app ./cmd/main.go

# Step 3: Final
FROM scratch

LABEL service="Antibruteforce"
LABEL mainteiner="tabularasa31@gmail.com"

COPY --from=builder /app/config /config
COPY --from=builder /bin/app /app

CMD ["/app"]