FROM golang:1.14-alpine AS build
ADD . /app
ENV CGO_ENABLED=0
WORKDIR /app
RUN go build -o app ./cmd

FROM alpine:latest
COPY --from=build /app/app /app/app
ENTRYPOINT ["/app/app"]
EXPOSE 9999