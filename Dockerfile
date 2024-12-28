#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/app
COPY . .
RUN go get -v ./...
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /go/bin/app -v .


#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/app /app
ENTRYPOINT /app
LABEL Name=misa-saa-ngapi Version=0.0.1
# RUN apk add --update nodejs npm
# RUN node --version
# RUN npm --version
