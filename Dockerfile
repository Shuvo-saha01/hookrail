FROM golang:1.26-alpine AS builder

WORKDIR /app

RUN go install github.com/air-verse/air@latest



COPY go.mod ./

COPY . .

RUN go build



# Build step for production
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /myapp .
# -------------------------------------------------------

FROM alpine:3.21

WORKDIR /

COPY --from=builder /myapp /myapp

EXPOSE 8080
ENTRYPOINT ["/myapp"]
