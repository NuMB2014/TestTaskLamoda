FROM golang:alpine

WORKDIR /srv/server
COPY ./ /srv/server

RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o server cmd/main.go

CMD ["/srv/server/server"]