FROM golang:1.12.0-alpine3.9 AS builder
WORKDIR /go/src/github.com/BiteLikeASnake/kiddy_line_processor
COPY . .
RUN go install ./...

FROM jwilder/dockerize AS production
COPY --from=builder /go/bin/cmd ./app

#docker build -t line_processor_img .
#docker stop kiddy_test
#docker rm kiddy_test
#docker run -it --name kiddy_test line_processor_img /bin/sh

#- /bin/sh -c "dockerize -wait http://db:5432 -timeout 30s"

#docker run -p8000:8000 -d --restart always antonboom/lines-provider
#