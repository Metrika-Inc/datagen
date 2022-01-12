FROM golang:1.17.6-alpine3.15 as builder

COPY . /src
WORKDIR /src

RUN go build -o data-service main.go

FROM alpine:3.15.0
COPY --from=builder /src/data-service /app/data-service
WORKDIR /app
ENTRYPOINT [ "/app/data-service" ]

