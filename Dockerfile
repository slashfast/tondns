FROM golang:1.23 AS build
WORKDIR /build
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /bin/tondns

FROM alpine:latest
COPY --from=build /bin/tondns /bin/tondns
RUN chmod +x /bin/tondns

ENTRYPOINT ["/bin/tondns"]
