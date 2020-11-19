build:
	docker-compose -f ./deploy/docker-compose.yml build --no-cache osnfs

run:
	docker-compose -f ./deploy/docker-compose.yml up --remove-orphans osnfs

test:
	docker-compose -f ./deploy/docker-compose.yml up tests

lint:
	docker-compose -f ./deploy/docker-compose.yml up linter

down:
	docker-compose -f deploy/docker-compose.yml down