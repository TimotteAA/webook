.PHONY: docker
docker:
	@rm webook || true
	@GOOS=linux GOARCH=arm go build -tags=k8s -o webook .
	@docker rmi -f timotte/webook:v0.0.1
	@docker build -t timotte/webook-live:v0.0.1 .