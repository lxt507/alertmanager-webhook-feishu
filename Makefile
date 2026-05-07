fmt:
	go fmt ./...
run:fmt
	go run main.go server -c config.yml -e -v
build: 
	goreleaser release --snapshot --clean
docker_build:
	docker build -t johnxu1989/alertmanager-webhook-feishu .
docker_push:docker_build
	docker push johnxu1989/alertmanager-webhook-feishu
