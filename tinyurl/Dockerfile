FROM golang:1.22

RUN mkdir app

COPY . /app

WORKDIR /app

RUN go mod download

WORKDIR /app/cmd

RUN CGO_ENABLED=0 GOOS=linux go build

ENTRYPOINT ["/app/cmd/cmd"]
