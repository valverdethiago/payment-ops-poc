FROM golang:alpine

RUN apk add --no-cache bash

# Set necessary environmet variables needed for our image
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

WORKDIR /env
COPY env/docker.env ./app.env

WORKDIR /src

COPY src .

RUN go mod download

RUN go build -o /bin/bank-transactional-ms

EXPOSE ${PORT}

CMD ["/bin/bank-transactional-ms"]
