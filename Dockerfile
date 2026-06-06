FROM golang:1.24-alpine AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags "-s -w -X main.version=${VERSION}" \
    -o /out/nzbget-metrics ./cmd/nzbget-metrics

FROM scratch
COPY root/ /
COPY --from=build /out/nzbget-metrics /usr/local/bin/nzbget-metrics
