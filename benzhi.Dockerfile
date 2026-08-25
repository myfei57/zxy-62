FROM golang:1.23 AS build
ENV GOPROXY=off GOSUMDB=off CGO_ENABLED=0
WORKDIR /app
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY cmd ./cmd
COPY internal ./internal
COPY web ./web
RUN go build -mod=vendor -o /tunnelnet ./cmd/tunnelnet

FROM golang:1.23
ENV GOPROXY=off GOSUMDB=off
WORKDIR /app
COPY . .
COPY --from=build /tunnelnet /usr/local/bin/tunnelnet
CMD ["tunnelnet"]
