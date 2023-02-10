FROM golang:alpine3.17 as build

WORKDIR /app

COPY . .

RUN go build -o main ./src

FROM alpine:3.17

WORKDIR /app

COPY --from=build /app/main .

ENTRYPOINT ["/app/main"]