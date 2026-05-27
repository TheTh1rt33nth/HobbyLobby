.PHONY: docker-build kind-push deploy all

docker-build:
	docker build -t hobby-lobby:latest .

kind-push:
	kind load docker-image hobby-lobby:latest

deploy:
	kubectl apply -f postgres-statefulset.yaml
	kubectl apply -f server-deployment.yaml

all: docker-build kind-push deploy
