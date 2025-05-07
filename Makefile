COMPOSE_PROJECT_NAME=a1exander256-pow
LINTER_VERSION = v2.1.5

lint:
	@if ! bin/golangci-lint --version | grep -q $(LINTER_VERSION); \
		then curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s $(LINTER_VERSION); fi;
	bin/golangci-lint run --fix

up:
	COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) docker-compose up --build --force-recreate -d

down:
	COMPOSE_PROJECT_NAME=$(COMPOSE_PROJECT_NAME) docker-compose  down -v