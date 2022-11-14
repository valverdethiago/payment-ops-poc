stack_start:
	docker-compose -f ./docker/docker-compose-full.yml  up -d --force-recreate --build --remove-orphans

stack_stop:
	docker-compose -f ./docker/docker-compose-full.yml  down

dev_start:
	docker-compose -f ./docker/docker-compose.yml  up -d --build --remove-orphans

dev_stop:
	docker-compose -f ./docker/docker-compose.yml  down 