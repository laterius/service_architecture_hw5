build:
	docker build -f docker/Dockerfile . -t 34234247632/otus-msa-hw5:v1.10

push:
	docker push 34234247632/otus-msa-hw5:v1.10

docker-start:
	cd docker && docker-compose up -d

docker-stop:
	cd docker && docker-compose down

k8s-pre-reqs:
	helm/pre-reqs.sh

k8s-deploy:
	helm/deploy.sh

k8s-remove:
	helm/remove.sh

# указать правильный порт перед запуском !!!
loadtest:
	ab -c 5 -n 1000000000 http://arch.homework:31841/api/user/1
