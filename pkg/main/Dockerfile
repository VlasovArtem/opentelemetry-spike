FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./pkg/main/*.go ./
COPY ./pkg/common/*.go ./pkg/common/
COPY ./pkg/opentelemetry/*.go ./pkg/opentelemetry/

RUN go build -o app

FROM alpine as main

COPY --from=builder /app/app /bin/app

EXPOSE 8080

ENTRYPOINT [ "/bin/app" ]