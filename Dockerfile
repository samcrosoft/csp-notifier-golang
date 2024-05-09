FROM golang:1.22.2-alpine as build
WORKDIR /app
COPY ./ /app

RUN go mod tidy \
    && go build -o /app/collector /app/main.go


FROM alpine:3.17
WORKDIR /app
COPY --from=build /app/collector /app
ENTRYPOINT ["/app/collector"]