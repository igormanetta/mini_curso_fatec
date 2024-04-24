FROM golang:alpine AS builder

WORKDIR /src

COPY go.mod /src
COPY go.sum /src

RUN go mod download
ADD . /src

RUN go build -o /root/fileserver ./cmd/fileserver

FROM ubuntu

WORKDIR /app

RUN mkdir -p /basedir/downloads

COPY --from=builder /root/fileserver .
CMD [ "./fileserver" ]
