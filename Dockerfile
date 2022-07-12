FROM golang:1.18rc1-alpine3.15
RUN mkdir post_loader_service
COPY . /post_loader_service
WORKDIR /post_loader_service
RUN go build -o main cmd/main.go
CMD ./main