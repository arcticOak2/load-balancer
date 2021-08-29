# base image
FROM golang:1.16-alpine

# init go
WORKDIR /app
COPY go.mod ./
COPY go.sum ./
RUN go mod download
COPY . ./
RUN go build -o /docker-load-balancer

CMD ["/docker-load-balancer"]