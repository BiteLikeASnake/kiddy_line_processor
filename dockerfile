FROM golang:1.12.0-alpine3.9 AS builder
WORKDIR /go/src/github.com/BiteLikeASnake/kiddy_line_processor
COPY . .
RUN go install ./...

#FROM jwilder/dockerize AS production
#COPY --from=builder /go/bin/cmd ./app

#docker build -t docker/line_processor
#docker stop kiddy_test
#docker rm kiddy_test
#docker run -it --name kiddy_test line_processor_img /bin/sh