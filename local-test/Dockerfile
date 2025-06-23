FROM golang:1.22-alpine AS builder

RUN apk add --no-cache ca-certificates build-base git

WORKDIR /app

# Copy source code
COPY . .

# Initialize modules if needed and build
RUN go mod init verana || true
RUN go mod tidy
RUN CGO_ENABLED=0 go build -o veranad ./cmd/veranad

# Final image
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /app/veranad /usr/local/bin/veranad
EXPOSE 26656 26657 1317 9090
ENTRYPOINT ["veranad"]