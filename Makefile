all: tests

storage_img:
	docker build -t postgr_storage_img docker_storage

line_processor_img:
	docker build -t line_processor_img .

lint:
	go get -u golang.org/x/lint/golint
	golint cmd internal

tests:
	go test ./... -v

run: storage_img line_processor_img
	docker-compose up -d

stop:
	docker-compose down

#all: tests