FROM golang:1.20-alpine AS builder

WORKDIR /app

COPY ./go.mod ./
COPY ./go.sum ./

RUN go mod download

COPY ./pkg/processor/*.go ./
COPY ./pkg/common/*.go ./pkg/common/
COPY ./pkg/opentelemetry/*.go ./pkg/opentelemetry/

RUN go build -o app

FROM alpine as main

RUN apk update && apk add bash

COPY ./pkg/processor/container_execution.sh /bin/container_execution.sh
RUN chmod +x /bin/container_execution.sh

COPY --from=builder /app/app /bin/app

ENTRYPOINT [ "/bin/app" ]