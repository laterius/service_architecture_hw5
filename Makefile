build:
	docker build -f docker/Dockerfile . -t 34234247632/otus-msa-hw5:v2.2

push:
	docker push 34234247632/otus-msa-hw5:v2.2

docker-start:
	cd docker && docker-compose up -d

docker-stop:
	cd docker && docker-compose down

k8s-deploy:
	kubectl create ns otus-msa-hw5
	helm upgrade --install -n otus-msa-hw5 otus-msa-hw5 helm/chart

k8s-remove:
	kubectl delete ns otus-msa-hw5

newman:
	newman run postman/collection.json

