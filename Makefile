.PHONY: docker-build kind-push db migrate deploy all

docker-build:
	docker build -t hobby-lobby:latest ./backend

kind-push:
	kind load docker-image hobby-lobby:latest

db:
	kubectl apply -f deploy/postgres-statefulset.yaml
	kubectl rollout status statefulset/postgres -n hobby-lobby-db --timeout=120s

migrate: db
	kubectl apply -f deploy/app-config.yaml
	kubectl delete job hobby-lobby-migrate -n hobby-lobby --ignore-not-found
	kubectl apply -f deploy/migrate-job.yaml
	kubectl wait --for=condition=complete job/hobby-lobby-migrate -n hobby-lobby --timeout=120s
	kubectl logs job/hobby-lobby-migrate -n hobby-lobby

deploy: migrate
	kubectl apply -f deploy/server-deployment.yaml

all: docker-build kind-push deploy
