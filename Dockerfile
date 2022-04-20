FROM golang:1-alpine AS build-env

RUN apk add git

RUN mkdir /app
COPY . /app

WORKDIR /app

RUN go build -o /tmp/groxy ./cmd/groxy

FROM alpine:latest

COPY --from=build-env /tmp/groxy /bin/groxy

RUN ls /bin/groxy

CMD ["/bin/groxy"]

