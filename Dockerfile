FROM golang:alpine as build-env

RUN apk update
RUN apk add -U --no-cache ca-certificates && update-ca-certificates
RUN apk add -U --no-cache tzdata

RUN mkdir /app
WORKDIR /app

COPY go.mod .
COPY go.sum .

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o /go/bin/ava-go

FROM scratch

COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build-env /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=build-env /go/bin/ava-go /go/bin/ava-go
COPY --from=build-env /app/config.json /config.json
COPY --from=build-env /app/channel_names.json /channel_names.json

ENTRYPOINT [ "/go/bin/ava-go" ]