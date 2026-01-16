FROM golang:1.21-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /build
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o backstage-gen .

FROM alpine:3.19
RUN apk add --no-cache ca-certificates git

COPY --from=builder /build/backstage-gen /usr/local/bin/backstage-gen

WORKDIR /workspace
ENTRYPOINT ["/usr/local/bin/backstage-gen"]
CMD ["--help"]
