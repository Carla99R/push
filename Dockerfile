from golang:1.16 AS builder

LABEL MAINTAINER="Carla R"

WORKDIR /src/
COPY . .
RUN go fmt $(go list ./... | grep -v /vendor/) &&\
    go vet $(go list ./... | grep -v /vendor/) &&\
    go test -race $(go list ./... | grep -v /vendor/) &&\
    GOOS=linux go build -a -o bin/app *.go

FROM debian:buster-slim
WORKDIR /push/
COPY --from=builder ["/src/bin/app", "/src/.env*", "/src/docs", "/v1"]

EXPOSE 80
CMD ["./app"]