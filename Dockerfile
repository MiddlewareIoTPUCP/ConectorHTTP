### Builder
FROM golang:1.14 as builder

WORKDIR /app

# Copying go modules files
COPY go.mod .
COPY go.sum .

# Downloading all dependencies
RUN go mod download

# Copying application files
COPY . .

# Building application
RUN go build -o /app/bin/ConectorHTTP

### Runner
FROM gcr.io/distroless/base

COPY --from=builder /app/bin/ConectorHTTP /ConectorHTTP

CMD ["/ConectorHTTP"]
