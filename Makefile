DOCKER_USERNAME ?= narmitag
APPLICATION_NAME ?= weather
 
build:
	docker build --tag ${DOCKER_USERNAME}/${APPLICATION_NAME} .
push:
	docker push ${DOCKER_USERNAME}/${APPLICATION_NAME}
run:
	docker run --publish 8081:8081 ${DOCKER_USERNAME}/${APPLICATION_NAME}