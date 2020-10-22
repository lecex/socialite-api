.PHONY: git
git:
	git add .
	git commit -m"自动提交 git 代码"
	git push
.PHONY: tag
tag:
	git push --tags
.PHONY: rpc
rpc:
	micro api  --handler=rpc  --namespace=go.micro.api --address=:8080
.PHONY: api
api:
	micro api  --handler=api  --namespace=go.micro.api --address=:8081
.PHONY: proto
proto:
	protoc -I . --micro_out=. --gogofaster_out=. proto/socialite/socialite.proto
	protoc -I . --micro_out=. --gogofaster_out=. proto/user/user.proto
	protoc -I . --micro_out=. --gogofaster_out=. proto/config/config.proto

.PHONY: docker
docker:
	docker build -f Dockerfile  -t socialite-api .
.PHONY: run
run:
	go run main.go
test:
	go test main_test.go -test.v