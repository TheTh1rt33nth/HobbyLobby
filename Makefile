.PHONY: docker-build kind-push db migrate deploy all \
        frontend-dev frontend-build frontend-test frontend-image frontend-kind-push frontend-deploy \
        proxy

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
	kubectl rollout restart deployment/hobby-lobby -n hobby-lobby
	kubectl rollout status deployment/hobby-lobby -n hobby-lobby --timeout=120s

all: docker-build kind-push deploy

frontend-dev:
	cd frontend && npm run dev

frontend-build:
	cd frontend && npm run build

frontend-test:
	cd frontend && npm run test

frontend-image:
	docker build -t hobby-lobby-frontend:latest ./frontend

frontend-kind-push: frontend-image
	kind load docker-image hobby-lobby-frontend:latest

frontend-deploy: frontend-kind-push
	kubectl apply -f deploy/frontend-deployment.yaml
	kubectl apply -f deploy/ingress.yaml
	kubectl rollout restart deployment/hobby-lobby-frontend -n hobby-lobby
	kubectl rollout status deployment/hobby-lobby-frontend -n hobby-lobby --timeout=120s

## Build and deploy both backend + frontend end-to-end
all-full: docker-build kind-push deploy frontend-kind-push frontend-deploy

## Port-forward the ingress controller — opens the full app at http://localhost:8888
proxy:
	kubectl port-forward svc/ingress-nginx-controller 8888:80 -n ingress-nginx
